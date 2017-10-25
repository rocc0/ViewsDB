package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"fmt"
)

func showIndexPage(c *gin.Context) {
	articles := getAllArticles()

	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title":   "Home Page",
		"payload": articles}, "index.html")
}



func getView(c *gin.Context) {
	// Check if the article ID is valid
	viewID, err := strconv.Atoi(c.Param("view_id"));
	if err == nil {
		// Check if the article exists
		view, err := getViewById(viewID)
		if err == nil {
			// Call the render function with the title, article and the name of the
			// template
			for _,v := range view.Fields {
				fmt.Sprintf("%v", v)
			}
			render(c, gin.H{
				"title": "Редагування | " + string(view.Fields["name_and_requisits"].([]uint8)),
				"pl": view.Fields}, "index.html")

		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusBadGateway)
	}
}

func viewCreationPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Додати відстеження",
	}, "index.html")
}

func showRatings(c *gin.Context) {
	render(c, gin.H{
		"title": "Таблиця по міністерствам",
	}, "index.html")
}

