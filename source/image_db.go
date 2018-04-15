package main

import (
	"mime/multipart"

	"log"

	pb "../imager.d/pb"
	"google.golang.org/grpc"
)

//var session *mgo.Session

type newImage struct {
	PhotoID string `json:"url"`
	DocID   string `json:"doc_id"`
	Thumb   string `json:"thumbUrl"`
}

type delImage struct {
	DocID   string `json:"doc_id"`
	PhotoID string `json:"photo_id"`
}

//Create image
func (i *newImage) uploadFilesToMinio(f *multipart.FileHeader) error {
	var b []byte
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
	_, err = file.Read(b)
	if err != nil {
		return err
	}

	img, err := addImage(pbclient, &pb.NewImageRequest{i.DocID, b})

	i.PhotoID = img.PhotoID
	i.Thumb = img.Thumb
	return nil
}

func getImageUrls(col string) ([]newImage, error) {
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	pbclient := pb.NewImagerClient(conn)

	images, err := getAllImages(pbclient, &pb.ImagesFilter{col})

	return images, nil
}

func (d delImage) deleteImage() error {
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewImagerClient(conn)

	if err := removeImage(client, &pb.RemoveRequest{d.DocID, d.PhotoID}); err != nil {
		return err
	}

	return nil
}
