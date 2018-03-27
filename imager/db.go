package main

import (
	"log"

	"github.com/disintegration/imaging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	pb "./imagegrpc"
)

var session *mgo.Session

type newImage struct {
	PhotoID  string `json:"photo_id"`
	DocID    string `json:"doc_id"`
	Original string `json:"original"`
	Thumb    string `json:"thumb"`
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

func getImageUrls(col string) ([]*pb.Images, error) {
	var result []*pb.Images

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

func (i newImage) deleteImage() error {

	c := session.DB("images").C("i" + i.DocID)
	err := c.Remove(bson.M{"photoid": i.PhotoID})

	if err != nil {
		return err
	}

	return nil
}
