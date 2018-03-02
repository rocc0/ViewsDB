package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

func initDB(connect string) error {
	var err error
	db, err = sql.Open("mysql", connect)
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
