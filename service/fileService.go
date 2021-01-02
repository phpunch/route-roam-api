package service

import (
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type fileService interface {
	UploadFile(objectName string, file *multipart.FileHeader, contentType string) error
	GetFile(objectName string) (*minio.Object, error)
}

func (s *service) UploadFile(
	objectName string,
	file *multipart.FileHeader,
	contentType string,
) error {
	return s.repository.UploadFile(objectName, file, contentType)
}

func (s *service) GetFile(
	objectName string,
) (*minio.Object, error) {
	return s.repository.GetFile(objectName)
}
