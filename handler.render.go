package main

import (
	"strconv"
	"net/http"

	"github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context) {
	render(c, gin.H{
		"title":   "Пошук відстежень",
		}, "index.html")
}

func viewCreationPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Додати відстеження",
	}, "index.html")
}

func showRatings(c *gin.Context) {
	render(c, gin.H{
		"title": "Якісні показники",
	}, "index.html")
}

func getView(c *gin.Context) {
	viewID, err := strconv.Atoi(c.Param("trk_id"));
	if err == nil {
		trace, err := getBasicData(viewID)
		if err == nil {
			render(c, gin.H{
				"title": "Редагування | " + string(trace.Fields["requisits"].([]uint8)),
				}, "index.html")
		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		c.AbortWithStatus(http.StatusBadGateway)
	}
}


