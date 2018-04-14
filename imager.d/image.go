package main

import (
	"log"

	pb "./pb"
	minio "github.com/minio/minio-go"
)

type Image struct {
	PhotoID string `json:"url"`
	DocID   string `json:"doc_id"`
	Thumb   string `json:"thumbUrl"`
}

func (i Image) getImages() ([]*pb.Images, error) {
	var result []*pb.Images

	return result, nil
}

func (i Image) addImage() error {

	return nil
}

func (i Image) deleteImage() error {
	client, err := minio.NewV4(config.MinioUrl, config.MinioKay, config.MinioSecret, false)
	if err != nil {
		log.Print(err)
		return err
	}

	err = client.RemoveObject(i.DocID, i.PhotoID)
	if err != nil {
		log.Print(err)
		return err
	}
	err = client.RemoveObject(i.DocID, i.PhotoID[8:])
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
