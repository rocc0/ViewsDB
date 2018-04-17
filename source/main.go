/*
 * Copyright (c) 2018.
 */

package main

import (
	"log"
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
	if err := config.getConf(); err != nil {
		log.Fatalf("Error when parsing config: %v", err)
	}
	err := initDB(config.MySQL)
	if err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}
	//err = mgoConnect()
	//if err != nil {
	//	log.Fatalf("Error initializing mongo: %v\n", err)
	//}
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error in main(): %v", err)
	}
}
