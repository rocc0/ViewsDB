package main

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

//Generate returns a random seq of symbols
func Generate(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
