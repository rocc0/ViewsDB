// handlers.user.go

package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

type userField struct {
	Field string `json:"field"`
	Data  string `json:"data"`
	ID    int    `json:"id"`
}

func cabinet(c *gin.Context) {
	var u user
	claims := jwt.ExtractClaims(c)
	u.Email, _ = claims["id"].(string)
	err := u.getUser()

	if err != nil {
		c.JSON(400, gin.H{
			"data": "user not found",
		})
	} else {
		c.JSON(200, gin.H{
			"data": &u,
		})
	}
}

func register(c *gin.Context) {
	var u user
	x, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal([]byte(x), &u)

	if err != nil {
		c.AbortWithStatus(400)
	}

	if err := u.register(); err == nil {
		c.JSON(200, gin.H{
			"title": "Реєстрація успішна!",
			"id":    u.Name,
		})
	} else {
		c.AbortWithStatus(400)
	}
}

func editField(c *gin.Context) {
	var f userField
	x, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal([]byte(x), &f)

	if err != nil {
		c.AbortWithStatus(400)
	}

	if err := f.editField(); err == nil {
		c.JSON(200, gin.H{
			"title": "Змінено!",
		})
	} else {
		c.AbortWithStatus(400)
	}
}
