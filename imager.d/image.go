package main

import (
	"log"

	pb "./pb"
	"github.com/minio/minio-go"
)

type Image struct {
	PhotoID string `json:"url"`
	DocID   string `json:"doc_id"`
	Thumb   string `json:"thumbUrl"`
}

func (i Image) getImages(filter *pb.ImagesFilter, stream pb.Imager_GetImagesServer) (<-chan minio.ObjectInfo, error) {
	var doneCh chan struct{}

	client, err := minio.NewV4(config.MinioUrl, config.MinioKay, config.MinioSecret, false)
	if err != nil {
		return nil, err
	}
	urls := client.ListObjectsV2(i.DocID, "resized/", true, doneCh)

	return urls, nil
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
