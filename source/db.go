package main

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

func initDB(connect string) error {
	var err error
	dbinfo := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=tracedb sslmode=disable",
		config.PgHost, config.PgUser, config.PgPass)
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = userInit()
	if err != nil {
		log.Print(err.Error())
		return err
	}
	return nil

}
