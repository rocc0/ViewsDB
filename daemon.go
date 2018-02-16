package main

import (
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

)

type rConfig struct {
	Cpuprofile string
	Memprofile string
}


type Config struct {
	ListenSpec string
	Rou        rConfig
	Db         dConfig
}

func Run(cfg *Config) error {
	e, c := make(chan int), make(chan int)
	err := InitDb(cfg.Db)
	if err != nil {
		log.Printf("Error initializing database: %v\n", err)
		return err
	}

	if cfg.Rou.Cpuprofile != "" {
		f, err := os.Create(cfg.Rou.Cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		f.Close()
	}

	go initializeRoutes(cfg.ListenSpec)

	go calculateRates(c)
	<-c

	go func() {
		err = elasticIndex(e)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StopCPUProfile()
		if cfg.Rou.Memprofile != "" {
			f, err := os.Create(cfg.Rou.Memprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.WriteHeapProfile(f)
			f.Close()
		}
	}()
	<-e

	WaitForSignal()

	return nil
}

func WaitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
