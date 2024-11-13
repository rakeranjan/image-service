package v1

import (
	"context"
	"mime/multipart"

	"github.com/rakeranjan/image-service/api/models"
	"github.com/stretchr/testify/mock"
)

type MockImageRepository struct {
	mock.Mock
}

func (m *MockImageRepository) SaveImageMetaData(ctx context.Context, data *models.ImageMetaData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockImageRepository) UploadToProcessing(ctx context.Context, metaData *models.ImageMetaData, fileHeader *multipart.FileHeader) error {
	args := m.Called(ctx, metaData, fileHeader)
	return args.Error(0)
}
func (m *MockImageRepository) SendToSqsForProcessing(ctx context.Context, imageMetaData *models.ImageMetaData) error {
	args := m.Called(ctx, imageMetaData)
	return args.Error(0)
}
func (m *MockImageRepository) GetImageMetaData(ctx context.Context, user models.User, imageID string) (*models.ImageResponse, error) {
	args := m.Called(ctx, user, imageID)
	return args.Get(0).(*models.ImageResponse), args.Error(1)
}
func (m *MockImageRepository) GetImageMetaDataByImageID(ctx context.Context, user *models.User, imageID string) (*models.ImageResponse, error) {
	args := m.Called(ctx, user, imageID)
	return args.Get(0).(*models.ImageResponse), args.Error(1)
}
func (m *MockImageRepository) GetAllImageMetaData(ctx context.Context, user *models.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}
func (m *MockImageRepository) UpdateImageMetaData(ctx context.Context, data *models.ImageMetaData) error {
	args := m.Called(ctx, data)
	return args.Error(1)
}
func (m *MockImageRepository) DeleteImageMetaData(ctx context.Context, data *models.ImageMetaData) error {
	args := m.Called(ctx, data)
	return args.Error(1)
}
