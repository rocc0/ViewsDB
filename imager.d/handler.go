package main

import (
	"log"

	pb "./pb"
	"golang.org/x/net/context"
)

type server struct {
}

func (s server) GetImages(filter *pb.ImagesFilter, stream pb.Imager_GetImagesServer) error {
	var i = Image{DocID: filter.ColID}
	urls, err := i.getImages()
	if err != nil {
		return err
	}
	for _, image := range urls {
		stream.Send(image)
	}
	return nil
}

func (s server) DeleteImage(ctx context.Context, rq *pb.RemoveRequest) (*pb.RemoveResponse, error) {

	var i Image
	i.PhotoID = rq.ImageID
	i.DocID = rq.ColID
	log.Print(i.PhotoID, " | ", i.DocID)
	if err := i.deleteImage(); err != nil {
		return nil, err
	}

	return &pb.RemoveResponse{true}, nil
}

func (s server) AddImage(ctx context.Context, rq *pb.NewImageRequest) (*pb.NewImageResponse, error) {
	i := Image{rq.PhotoID, rq.DocID, rq.Thumb}

	if err := i.addImage(); err != nil {
		return nil, err
	}
	return &pb.NewImageResponse{true}, nil
}
