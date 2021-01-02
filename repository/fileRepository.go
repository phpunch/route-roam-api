package repository

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type fileRepository interface {
	UploadFile(ctx context.Context, objectName string, file *multipart.FileHeader, contentType string) error
	GetFile(ctx context.Context, objectName string) (*minio.Object, error)
}

func (r *repository) UploadFile(
	ctx context.Context,
	objectName string,
	file *multipart.FileHeader,
	contentType string,
) error {
	return r.Ds.MinioDB.UploadFile(ctx, objectName, file, contentType)
}

func (r *repository) GetFile(ctx context.Context, objectName string) (*minio.Object, error) {
	return r.Ds.MinioDB.GetFile(ctx, objectName)
}
