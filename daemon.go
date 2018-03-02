package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"syscall"
)

func run() error {
	ch := make(chan bool, 1)

	//Flag parsing
	reindex := flag.NewFlagSet("reindex", flag.ExitOnError)
	reindexWorkers := reindex.Int("qty", 1, "Number of workers")
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

	go initializeRoutes()

	//Reindexing
	if reindex.Parsed() {
		if *reindexWorkers <= 0 {
			reindex.Usage()
			os.Exit(1)
		}
		doReindex(reindexWorkers, ch)
	}

	//Calculation of reports
	//go calculateRates(c)
	//<-c
	waitForSignal()

	return nil
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
