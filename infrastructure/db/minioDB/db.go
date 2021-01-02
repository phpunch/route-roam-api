package miniodb

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/phpunch/route-roam-api/log"
	"mime/multipart"
)

type DB interface {
	CreateBucket(ctx context.Context, name string) error
	UploadFile(ctx context.Context, objectName string, file *multipart.FileHeader, contentType string) error
	GetFile(ctx context.Context, objectName string) (*minio.Object, error)
}

type db struct {
	Client *minio.Client
	config *Config
}

func New(config *Config) DB {
	// Initialize minio client object.
	minioClient, err := minio.New(config.Endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				config.AccessKeyID,
				config.SecretAccessKey,
				"",
			),
			Secure: config.UseSSL,
		})
	if err != nil {
		log.Log.Fatalf("%v", err)
	}
	return &db{
		Client: minioClient,
		config: config,
	}
}

func (d *db) CreateBucket(ctx context.Context, name string) error {
	err := d.Client.MakeBucket(ctx, name, minio.MakeBucketOptions{})
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
}

func (d *db) UploadFile(
	ctx context.Context,
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

	_, err = d.Client.PutObject(ctx, d.config.BucketName, objectName, src, file.Size, minio.PutObjectOptions{ContentType: contentType, UserMetadata: userMetaData})
	if err != nil {
		return err
	}
	return nil
}

func (d *db) GetFile(ctx context.Context, objectName string) (*minio.Object, error) {
	object, err := d.Client.GetObject(ctx, d.config.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}
