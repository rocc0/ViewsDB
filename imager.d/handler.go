package main

import (
	pb "./pb"
	"golang.org/x/net/context"
)

type server struct {
}

func (s server) AddImage(ctx context.Context, rq *pb.NewImageRequest) (*pb.NewImageResponse, error) {
	img := pb.NewImageRequest{DocID: rq.DocID, Photo: rq.Photo}
	if res, err := uploadFilesToMinio(img); err != nil {
		return nil, err
	} else {
		return &pb.NewImageResponse{res.DocID, res.PhotoID, res.Thumb}, nil
	}
}

func (s server) GetImages(filter *pb.ImagesFilter, stream pb.Imager_GetImagesServer) error {
	imgDir := Image{DocID: filter.ColID}
	urls, err := imgDir.getImages(filter, stream)
	if err != nil {
		return err
	}

	for i := range urls {
		img := pb.Image{PhotoID: "http://192.168.99.100:9000/" + filter.ColID + "/" + i.Key[8:],
			DocID: filter.ColID, Thumb: i.Key}
		if err := stream.Send(&img); err != nil {
			return err
		}
	}
	return nil
}

func (s server) DeleteImage(ctx context.Context, rq *pb.RemoveRequest) (*pb.RemoveResponse, error) {
	img := Image{PhotoID: rq.ImageID, DocID: rq.ColID}

	if err := img.deleteImage(); err != nil {
		return nil, err
	}

	return &pb.RemoveResponse{Success: true}, nil
}
