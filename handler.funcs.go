package main

import (
	"math/rand"
	"time"
	"html/template"
	"log"
	"net/http"
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)


var db *sql.DB

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var TemplateHelpers = template.FuncMap {
	"toString": func(s []uint8) string {
		return string(s)
	},
}

func check(e error) {
	if e != nil {
		log.Print(e.Error())
	}
}

func init() {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(192.168.99.100:3306)/db")
	check(err)

	err = db.Ping()
	check(err)
}


func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


func render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["pl"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["pl"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}