package main

import (
	"bytes"
	"image/jpeg"

	pb "./pb"
	"github.com/minio/minio-go"
	"github.com/nfnt/resize"
	"github.com/rocc0/TraceDB/source/gen"
)

type Image struct {
	PhotoID string `json:"url"`
	DocID   string `json:"doc_id"`
	Thumb   string `json:"thumbUrl"`
}

func checkBucketExists(name string) error {
	client, err := minio.NewV4(config.MinioUrl, config.MinioKay, config.MinioSecret, false)
	if err != nil {
		return err
	}
	if ok, _ := client.BucketExists(name); !ok {
		if err = client.MakeBucket(name, "us-east-1"); err != nil {
			return err
		}
	}
	return nil
}

func uploadFilesToMinio(img pb.NewImageRequest) (*pb.NewImageResponse, error) {
	r := bytes.NewReader(img.Photo)
	photoID := gen.Generate(20) + ".jpg"
	var i = Image{PhotoID: photoID, DocID: img.DocID, Thumb: "resized/" + photoID}

	client, err := minio.NewV4(config.MinioUrl, config.MinioKay, config.MinioSecret, false)
	if err != nil {
		return nil, err
	}

	resized, err := resizeImage(img.Photo)
	if err != nil {
		return nil, err
	}

	if err = checkBucketExists(i.DocID); err != nil {
		return nil, err
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

func (i Image) getImages(filter *pb.ImagesFilter, stream pb.Imager_GetImagesServer) (<-chan minio.ObjectInfo, error) {
	var doneCh chan struct{}

	client, err := minio.NewV4(config.MinioUrl, config.MinioKay, config.MinioSecret, false)
	if err != nil {
		return nil, err
	}
	if err = checkBucketExists(filter.ColID); err != nil {
		return nil, err
	}
	urls := client.ListObjectsV2(i.DocID, "resized/", true, doneCh)

	return urls, nil
}

func (i Image) deleteImage() error {
	client, err := minio.NewV4(config.MinioUrl, config.MinioKay, config.MinioSecret, false)
	if err != nil {
		return err
	}

	err = client.RemoveObject(i.DocID, i.PhotoID)
	if err != nil {
		return err
	}
	err = client.RemoveObject(i.DocID, i.PhotoID[8:])
	if err != nil {
		return err
	}

	return nil
}
