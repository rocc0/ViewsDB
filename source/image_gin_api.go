package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	pb "./imager/imagegrpc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

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

func postAddImage(c *gin.Context) {
	var i newImage

	i.DocID = c.PostForm("doc_id")
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
		"photo_id": i.PhotoID,
	})
}

func postDelImage(c *gin.Context) {
	var d delImage
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	defer conn.Close()

	client := pb.NewImagerClient(conn)

	colID := c.Param("trk_id")

	x, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal([]byte(x), &d); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if err := removeImage(client, &pb.RemoveRequest{colID, d.PhotoID}); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "deleted",
	})

}
