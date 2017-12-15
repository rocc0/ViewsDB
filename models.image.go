package main

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Image struct {
	PhotoId string `json:"photo_id"`
	Original string `json:"original"`
	Thumb string `json:"thumb"`
}


func addImageUrls(col, original, thumb, phid string) error {
	session, err := mgo.Dial("mongodb://adder:password@192.168.99.100:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("images").C("i"+col)
	err = c.Insert(&Image{phid,original, thumb })
	if err != nil {
		log.Print(err.Error(),"error inserting data")
	}

	return nil
}

func getImageUrls(col string) ([]Image, error) {
	session, err := mgo.Dial("mongodb://adder:password@192.168.99.100:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var result[]Image

	c := session.DB("images").C("i"+col)

	err = c.Find(nil).All(&result)

	if err != nil {
		log.Print(err.Error(),"image not found")
	}

	return result,  nil
}

func deleteImage(photo_id, col string) error {

	session, err := mgo.Dial("mongodb://adder:password@192.168.99.100:27017")
	session.SetMode(mgo.Monotonic, true)

	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("images").C("i"+col)

	err = c.Remove(bson.M{"photoid": photo_id})

	check(err)

	return nil
}