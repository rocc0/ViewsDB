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
		u := user{Email: email, Password: password}

		if u.loginCheck() == true {
			return email, true
		}

		return email, false
	},
	Authorizator: func(email string, c *gin.Context) bool {
		u := user{Email: email}
		return u.authCheck()
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
