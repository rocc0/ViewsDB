/*
 * Copyright (c) 2018.
 */

package main

import (
	"flag"
	"log"
	"os"
)

// @title TraceDB API
// @version 1.0
// @description This is a TraceDB celler server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8888
// @BasePath /api

func init() {
	initLogger()
	if err := config.getConf(); err != nil {
		log.Fatalf("Error when parsing config: %v", err)
	}

	if err := initDB(); err != nil {
		log.Printf("Error initializing database: %v\n", err)
	}

	ch := make(chan bool, 1)
	//Flag parsing
	reindex := flag.NewFlagSet("reindex", flag.ExitOnError)
	reindexWorkers := reindex.Int("qty", 1, "Number of workers")
	if len(os.Args) > 2 {
		switch os.Args[1] {
		case "reindex":
			if err := reindex.Parse(os.Args[2:]); err != nil {
				log.Fatal(err)
			}
		default:
			log.Print("Starting...")
		}
	}

	go func() {
		if err := initializeRoutes(); err != nil {
			log.Fatal(err)
		}
	}()

	//Reindexing
	if reindex.Parsed() {
		if *reindexWorkers <= 0 {
			reindex.Usage()
			os.Exit(1)
		}
		doReindex(reindexWorkers, ch)
	}
}

func main() {
	waitForSignal()
}
