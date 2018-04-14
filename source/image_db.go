package main

import (
	"bytes"
	"image/jpeg"
	"mime/multipart"

	"github.com/minio/minio-go"
	"github.com/nfnt/resize"

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

//func mgoConnect() error {
//	s, err := mgo.Dial(config.Mongo)
//	if err != nil {
//		return err
//	}
//	s.SetMode(mgo.Monotonic, true)
//	session = s.Copy()
//	return nil
//}

//Create image
func (i *newImage) uploadFilesToMinio(file *multipart.FileHeader) error {
	i.PhotoID = generate(20) + ".jpg"
	i.Thumb = "resized/" + i.PhotoID
	client, err := minio.NewV4(config.MinioUrl, config.MinioKay, config.MinioSecret, false)
	if err != nil {
		return err
	}
	original, err := file.Open()
	if err != nil {
		return err
	}
	resized, err := resizeImage(file)
	if err != nil {
		return err
	}

	if ok, _ := client.BucketExists(i.DocID); !ok {
		if err = client.MakeBucket(i.DocID, "us-east-1"); err != nil {
			return err
		}
	}
	if err = client.SetBucketPolicy(i.DocID, "", "readonly"); err != nil {
		return err
	}
	//Save original
	if _, err = client.PutObject(i.DocID, i.PhotoID, original, file.Size,
		minio.PutObjectOptions{ContentType: "multipart/form-data"}); err != nil {
		return err
	}
	//Save thumb
	if _, err = client.PutObject(i.DocID, "resized/"+i.PhotoID, resized, int64(resized.Len()),
		minio.PutObjectOptions{ContentType: "image/jpeg"}); err != nil {
		return err
	}

	return nil
}

func resizeImage(file *multipart.FileHeader) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	img, err := jpeg.Decode(src)
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

func getImageUrls(col string) ([]newImage, error) {
	var doneCh chan struct{}
	var result []newImage
	client, err := minio.NewV4(config.MinioUrl, config.MinioKay, config.MinioSecret, false)
	if err != nil {
		return nil, err
	}
	list := client.ListObjectsV2(col, "resized/", true, doneCh)
	for i := range list {
		result = append(result, newImage{PhotoID: "http://192.168.99.100:9000/" + col + "/" + i.Key[8:], DocID: col, Thumb: i.Key})
	}
	return result, nil
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
