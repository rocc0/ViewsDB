package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"log"

	"./httputil"
	"github.com/gin-gonic/gin"
)

// postAddImage godoc
// @Summary Image uploading endpoint
// @Description putting new image to minio storage
// @Accept   multipart/form-data
// @Produce  json
// @Param doc_id path string true "doc id"
// @Success 200 {object} main.newImage
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/upload [post]
func postAddImage(c *gin.Context) {
	var i newImage

	i.DocID = c.PostForm("doc_id")
	file, err := c.FormFile("file")

	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	if err = i.createAddImage(file); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":      "http://192.168.99.100:9000/" + i.DocID + "/" + i.PhotoID,
		"doc_id":   i.DocID,
		"thumbUrl": i.Thumb,
	})
}

// getTraceImages godoc
// @Summary Getting all images of document
// @Description Getting all images of document from minio through gRPC
// @Accept   json
// @Produce  json
// @Param trk_id path string true "trace id"
// @Success 200 {string} string "ok"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /img/{trk_id} [get]
func getTraceImages(c *gin.Context) {
	id := c.Param("trk_id")
	urls, err := getImageUrls(id)

	if err != nil {
		httputil.NewError(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"images": urls,
	})
}

// postDelImage godoc
// @Summary Deleting image of a document
// @Description Deleting image of a document
// @Accept   json
// @Produce  json
// @Param trk_id path string true "trace id"
// @Success 200 {string} string "ok"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/{trk_id}/delete [post]
func postDelImage(c *gin.Context) {
	var d delImage
	x, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal([]byte(x), &d); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	d.DocID = c.Param("trk_id")
	log.Print(d.DocID, " | ", d.PhotoID)
	if err := d.createRemoveRequest(); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "deleted",
	})

}
