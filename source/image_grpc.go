package main

import (
	"io"
	"log"

	pb "../imager.d/pb"
	"golang.org/x/net/context"
)

func addImage(client pb.ImagerClient, img *pb.NewImageRequest) {
	resp, err := client.AddImage(context.Background(), img)
	if err != nil {
		log.Fatalf("Could not add Image: %v", err)
	}
	if resp.Success {
		log.Printf("A new Image has been added")
	}
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
func getImagesGRPC(client pb.ImagerClient, filter *pb.ImagesFilter) {
	// calling the streaming API
	imgs := []*pb.Images{}
	stream, err := client.GetImages(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}
	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetCustomers(_) = _, %v", client, err)
		}
		imgs = append(imgs, customer)
	}

}
