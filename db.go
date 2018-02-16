package main

import (
	"database/sql"
	"log"
)


var db *sql.DB

type dConfig struct {
	ConnectString string
}

func InitDb(cfg dConfig) error {
	var err error
	db, err = sql.Open("mysql", cfg.ConnectString)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = UserInit()
	if err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}




