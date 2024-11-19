package v1

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/rakeranjan/image-service/api/models"
	"github.com/rakeranjan/image-service/api/repository"
)

type ImageServiceImpl struct {
	imageRepository ImageRepository
}

func NewImageServiceV1() *ImageServiceImpl {
	return &ImageServiceImpl{
		imageRepository: repository.NewImageRepositoryImpl(),
	}
}

func (i ImageServiceImpl) Upload(ctx context.Context, user *models.User, fileHeader *multipart.FileHeader) (*models.ImageMetaData, error) {
	if user == nil {
		return nil, errors.New("user not found")
	}
	if fileHeader == nil {
		return nil, errors.New("no upload found")
	}
	imageMetaData := getImageAnalysis(user, fileHeader)
	err := i.imageRepository.SaveImageMetaData(ctx, imageMetaData)
	if err != nil {
		log.Println("Failed while saving imageMetadata to database, imageMetaData:", imageMetaData)
		return nil, err
	}

	err = i.imageRepository.UploadToProcessing(ctx, imageMetaData, fileHeader)
	if err != nil {
		log.Println("Failed while saving image to Object storage, imageMetaData:", imageMetaData)
		return nil, err
	}

	err = i.imageRepository.SendToSqsForProcessing(ctx, imageMetaData)
	if err != nil {
		log.Println("Failed while pushing message to Queue, imageMetaData:", imageMetaData)
		return nil, err
	}
	return imageMetaData, nil
}

func (i ImageServiceImpl) GetByID(ctx context.Context, user *models.User, imageID string) (*models.ImageResponse, error) {
	imageReponse, err := i.imageRepository.GetImageMetaDataByImageID(ctx, user, imageID)
	if err != nil {
		return nil, err
	}
	return imageReponse, nil
}

func (i ImageServiceImpl) List(ctx context.Context, user *models.User) (string, error) {
	imageReponses, err := i.imageRepository.GetAllImageMetaData(ctx, user)
	if err != nil {
		return "", err
	}
	return imageReponses, nil
}

func (i ImageServiceImpl) Update(ctx context.Context, user *models.User, imageID string, fileHeader *multipart.FileHeader) (*models.ImageMetaData, error) {
	data, err := i.imageRepository.GetImageMetaDataByImageID(ctx, user, imageID)
	if err != nil {
		return nil, err
	}
	// upload new image
	err = i.imageRepository.UploadToProcessing(ctx, &data.ImageMetaData, fileHeader)
	if err != nil {
		return nil, err
	}
	data.ImageMetaData.IsProcessed = false

	err = i.imageRepository.SendToSqsForProcessing(ctx, &data.ImageMetaData)
	if err != nil {
		return nil, err
	}

	err = i.imageRepository.DeleteProcessedObjext(ctx, &data.ImageMetaData)
	if err != nil {
		//  for now // return nil, err
	}
	return &data.ImageMetaData, err
}

func (i ImageServiceImpl) Delete(ctx context.Context, user *models.User, imageID string) bool {
	data, err := i.imageRepository.GetImageMetaDataByImageID(ctx, user, imageID)
	if err != nil {
		return false
	}
	err = i.imageRepository.DeleteProcessedObjext(ctx, &data.ImageMetaData)
	if err != nil {
		return false
	}
	return true
}

func getImageAnalysis(user *models.User, file *multipart.FileHeader) *models.ImageMetaData {
	imageId := uuid.NewString()
	return &models.ImageMetaData{
		FileName:        file.Filename,
		ImageId:         imageId,
		CreatedAt:       time.Now(),
		UserId:          user.ID,
		UserPhoneNumber: user.PhoneNumber,
		SizeInKb:        int(file.Size / 1024),
	}
}
