package main

import (
	"net/http"
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
		trace BasicTrace
		title string
	)
	traceID := c.Param("trk_id")
	url := c.Request.URL.Path
	if strings.Contains(url, "edit") {
		title = "Редагування | "
	} else {
		title = ""
	}

	if err := trace.getBasicData(traceID); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}

	render(c, gin.H{
		"title": title + string(trace.Fields["reg_name"].([]uint8)),
	}, "index.html")

}
