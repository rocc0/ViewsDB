package main


import (
	"net/http"
	"github.com/gin-gonic/gin"
	"html/template"

	bleveHttp "github.com/blevesearch/bleve/http"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var router *gin.Engine

func main(){
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	router = gin.Default()


	router.Static("js/", "js/")
	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	//router.LoadHTMLGlob("templates/*")
	if tmpl, err := template.New("projectViews").Funcs(TemplateHelpers).ParseGlob("templates/*"); err == nil {
		router.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}
	userInit()
	// Initialize the routes
	initializeRoutes()
	generateIndexes()
	idx, _ := Bleve(viewsIdx)
	bleveHttp.RegisterIndexName("view", idx)
	log.Printf("Indexing complited!")

	// Start serving the application
	router.Run(":8888")
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