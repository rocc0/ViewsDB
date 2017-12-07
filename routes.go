package main


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

	viewRoutes := router.Group("/track")
	{
		viewRoutes.GET("/id/:trk_id", getView)
		viewRoutes.GET("/id/:trk_id/edit", getView)

		//viewRoutes.POST("/view/:view_id/edit", ensureLoggedIn(), showEditView)
		//
		//viewRoutes.POST("/view/:view_id/edit", ensureLoggedIn(), editView)

		viewRoutes.GET("/create",  viewCreationPage)
	}

	apiRoutes := router.Group("/api")
	{
		//Get goverments names and ids
		apiRoutes.GET("/govs", getGovernments)

		//View ratings
		apiRoutes.GET("/ratings", getRatings)

		//Show and edit view
		apiRoutes.GET("/v/:trk_id", getTrack)
		apiRoutes.POST("/v/:trk_id", postTrackField)

		//Creation of view
		apiRoutes.POST("/create", postCreateItem)

		//Delete handling
		apiRoutes.POST("/delete", postDeleteItem)
	}
}
