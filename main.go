package main

import (


	"github.com/gin-gonic/gin"
	"flag"
	"log"

)

var router *gin.Engine

const imgpath = "/static/images/"

var assetsPath string

func processFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ListenSpec, "listen", ":8888", "HTTP listen spec")
	flag.StringVar(&cfg.Db.ConnectString, "db-connect",
		"root:password@tcp(192.168.99.100:3306)/trackdb", "DB Connect String")
	flag.StringVar(&assetsPath, "assets-path", "static/", "Path to assets dir")
	flag.StringVar(&cfg.Rou.Cpuprofile,"cpuprofile", "./static/cpu.out", "write cpu profile to file")
	flag.StringVar(&cfg.Rou.Memprofile,"memprofile", "./static/mem.out", "write mem profile to file")

	flag.Parse()
	return cfg
}

func main() {

	cfg := processFlags()

	if err := Run(cfg); err != nil {
		log.Printf("Error in main(): %v", err)
	}
}
