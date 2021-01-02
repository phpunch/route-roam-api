package repository

import (
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type fileRepository interface {
	UploadFile(objectName string, file *multipart.FileHeader, contentType string) error
	GetFile(objectName string) (*minio.Object, error)
}

func (r *repository) UploadFile(
	objectName string,
	file *multipart.FileHeader,
	contentType string,
) error {
	return r.Ds.MinioDB.UploadFile("image", objectName, file, contentType)
}

func (r *repository) GetFile(objectName string) (*minio.Object, error) {
	return r.Ds.MinioDB.GetFile("image", objectName)
}
