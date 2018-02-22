package main

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/disintegration/imaging"
)


func (i *Image) addImages() error {

	i.PhotoId = RandSeq(20)
	i.Original = config.ImagePath + i.DocId + "/" + i.PhotoId + ".jpg"
	i.Thumb = config.ImagePath + i.DocId + "/resized/" + i.PhotoId + ".jpg"

	os.Mkdir("."+config.ImagePath+i.DocId, 755)
	os.Mkdir("."+config.ImagePath + i.DocId + "/resized", 755)

	return nil
}

func (i *Image) resizeImage() error {

	src, err := imaging.Open("." + i.Original)
	if err != nil {
		return err
	}

	dstImage128 := imaging.Resize(src, 128, 0, imaging.Lanczos)

	err = imaging.Save(dstImage128, "." + i.Thumb)
	if err != nil {
		return err
	}

	return nil
}

func (i *Image) addImageUrls() error {
	session, err := mgo.Dial(config.Mongo)
	if err != nil {
		return err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("images").C("i" + i.DocId)
	err = c.Insert(&i)
	if err != nil {
		return err
	}

	return nil
}

func getImageUrls(col string) ([]Image, error) {
	session, err := mgo.Dial(config.Mongo)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var result []Image

	c := session.DB("images").C("i" + col)
	err = c.Find(nil).All(&result)

	if err != nil {
		log.Print(err.Error(), "image not found")
	}

	return result, nil
}

func (d DelImage) deleteImage(col string) error {

	session, err := mgo.Dial(config.Mongo)
	session.SetMode(mgo.Monotonic, true)

	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("images").C("i" + col)

	err = c.Remove(bson.M{"photoid": d.PhotoId})

	if err != nil {
		return err
	}

	return nil
}
