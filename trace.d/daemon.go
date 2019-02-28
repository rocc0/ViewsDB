package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() error {
	dbInfo := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=tracedb sslmode=disable",
		config.PgHost, config.PgUser, config.PgPass)
	newDB, err := sql.Open("postgres", dbInfo)
	if err != nil {
		Logger.ERROR(err)
		return err
	}
	db = newDB
	if err = db.Ping(); err != nil {
		Logger.ERROR(err)
		return err
	}

	if err = userInit(); err != nil {
		Logger.ERROR(err)
		return err
	}
	return nil

}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
