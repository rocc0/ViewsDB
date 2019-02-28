package main

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/rocc0/TraceDB/trace.d/docs" // docs is generated by Swag CLI, you have to import it.
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var router *gin.Engine

func initializeRoutes() error {
	gin.SetMode(gin.TestMode)
	router = gin.Default()
	router.MaxMultipartMemory = 8 << 20
	router.Static("static/", config.Assets)
	router.StaticFile("/favicon.ico", "static/favicon.ico")
	if tmpl, err := template.New("projectViews").ParseGlob("trace.d/static/templates/*"); err == nil {
		router.SetHTMLTemplate(tmpl)
	} else {
		return err
	}

	router.GET("/", showIndexPage)
	router.GET("/ratings", showRatingsPage)
	router.GET("/govs/edit", showEditGovsNames)

	userRoutes := router.RouterGroup.Group("/u")
	userRoutes.GET("/login", showLoginPage)
	userRoutes.POST("/login", authMiddleware.LoginHandler)
	userRoutes.GET("/register", showRegisterPage)
	userRoutes.POST("/register", register)
	userRoutes.GET("/cabinet", showUserPage)

	traceRoutes := router.Group("/track")
	traceRoutes.GET("/id/:trk_id", showTracePage)
	traceRoutes.GET("/id/:trk_id/edit", showTracePage)
	traceRoutes.GET("/create", showTraceCreationPage)

	apiTraceRoutes := router.Group("/api")
	apiTraceRoutes.GET("/govs", getGoverns)
	apiTraceRoutes.POST("/govs/edit", authMiddleware.MiddlewareFunc(), postEditGovernments)
	apiTraceRoutes.GET("/ratings", getRatings)
	apiTraceRoutes.GET("/v/:trk_id", getTrace)
	apiTraceRoutes.POST("/v/:trk_id", authMiddleware.MiddlewareFunc(), postEditTrackField)
	apiTraceRoutes.POST("/create", authMiddleware.MiddlewareFunc(), postCreateItem)
	apiTraceRoutes.POST("/create-period", authMiddleware.MiddlewareFunc(), postCreateItem)
	apiTraceRoutes.POST("/delete", authMiddleware.MiddlewareFunc(), postDeleteItem)

	apiImageRoutes := router.Group("/api")
	apiImageRoutes.POST("/upload", postAddImage)
	apiImageRoutes.GET("/img/:trk_id", getTraceImages)
	apiImageRoutes.POST("/img/:trk_id/delete", authMiddleware.MiddlewareFunc(), postDelImage)

	apiUserRoutes := router.Group("/api")
	apiUserRoutes.GET("/cabinet", authMiddleware.MiddlewareFunc(), cabinet)
	apiUserRoutes.POST("/edituser", authMiddleware.MiddlewareFunc(), editField)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	Logger.INFO("Starting, HTTP on:", config.Listen)
	return router.Run(config.Listen)
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