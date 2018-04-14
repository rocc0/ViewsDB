package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
)

func postAddImage(c *gin.Context) {
	var i newImage

	i.DocID = c.PostForm("doc_id")
	file, err := c.FormFile("file")

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if err = i.uploadFilesToMinio(file); err != nil {
		log.Print(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	c.JSON(http.StatusOK, gin.H{
		"url":      "http://192.168.99.100:9000/" + i.DocID + "/" + i.PhotoID,
		"doc_id":   i.DocID,
		"thumbUrl": i.Thumb,
	})
}

func getTraceImages(c *gin.Context) {
	id := c.Param("trk_id")
	urls, err := getImageUrls(id)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	c.JSON(http.StatusOK, gin.H{
		"images": urls,
	})
}

func postDelImage(c *gin.Context) {
	var d delImage
	x, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal([]byte(x), &d); err != nil {
		log.Print(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	d.DocID = c.Param("trk_id")
	log.Print(d.DocID, " | ", d.PhotoID)
	if err := d.deleteImage(); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "deleted",
	})

}
