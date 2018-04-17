// handlers.user.go

package main

import (
	"encoding/json"
	"io/ioutil"

	"net/http"

	"./httputil"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// Cabinet godoc
// @Summary Getting user data for displaying in user cabinet
// @Description get string by ID
// @ID get-string-by-string
// @Accept  json
// @Produce  json
// @Param id path string true "JWT ID"
// @Success 200 {object} main.User
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/cabinet [get]
func cabinet(c *gin.Context) {
	var u User
	claims := jwt.ExtractClaims(c)
	u.Email, _ = claims["id"].(string)
	err := u.getUser()
	if err != nil {
		httputil.NewError(c, http.StatusNotFound, err)
		return
	}
	c.JSON(200, gin.H{
		"data": &u,
	})
}

// Register godoc
// @Summary Handling user registration
// @Description register a new user
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Реєстрація успішна!"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /u/register [get]
func register(c *gin.Context) {
	var u User
	x, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal([]byte(x), &u)

	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	if err := u.register(); err != nil {
		httputil.NewError(c, http.StatusNotFound, err)
		return
	}

	c.JSON(200, gin.H{
		"title": "Реєстрація успішна!",
		"id":    u.Name,
	})
}

// editField godoc
// @Summary Handling changes in user fields
// @Description this endpoint used to hendling passchange or first/lastname changing
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Змінено!"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/edituser [post]
func editField(c *gin.Context) {
	var f userField
	x, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal([]byte(x), &f)

	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	if err := f.editField(); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, gin.H{
		"title": "Змінено!",
	})
}
