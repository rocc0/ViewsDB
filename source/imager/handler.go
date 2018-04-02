package main

import (
	"os"

	pb "./imagegrpc"
	"golang.org/x/net/context"
	"log"
)

type server struct {
}

func (s server) GetImages(filter *pb.ImagesFilter, stream pb.Imager_GetImagesServer) error {
	urls, err := getImageUrls(filter.ColID)
	if err != nil {
		return err
	}
	for _, image := range urls {
		stream.Send(image)
	}
	return nil
}

func (s server) DeleteImage(ctx context.Context, rq *pb.RemoveRequest) (*pb.RemoveResponse, error) {

	var i newImage
	i.PhotoID = rq.ImageID
	i.DocID = rq.ColID
	original := "." + config.ImagePath + rq.ColID + "/" + rq.ImageID
	log.Print(original)
	if err := i.deleteImage(); err != nil {
		return nil, err
	}

	if err := os.RemoveAll(original); err != nil {
		return nil, err
	}

	return &pb.RemoveResponse{true}, nil
}

func (s server) AddImage(ctx context.Context, rq *pb.NewImageRequest) (*pb.NewImageResponse, error) {
	i := newImage{rq.PhotoID, rq.DocID, rq.Original, rq.Thumb}

	if err := i.resizeImage(); err != nil {
		return nil, err
	}

	if err := i.addImageUrls(); err != nil {
		return nil, err
	}
	return &pb.NewImageResponse{true}, nil
}
