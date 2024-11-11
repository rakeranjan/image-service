package v1

import (
	"context"
	"mime/multipart"

	"github.com/rakeranjan/image-service/api/models"
)

type ImageRepository interface {
	SaveImageMetaData(ctx context.Context, data *models.ImageMetaData) error
	UploadToProcessing(ctx context.Context, metaData *models.ImageMetaData, fileHeader *multipart.FileHeader) error
	SendToSqsForProcessing(ctx context.Context, imageMetaData *models.ImageMetaData) error
	GetImageMetaData(ctx context.Context, user models.User, imageID string) (*models.ImageResponse, error)
	GetImageMetaDataByImageID(ctx context.Context, user *models.User, imageID string) (*models.ImageResponse, error)
	GetAllImageMetaData(ctx context.Context, user *models.User) (string, error)
	UpdateImageMetaData(ctx context.Context, data *models.ImageMetaData) error
	DeleteImageMetaData(ctx context.Context, data *models.ImageMetaData) error
}
