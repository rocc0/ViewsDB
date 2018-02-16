package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Пошук відстежень",
	}, "index.html")
}

func showTraceCreationPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Додати відстеження",
	}, "index.html")
}

func showRatings(c *gin.Context) {
	render(c, gin.H{
		"title": "Якісні показники",
	}, "index.html")
}

func showUserPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Кабінет користувача",
	}, "index.html")
}

func showEditGovsNames(c *gin.Context) {
	render(c, gin.H{
		"title": "Редагувати назви державних органів",
	}, "index.html")
}

func showTracePage(c *gin.Context) {
	var (
		trace Trace
		title string
	)
	trackId, err := strconv.Atoi(c.Param("trk_id"))
	url := c.Request.URL.Path
	if strings.Contains(url, "edit") {
		title = "Редагування | "
	} else {
		title = ""
	}
	if err == nil {
		err := trace.getBasicData(trackId)
		if err == nil {
			render(c, gin.H{
				"title": title + string(trace.Fields["requisits"].([]uint8)),
			}, "index.html")
		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		c.AbortWithStatus(http.StatusBadGateway)
	}
}
