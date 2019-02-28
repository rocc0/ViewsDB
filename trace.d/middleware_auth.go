// middleware.auth.go

package main

import (
	"encoding/json"
	"errors"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

var authMiddleware = &jwt.GinJWTMiddleware{
	Realm:         "test zone",
	Key:           []byte("secret key"),
	Timeout:       time.Hour,
	MaxRefresh:    time.Hour,
	Authenticator: authenticator,
	Authorizator:  authorizator,
	Unauthorized:  unauthorizator,
	TokenLookup:   "header:Authorization",
	TokenHeadName: "Bearer",
	TimeFunc:      time.Now,
}

func authenticator(c *gin.Context) (interface{}, error) {
	var u User
	var b []byte
	if _, err := c.Request.Body.Read(b); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &u); err != nil {
		return nil, err
	}

	if u.loginCheck() {
		return u.Email, nil
	}

	return nil, errors.New("unknown user")
}

func authorizator(email interface{}, c *gin.Context) bool {
	u := User{Email: email.(string)}
	return u.authCheck()
}

func unauthorizator(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
