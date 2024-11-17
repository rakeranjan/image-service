package repository

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
	"github.com/rakeranjan/image-service/api/models"
	"github.com/rakeranjan/image-service/internal/configuration"
	db "github.com/rakeranjan/image-service/internal/database/dynamodb"
	objectdatabase "github.com/rakeranjan/image-service/internal/database/object_database"
	"github.com/rakeranjan/image-service/utils"
)

type ImageRepository struct {
	conf          *configuration.Config
	objectStorage *objectdatabase.ObjectStorage
	dbClient      *db.DataBase
}

func NewImageRepository(conf *configuration.Config, objectStorage *objectdatabase.ObjectStorage, dbClient *db.DataBase) *ImageRepository {
	return &ImageRepository{
		conf:          conf,
		objectStorage: objectStorage,
		dbClient:      dbClient,
	}
}

func (i *ImageRepository) UpdateImageMetaData(ctx context.Context, data *models.ImageMetaData) error {
	items, err := attributevalue.MarshalMap(data)
	if err != nil {
		return err
	}
	_, err = i.dbClient.GetDB().PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(utils.IMAGE_TABLE),
		Item:      items,
	})
	if err != nil {
		return err
	}
	return nil
}

func (i *ImageRepository) GetImage(ctx context.Context, data *models.ImageMetaData) (*os.File, error) {
	objectKey := data.GetObjectKey()
	result, err := i.objectStorage.GetClient().GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(utils.PROCESSING_BUCKET),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	file, err := os.Create(data.ImageId + "-" + data.FileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", data.FileName, err)
		return nil, err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return file, nil
}

func (i *ImageRepository) MoveImageToProcessed(ctx context.Context, data *models.ImageMetaData, filePath string) error {
	err := i.UploadImageToProcessed(ctx, data, filePath)
	if err != nil {
		log.Printf("error while uploading file: %s to processed bucket, err: %s", filePath, err.Error())
		return err
	}
	err = i.DeleteFromProcessing(ctx, data)
	if err != nil {
		log.Printf("error while deleting from processing bucket, err: %s", err.Error())
	}
	return err
}

func (i *ImageRepository) DeleteFromProcessing(ctx context.Context, data *models.ImageMetaData) error {
	objectKey := data.GetObjectKey()
	_, err := i.objectStorage.GetClient().DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(utils.PROCESSING_BUCKET),
		Key:    aws.String(objectKey),
	})
	return err
}

func (i *ImageRepository) UploadImageToProcessed(ctx context.Context, metaData *models.ImageMetaData, filePath string) error {
	// fileName := fileHeader.Filename
	objectKey := metaData.GetObjectKey()
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", filePath, err)
		return err
	} else {
		defer file.Close()
		_, err = i.objectStorage.GetClient().PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(utils.PROCESSED_BUCKET),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
				log.Printf("Error while uploading object to %s. The object is too large.", utils.PROCESSED_BUCKET)
			} else {
				log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
					filePath, utils.PROCESSED_BUCKET, objectKey, err)
			}
		} else {
			err = s3.NewObjectExistsWaiter(i.objectStorage.GetClient()).Wait(
				ctx, &s3.HeadObjectInput{Bucket: aws.String(utils.PROCESSED_BUCKET), Key: aws.String(objectKey)}, time.Minute)
			if err != nil {
				log.Printf("Failed attempt to wait for object %s to exist.\n", objectKey)
			}
		}
	}
	return err
}
