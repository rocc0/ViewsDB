package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Image struct {
	PhotoId  string `json:"photo_id"`
	DocId    string `json:"doc_id"`
	Original string `json:"original"`
	Thumb    string `json:"thumb"`
}

type DelImage struct {
	PhotoId string `json:"photo_id"`
}

func getImages(c *gin.Context) {
	id := c.Param("trk_id")
	urls, err := getImageUrls(id)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	c.JSON(http.StatusOK, gin.H{
		"images": urls,
	})
}

func postImage(c *gin.Context) {
	var i Image
	i.DocId = c.PostForm("doc_id")
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if err := i.addImages(); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if err := c.SaveUploadedFile(file, "."+i.Original); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if err := i.resizeImage(); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if err := i.addImageUrls(); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	c.JSON(http.StatusOK, gin.H{
		"original": i.Original,
		"thumb":    i.Thumb,
		"photo_id": i.PhotoId,
	})
}

func postDelImage(c *gin.Context) {

	var d DelImage

	col := c.Param("trk_id")
	original := "." + imgpath + col

	x, _ := ioutil.ReadAll(c.Request.Body)

	if err := json.Unmarshal([]byte(x), &d); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if err := d.deleteImage(col); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if err := os.RemoveAll(original); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "deleted",
	})

}
