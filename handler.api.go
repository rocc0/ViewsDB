package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"log"

)


func getViewJson(c *gin.Context) {
	viewID, err := strconv.Atoi(c.Param("view_id"));
	if err == nil {
		// Check if the article exists
		view, err := getViewById(viewID)
		if err == nil {
			// Call the render function with the title, article and the name of the
			// template
			c.SecureJSON(http.StatusOK, gin.H{
				"pl": view.Fields})
			log.Print(view.Fields)

		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusBadGateway)
	}
}