/*
 * Copyright (c) 2018.
 */

package main

import (
	"log"
)

func init() {
	if err := config.getConf(); err != nil {
		log.Fatalf("Error when parsing config: %v", err)
	}
	err := initDB(config.MySQL)
	if err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}
	err = mgoConnect()
	if err != nil {
		log.Fatalf("Error initializing mongo: %v\n", err)
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error in main(): %v", err)
	}
}
