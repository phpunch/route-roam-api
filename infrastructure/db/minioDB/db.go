package minioDB

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/phpunch/route-roam-api/log"
	"mime/multipart"
)

type DB interface {
	CreateBucket(name string) error
	UploadFile(bucketName string, objectName string, file *multipart.FileHeader, contentType string) error
	GetFile(bucketName, objectName string) (*minio.Object, error)
}

type db struct {
	Client *minio.Client
}

func New() DB {
	endpoint := "localhost:9000"
	accessKeyID := "route-roam"
	secretAccessKey := "route-roam"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Log.Fatalf("%v", err)
	}
	return &db{
		Client: minioClient,
	}
}

func (d *db) CreateBucket(name string) error {
	location := "us-east-1"

	ctx := context.Background()
	err := d.Client.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := d.Client.BucketExists(ctx, name)
		if errBucketExists == nil && exists {
			// return fmt.Errorf("We already own %s\n", name)
			log.Log.Infof("We already own %s\n", name)
			return nil
		}
		return fmt.Errorf("%v", err)
	}
	return nil

	// // Upload the zip file
	// objectName := "golden-oldies.zip"
	// filePath := "/tmp/golden-oldies.zip"
	// contentType := "application/zip"

	// // Upload the zip file with FPutObject
	// n, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
}

func (d *db) UploadFile(
	bucketName string,
	objectName string,
	file *multipart.FileHeader,
	contentType string,
) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	userMetaData := map[string]string{"x-amz-acl": "public-read"}

	// Use PutObject to upload a zip file
	ctx := context.Background()

	_, err = d.Client.PutObject(ctx, bucketName, objectName, src, file.Size, minio.PutObjectOptions{ContentType: contentType, UserMetadata: userMetaData})
	if err != nil {
		return err
	}
	return nil
}

func (d *db) GetFile(bucketName, objectName string) (*minio.Object, error) {
	object, err := d.Client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}
