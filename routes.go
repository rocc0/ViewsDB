package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
)

func initializeRoutes() {

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not

	authMiddleware := &jwt.GinJWTMiddleware {
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(eMail string, password string, c *gin.Context) (string, bool) {
			if (loginCheck(eMail, password) == true) {
				return eMail, true
			}

			return eMail, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			return authCheck(userId)
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		TokenLookup: "header:Authorization",

		TokenHeadName: "Bearer",

		TimeFunc: time.Now,
	}

	// Handle the index route
	router.GET("/", showIndexPage)
	router.GET("/ratings", showRatings)
	router.GET("/govs/edit", showEditGovsNames)

	userRoutes := router.Group("/u")

	{
		userRoutes.GET("/login", showIndexPage)

		userRoutes.POST("/login", authMiddleware.LoginHandler)

		userRoutes.GET("/register", showIndexPage)

		userRoutes.POST("/register", registerHandler)

		userRoutes.GET("/cabinet", showUserPage)

	}

	trackRoutes := router.Group("/track")
	{
		trackRoutes.GET("/id/:trk_id", getView)

		trackRoutes.GET("/id/:trk_id/edit", getView)

		trackRoutes.GET("/create", viewCreationPage)

	}

	apiRoutes := router.Group("/api")
	{
		//Get goverments names and ids
		apiRoutes.GET("/govs", getGovernments)
		apiRoutes.POST("/govs/edit", authMiddleware.MiddlewareFunc(), postEditGovernments)

		//View ratings
		apiRoutes.GET("/ratings", getRatings)

		//Show and edit view
		apiRoutes.GET("/v/:trk_id", getTrack)
		apiRoutes.POST("/v/:trk_id", authMiddleware.MiddlewareFunc(),postTrackField)

		//Creation of view
		apiRoutes.POST("/create", authMiddleware.MiddlewareFunc(),postCreateItem)

		//Delete handling
		apiRoutes.POST("/delete", authMiddleware.MiddlewareFunc(),postDeleteItem)

		//Images
		apiRoutes.POST("/upload",  postImage)

		apiRoutes.GET("/img/:trk_id", getImages)

		apiRoutes.POST("/img/:trk_id/delete", authMiddleware.MiddlewareFunc(), postDelImage)

		//user
		apiRoutes.GET("/cabinet", authMiddleware.MiddlewareFunc(), cabinetHandler)
		apiRoutes.POST("/edituser",authMiddleware.MiddlewareFunc(), editUserField)
	}

}
