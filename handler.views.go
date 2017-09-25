package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
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
			render(c, gin.H{
				"title":  string(view.Fields["name_and_requisits"].([]uint8)),
				"pl": view.Fields}, "view.html")

		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusBadGateway)
	}
}

func showViewCreationPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Create view",
	}, "create.html")

}

func  createView(c *gin.Context) {
	cols := make([]string, len(getColsNames()))
	formData := make([]interface{}, len(getColsNames()))
	colsNames := getColsNames()
	for i, _ := range cols{
		formData[i] = c.PostForm(colsNames[i])
	}
	if a, err := createNewView(formData); err == nil {
		render(c, gin.H{
			"title":   "Submission Successful",
			"pl": a}, "success.html")
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}