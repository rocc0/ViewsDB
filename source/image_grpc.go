package main

import (
	"io"
	"log"

	pb "../imager.d/pb"
	"golang.org/x/net/context"
)

func addImage(client pb.ImagerClient, img *pb.NewImageRequest) (*pb.NewImageResponse, error) {
	resp, err := client.AddImage(context.Background(), img)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func removeImage(client pb.ImagerClient, img *pb.RemoveRequest) error {
	resp, err := client.DeleteImage(context.Background(), img)
	if err != nil {
		return err
	}
	if resp.Success {
		log.Printf("Image has been removed")
	}
	return nil
}

// getImages calls the RPC method GetIage of ImagerServer
func getImagesGRPC(client pb.ImagerClient, filter *pb.ImagesFilter) ([]newImage, error) {
	// calling the streaming API
	images := []newImage{}

	stream, err := client.GetImages(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get images: %v", err)
	}
	for {
		image, err := stream.Recv()
		if err == io.EOF {
			return images, nil
		}
		if err != nil {
			log.Fatalf("%v.GetImages(_) = _, %v", client, err)
		}
		images = append(images, newImage{image.PhotoID, image.DocID, image.Thumb})
	}
	return images, nil
}
