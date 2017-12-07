package main


import (
	"html/template"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main(){
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Set static routes
	router.Static("static/", "static/")
	router.StaticFile("/favicon.ico", "static/favicon.ico")


	if tmpl, err := template.New("projectViews").Funcs(TemplateHelpers).ParseGlob("templates/*"); err == nil {
		router.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}
	userInit()
	// Initialize the routes
	initializeRoutes()
	//Search indexing
	elasticIndex()

	calculateRates()
	// Start serving the application
	router.Run(":8888")
}
