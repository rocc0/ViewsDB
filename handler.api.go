package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"reflect"
)


func getViewJson(c *gin.Context) {
	viewID, err := strconv.Atoi(c.Param("view_id"));
	if err == nil {
		// Check if the article exists
		view, err := getViewById(viewID)
		arr := make(map[string]string)
		for k,v := range view.Fields {
			s := reflect.ValueOf(v)
			switch s.Interface().(type){
			case int64:
				arr[k] = string(v.(int64))
			default:
				arr[k] = string(v.([]uint8))
			}
		}
		if err == nil {
			// Call the render function with the title, article and the name of the
			// template
			c.JSON(http.StatusOK, gin.H{
				"pl": arr})
		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusBadGateway)
	}
}
