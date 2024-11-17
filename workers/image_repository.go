package workers

import (
	"context"
	"os"

	"github.com/rakeranjan/image-service/api/models"
)

type ImageRepository interface {
	UpdateImageMetaData(ctx context.Context, data *models.ImageMetaData) error
	GetImage(ctx context.Context, data *models.ImageMetaData) (*os.File, error)
	MoveImageToProcessed(ctx context.Context, data *models.ImageMetaData, filePath string) error
}
