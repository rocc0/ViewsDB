package main

import (
	"log"
	"os"

	"./src/gen"
	"github.com/disintegration/imaging"
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
	g := gen.SeqLength{20}
	i.PhotoID = g.Generate()
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

	c := session.DB("images").C("i" + col)
	err := c.Find(nil).All(&result)

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
