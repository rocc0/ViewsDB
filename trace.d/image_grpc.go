package main

import (
	"io"
	"log"

	"mime/multipart"

	pb "github.com/rocc0/TraceDB/imager.d/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type (
	newImage struct {
		PhotoID string `json:"url"`
		DocID   string `json:"doc_id"`
		Thumb   string `json:"thumbUrl"`
	}
	delImage struct {
		DocID   string `json:"doc_id"`
		PhotoID string `json:"photo_id"`
	}
)

func sendAddRequest(client pb.ImagerClient, img *pb.NewImageRequest) (*pb.NewImageResponse, error) {
	resp, err := client.AddImage(context.Background(), img)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//Create image
func (i *newImage) createAddImage(f *multipart.FileHeader) error {
	var b2 []byte
	b := make([]byte, 100000)
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Print(err)
		}
	}()
	pbClient := pb.NewImagerClient(conn)

	file, err := f.Open()
	if err != nil {
		return err
	}

	for {
		if _, err := file.Read(b); err == io.EOF {
			log.Print("ok", len(b2))
			break
		}
		b2 = append(b2, b...)
	}
	img, err := sendAddRequest(pbClient, &pb.NewImageRequest{DocID: i.DocID, Photo: b2})

	i.PhotoID = img.PhotoID
	i.Thumb = img.Thumb
	return nil
}

// getImages calls the RPC method GetIage of ImagerServer
func readUrlsFromStream(client pb.ImagerClient, filter *pb.ImagesFilter) ([]newImage, error) {
	// calling the streaming API
	var images []newImage

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
			log.Printf("%v.GetImages(_) = _, %v", client, err)
			return nil, err
		}
		images = append(images, newImage{image.PhotoID, image.DocID, image.Thumb})
	}
	return images, nil
}

func getImageUrls(colID string) ([]newImage, error) {
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Print(err)
		}
	}()
	pbClient := pb.NewImagerClient(conn)

	images, err := readUrlsFromStream(pbClient, &pb.ImagesFilter{ColID: colID})
	if err != nil {
		return nil, err
	}

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
	defer func() {
		if err := conn.Close(); err != nil {
			log.Print(err)
		}
	}()

	client := pb.NewImagerClient(conn)
	if err := sendRemoveRequest(client, &pb.RemoveRequest{ColID: d.DocID, ImageID: d.PhotoID}); err != nil {
		return err
	}
	return nil
}
