package service

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type fileService interface {
	UploadFile(ctx context.Context, objectName string, file *multipart.FileHeader, contentType string) (string, error)
	GetFile(ctx context.Context, objectName string) (*minio.Object, error)
}

func (s *service) UploadFile(
	ctx context.Context,
	objectName string,
	file *multipart.FileHeader,
	contentType string,
) (string, error) {
	return s.repository.UploadFile(ctx, objectName, file, contentType)
}

func (s *service) GetFile(
	ctx context.Context,
	objectName string,
) (*minio.Object, error) {
	return s.repository.GetFile(ctx, objectName)
}
