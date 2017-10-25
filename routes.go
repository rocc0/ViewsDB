package main

import (
	bleveHttp "github.com/blevesearch/bleve/http"
	"github.com/gin-gonic/gin"
)



func initializeRoutes() {

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not

	router.Use(setUserStatus())

	// Handle the index route
	router.GET("/", showIndexPage)
	router.GET("/ratings", showRatings)
	userRoutes := router.Group("/u")
	{
		// Handle the GET requests at /u/login
		// Show the login page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)

		// Handle POST requests at /u/login
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/login", ensureNotLoggedIn(), performLogin)

		// Handle GET requests at /u/logout
		// Ensure that the user is logged in by using the middleware
		userRoutes.GET("/logout", ensureLoggedIn(), logout)

		// Handle the GET requests at /u/register
		// Show the registration page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/register", ensureNotLoggedIn(), showRegistrationPage)

		// Handle POST requests at /u/register
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/register", ensureNotLoggedIn(), register)
	}

	viewRoutes := router.Group("/views")
	{
		viewRoutes.GET("/view/:view_id", getView)

		//viewRoutes.POST("/view/:view_id/edit", ensureLoggedIn(), showEditView)
		//
		//viewRoutes.POST("/view/:view_id/edit", ensureLoggedIn(), editView)

		viewRoutes.GET("/create",  viewCreationPage)

		// Handle POST requests at /article/create
		// Ensure that the user is logged in by using the middleware
		viewRoutes.POST("/create",  createView)

	}

	apiRoutes := router.Group("/api")
	{
		//Get goverments names and ids
		apiRoutes.GET("/govs", allGovs)

		//View ratings
		apiRoutes.GET("/ratings", viewRatings)

		//Show and edit view
		apiRoutes.GET("/v/:view_id", viewData)
		apiRoutes.POST("/v/:view_id", saveViewRow)

		//Creation of view
		apiRoutes.POST("/create", createView)

		//Search handling
		apiRoutes.POST("/search", gin.WrapH(bleveHttp.NewSearchHandler("view")))
		apiRoutes.GET("/fields", gin.WrapH(bleveHttp.NewListFieldsHandler("view")))

	}
}
