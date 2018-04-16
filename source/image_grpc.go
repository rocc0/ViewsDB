package main

import (
	"io"
	"log"

	"mime/multipart"

	pb "../imager.d/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type newImage struct {
	PhotoID string `json:"url"`
	DocID   string `json:"doc_id"`
	Thumb   string `json:"thumbUrl"`
}

type delImage struct {
	DocID   string `json:"doc_id"`
	PhotoID string `json:"photo_id"`
}

func sendAddRequest(client pb.ImagerClient, img *pb.NewImageRequest) (*pb.NewImageResponse, error) {
	resp, err := client.AddImage(context.Background(), img)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//Create image
func (i *newImage) createAddImage(f *multipart.FileHeader) error {
	b := make([]byte, 100000)
	var b2 []byte
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	pbclient := pb.NewImagerClient(conn)

	file, err := f.Open()
	if err != nil {
		return err
	}

	for {
		_, err := file.Read(b)
		if err == io.EOF {
			log.Print("ok", len(b2))
			break
		}
		b2 = append(b2, b...)
	}
	img, err := sendAddRequest(pbclient, &pb.NewImageRequest{i.DocID, b2})

	i.PhotoID = img.PhotoID
	i.Thumb = img.Thumb
	return nil
}

// getImages calls the RPC method GetIage of ImagerServer
func readUrlsFromStream(client pb.ImagerClient, filter *pb.ImagesFilter) ([]newImage, error) {
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

func getImageUrls(docID string) ([]newImage, error) {
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	pbclient := pb.NewImagerClient(conn)

	images, err := readUrlsFromStream(pbclient, &pb.ImagesFilter{docID})

	return images, nil
}

func sendRemoveRequest(client pb.ImagerClient, img *pb.RemoveRequest) error {
	resp, err := client.DeleteImage(context.Background(), img)
	if err != nil {
		return err
	}
	if resp.Success {
		log.Printf("Image has been removed")
	}
	return nil
}

func (d delImage) createRemoveRequest() error {
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewImagerClient(conn)

	if err := sendRemoveRequest(client, &pb.RemoveRequest{d.DocID, d.PhotoID}); err != nil {
		return err
	}

	return nil
}
