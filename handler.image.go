package main

import (
	"net/http"
	"os"
	"log"
	"io/ioutil"
	"encoding/json"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)


type delImage struct {
	Photo_id string
}


func getImages(c *gin.Context) {
	id := c.Param("trk_id")
	res, err := getImageUrls(id)
	check(err)

	c.JSON(http.StatusOK, gin.H{
		"images": res,
	})
}

func postImage(c *gin.Context) {

	file, err := c.FormFile("file")
	check(err)

	id := c.PostForm("docid")

	photo_id := randSeq(20)

	os.Mkdir("./static/images/"+id, 755)
	os.Mkdir("./static/images/"+id+"/resized", 755)

	original := "/static/images/"+id+"/"+photo_id + ".jpg"
	thumb := "/static/images/"+id+"/resized/"+photo_id + ".jpg"

	if err := c.SaveUploadedFile(file, "."+original); err != nil {
		log.Print(err.Error())
	}

	src, err := imaging.Open("."+original)
	dstImage128 := imaging.Resize(src, 128, 0, imaging.Lanczos)

	err = imaging.Save(dstImage128, "."+thumb)
	if err != nil {
		log.Printf("Save failed: %v", err)
	}

	addImageUrls(id, original, thumb, photo_id)

	c.JSON(http.StatusOK, gin.H{
		"original": original,
		"thumb": thumb,
		"photo_id": photo_id,
	})
}

func postDelImage(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var rq delImage
	err := json.Unmarshal([]byte(x), &rq)

	check(err)
	col := c.Param("trk_id")

	err = deleteImage(rq.Photo_id, col)
	check(err)

	original := "./static/images/"+col+"/"+ rq.Photo_id + ".jpg"
	thumb := "./static/images/"+col+"/resized/"+rq.Photo_id + ".jpg"

	err = os.Remove(original)
	check(err)

	err = os.Remove(thumb)
	check(err)

	c.JSON(http.StatusOK, gin.H{
		"msg": "deleted",
	})
}