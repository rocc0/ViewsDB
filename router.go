package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	//gpprof "github.com/gin-contrib/pprof"
)

// Use the setUserStatus middleware for every route to set a flag
// indicating whether the request was from an authenticated user or not

var router *gin.Engine

func initializeRoutes() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)
	// Set the router as the default one provided by Gin
	router = gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20
	// Set static routes
	router.Static("static/", config.Assets)
	// Set favicon path
	router.StaticFile("/favicon.ico", "static/favicon.ico")
	//Set templates path
	if tmpl, err := template.New("projectViews").ParseGlob("templates/*"); err == nil {
		router.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}

	// Handle the index route
	router.GET("/", showIndexPage)
	router.GET("/ratings", showRatings)
	router.GET("/govs/edit", showEditGovsNames)

	userRoutes := router.RouterGroup.Group("/u")

	{
		userRoutes.GET("/login", showIndexPage)

		userRoutes.POST("/login", authMiddleware.LoginHandler)

		userRoutes.GET("/register", showIndexPage)

		userRoutes.POST("/register", register)

		userRoutes.GET("/cabinet", showUserPage)

	}

	traceRoutes := router.Group("/track")
	{
		traceRoutes.GET("/id/:trk_id", showTracePage)

		traceRoutes.GET("/id/:trk_id/edit", showTracePage)

		traceRoutes.GET("/create", showTraceCreationPage)

	}

	apiTraceRoutes := router.Group("/api")
	{
		//Get goverments names and ids
		apiTraceRoutes.GET("/govs", getGoverns)
		apiTraceRoutes.POST("/govs/edit", authMiddleware.MiddlewareFunc(), postEditGovernments)

		//View ratings
		apiTraceRoutes.GET("/ratings", getRatings)

		//Show and edit view
		apiTraceRoutes.GET("/v/:trk_id", getTrace)
		apiTraceRoutes.POST("/v/:trk_id", authMiddleware.MiddlewareFunc(), postTrackField)

		//Creation of a new trace
		apiTraceRoutes.POST("/create", authMiddleware.MiddlewareFunc(), postCreateItem)

		//Creation of periodic trace
		apiTraceRoutes.POST("/create-period", authMiddleware.MiddlewareFunc(), postCreateItem)

		//Delete handling
		apiTraceRoutes.POST("/delete", authMiddleware.MiddlewareFunc(), postDeleteItem)

	}
	apiImageRoutes := router.Group("/api")
	{
		//Images
		apiImageRoutes.POST("/upload", postAddImage)

		apiImageRoutes.GET("/img/:trk_id", getTraceImages)

		apiImageRoutes.POST("/img/:trk_id/delete", authMiddleware.MiddlewareFunc(), postDelImage)
	}

	apiUserRoutes := router.Group("/api")
	{
		//user
		apiUserRoutes.GET("/cabinet", authMiddleware.MiddlewareFunc(), cabinet)
		apiUserRoutes.POST("/edituser", authMiddleware.MiddlewareFunc(), editField)
	}
	log.Printf("Starting, HTTP on: %s\n", config.Listen)
	//gpprof.Register(router, &gpprof.Options{
	//	// default is "debug/pprof"
	//	RoutePrefix: "debug/pprof",
	//})
	router.Run(config.Listen)
}

func render(c *gin.Context, data gin.H, templateName string) {

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
