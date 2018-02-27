package main

import (
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)


func Run() error {
	ch, _ := make(chan bool), make(chan int)
	err := InitDb(config.MySql)
	if err != nil {
		log.Printf("Error initializing database: %v\n", err)
		return err
	}

	err = mgoConnect()
	if err != nil {
		log.Printf("Error initializing mongo: %v\n", err)
		return err
	}

	if config.CpuProf != "" {
		f, err := os.Create(config.CpuProf)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		f.Close()
	}

	go initializeRoutes()
	//go calculateRates(c)
	//<-c

	createWorkerPool(10, ch)

	pprof.StopCPUProfile()
	if config.MemProf != "" {
		f, err := os.Create(config.MemProf)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
	<-ch

	WaitForSignal()

	return nil
}

func WaitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
