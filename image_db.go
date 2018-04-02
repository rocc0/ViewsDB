package main

import (
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/minio/minio-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session

type newImage struct {
	PhotoID  string `json:"photo_id"`
	DocID    string `json:"doc_id"`
	Original string `json:"original"`
	Thumb    string `json:"thumb"`
}

type delImage struct {
	PhotoID string `json:"photo_id"`
}

func mgoConnect() error {
	s, err := mgo.Dial(config.Mongo)
	if err != nil {
		return err
	}
	s.SetMode(mgo.Monotonic, true)
	session = s.Copy()
	return nil
}

func (i *newImage) addImages() error {
	i.PhotoID = generate(20)
	i.Original = config.ImagePath + i.DocID + "/" + i.PhotoID + ".jpg"
	i.Thumb = config.ImagePath + i.DocID + "/resized/" + i.PhotoID + ".jpg"

	os.Mkdir("."+config.ImagePath+i.DocID, 0700)
	os.Mkdir("."+config.ImagePath+i.DocID+"/resized", 0700)

	return nil
}

func (i newImage) resizeImage() error {
	src, err := imaging.Open("." + i.Original)
	if err != nil {
		return err
	}

	dstImage128 := imaging.Resize(src, 128, 0, imaging.Lanczos)

	err = imaging.Save(dstImage128, "."+i.Thumb)
	if err != nil {
		return err
	}

	return nil
}

func getImageUrls(col string) ([]newImage, error) {
	var result []newImage
	client, err := minio.NewV4("192.168.99.100:9000", config.MinioKay,
		config.MinioSecret, false)
	if err != nil {
		return nil, err
	}

	if err != nil {
		log.Print(err.Error(), "minio error")
	}
	bkts, err := client.ListBuckets()
	log.Print(bkts, len(bkts))

	c := session.DB("images").C("i" + col)
	err = c.Find(nil).All(&result)

	if err != nil {
		log.Print(err.Error(), "image not found")
	}

	return result, nil
}

func (i newImage) addImageUrls() error {
	c := session.DB("images").C("i" + i.DocID)
	err := c.Insert(&i)

	if err != nil {
		return err
	}

	return nil
}

func (d delImage) deleteImage(col string) error {

	c := session.DB("images").C("i" + col)
	err := c.Remove(bson.M{"photoid": d.PhotoID})

	if err != nil {
		return err
	}

	return nil
}
