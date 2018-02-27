package main

import (
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"flag"
)


func Run() error {
	ch, _ := make(chan bool), make(chan int)

	//Flag parsing
	reindex := flag.NewFlagSet("reindex", flag.ExitOnError)
	reindexWorkers := reindex.Int("qty",  1, "Number of workers")
	if len(os.Args) > 2 {
		switch os.Args[1] {
		case "reindex":
			err := reindex.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		default:
			log.Print("Strating...")
		}
	}

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

	go initializeRoutes()

	//Calculation of reports
	//go calculateRates(c)
	//<-c

	//Reindexing
	if reindex.Parsed() {
		if *reindexWorkers <= 0 {
			reindex.Usage()
			os.Exit(1)
		}

		createWorkerPool(*reindexWorkers, ch)

		if config.CpuProf != "" {
			f, err := os.Create(config.CpuProf)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			f.Close()
		}

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
	}

	WaitForSignal()

	return nil
}

func WaitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
