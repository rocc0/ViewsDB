package main

import (
	"log"
	"math/rand"
	"time"

	"bytes"
	"image/jpeg"

	pb "./pb"

	"github.com/minio/minio-go"
	"github.com/nfnt/resize"
)

type Image struct {
	PhotoID string `json:"url"`
	DocID   string `json:"doc_id"`
	Thumb   string `json:"thumbUrl"`
}

func uploadFilesToMinio(img pb.NewImageRequest) (*pb.NewImageResponse, error) {
	r := bytes.NewReader(img.Photo)
	log.Print(len(img.Photo))
	photoID := generate(20) + ".jpg"
	var i = Image{PhotoID: photoID, DocID: img.DocID, Thumb: "resized/" + photoID}

	client, err := minio.NewV4(config.MinioUrl, config.MinioKay, config.MinioSecret, false)
	if err != nil {
		return nil, err
	}

	resized, err := resizeImage(img.Photo)
	if err != nil {
		return nil, err
	}

	if ok, _ := client.BucketExists(i.DocID); !ok {
		if err = client.MakeBucket(i.DocID, "us-east-1"); err != nil {
			return nil, err
		}
	}
	if err = client.SetBucketPolicy(i.DocID, "", "readonly"); err != nil {
		return nil, err
	}
	//Save original
	if _, err = client.PutObject(i.DocID, i.PhotoID, r, r.Size(),
		minio.PutObjectOptions{ContentType: "multipart/form-data"}); err != nil {
		return nil, err
	}
	//Save thumb
	if _, err = client.PutObject(i.DocID, "resized/"+i.PhotoID, resized, int64(resized.Len()),
		minio.PutObjectOptions{ContentType: "image/jpeg"}); err != nil {
		return nil, err
	}
	return &pb.NewImageResponse{i.DocID, i.PhotoID, i.Thumb}, nil
}

func resizeImage(file []byte) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	r := bytes.NewReader(file)
	img, err := jpeg.Decode(r)
	if err != nil {
		return nil, err
	}

	dstImage128 := resize.Resize(128, 0, img, resize.Lanczos3)
	err = jpeg.Encode(buf, dstImage128, nil)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

//Generate returns a random seq of symbols
func generate(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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
