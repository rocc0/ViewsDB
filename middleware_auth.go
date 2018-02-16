// middleware.auth.go

package main

import (
	"time"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

var authMiddleware = &jwt.GinJWTMiddleware{
	Realm:      "test zone",
	Key:        []byte("secret key"),
	Timeout:    time.Hour,
	MaxRefresh: time.Hour,
	Authenticator: func(email string, password string, c *gin.Context) (string, bool) {
		var u User
		u.Email = email
		u.Password = password
		if u.LoginCheck() == true {
			return email, true
		}

		return email, false
	},
	Authorizator: func(email string, c *gin.Context) bool {
		var u User
		return u.AuthCheck()
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