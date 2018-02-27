package main

import (
	"log"
)

func init() {
	var c Config

	if err := c.getConf(); err != nil {
		log.Fatal("Error when parsing config: %v", err)
	}
}

func main() {
	if err := Run(); err != nil {
		log.Printf("Error in main(): %v", err)
	}
}
