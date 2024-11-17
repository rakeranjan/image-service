package v1

import (
	"context"
	"mime/multipart"

	"github.com/rakeranjan/image-service/api/models"
)

type ImageServiceV1 interface {
	Upload(ctx context.Context, user *models.User, fileHeader *multipart.FileHeader) (*models.ImageMetaData, error)
	GetByID(ctx context.Context, user *models.User, imageID string) (*models.ImageResponse, error)
	List(ctx context.Context, user *models.User) (string, error)
	Update(ctx context.Context, user *models.User, imageID string, fileHeader *multipart.FileHeader) (*models.ImageMetaData, error)
	Delete(ctx context.Context, user *models.User, imageID string) bool
}
