package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/smithy-go"
	apiutils "github.com/rakeranjan/image-service/api/api_utils"
	"github.com/rakeranjan/image-service/api/models"
	"github.com/rakeranjan/image-service/internal/configuration"
	db "github.com/rakeranjan/image-service/internal/database/dynamodb"
	objectdatabase "github.com/rakeranjan/image-service/internal/database/object_database"
	"github.com/rakeranjan/image-service/internal/queue"
	"github.com/rakeranjan/image-service/utils"
)

type ImageRepositoryImpl struct {
	objectStorage *objectdatabase.ObjectStorage
	dbClient      *db.DataBase
	queueClient   *queue.Queue
	conf          *configuration.Config
}

func NewImageRepositoryImpl() *ImageRepositoryImpl {
	conf, _ := configuration.NewConfiguration()
	return &ImageRepositoryImpl{
		objectStorage: objectdatabase.NewObjectStorage(),
		dbClient:      db.NewDataBase(),
		queueClient:   queue.NewQueue(),
		conf:          conf,
	}
}

// func (i *ImageRepositoryImpl)
func (i *ImageRepositoryImpl) UploadToProcessing(ctx context.Context, metaData *models.ImageMetaData, fileHeader *multipart.FileHeader) error {
	fileName := fileHeader.Filename
	objectKey := metaData.GetObjectKey()
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
		return err
	} else {
		defer file.Close()
		_, err = i.objectStorage.GetClient().PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(utils.PROCESSING_BUCKET),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
				log.Printf("Error while uploading object to %s. The object is too large.", utils.PROCESSING_BUCKET)
			} else {
				log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
					fileName, utils.PROCESSING_BUCKET, objectKey, err)
			}
		} else {
			err = s3.NewObjectExistsWaiter(i.objectStorage.GetClient()).Wait(
				ctx, &s3.HeadObjectInput{Bucket: aws.String(utils.PROCESSING_BUCKET), Key: aws.String(objectKey)}, time.Minute)
			if err != nil {
				log.Printf("Failed attempt to wait for object %s to exist.\n", objectKey)
			}
		}
	}
	return err
}

func (i *ImageRepositoryImpl) SendToSqsForProcessing(ctx context.Context, imageMetaData *models.ImageMetaData) error {
	message, err := json.Marshal(imageMetaData)
	if err != nil {
		return nil
	}
	_, err = i.queueClient.GetClient().SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(string(message)),
		QueueUrl:    aws.String(i.conf.ProcessingQueueURL),
		// MessageGroupId: aws.String(imageMetaData.UserId),
		// MessageDeduplicationId: aws.String(imageMetaData.UserId),
	})
	return err
}

func (i *ImageRepositoryImpl) GetImageMetaData(ctx context.Context, user *models.User, imageID string) (*models.ImageResponse, error) {
	imageMetaData, err := i.GetImageMetaDataWIthID(ctx, imageID, user.ID)
	if err != nil {
		return nil, err
	}
	link, err := i.GetDownLoadLink(ctx, imageMetaData)
	if err != nil {
		return nil, err
	}
	return &models.ImageResponse{
		ImageMetaData: *imageMetaData,
		File:          link,
	}, nil
}

func (i *ImageRepositoryImpl) GetImageMetaDataByImageID(ctx context.Context, user *models.User, imageID string) (*models.ImageResponse, error) {
	imageMetaData, err := i.GetImageMetaDataWIthID(ctx, imageID, user.ID)
	if err != nil {
		return nil, err
	}
	link, err := i.GetDownLoadLink(ctx, imageMetaData)
	if err != nil {
		return nil, err
	}
	return &models.ImageResponse{
		ImageMetaData: *imageMetaData,
		File:          link,
	}, nil
}

func (i *ImageRepositoryImpl) GetAllImageMetaData(ctx context.Context, user *models.User) (string, error) {
	imageMetaData, err := i.GetAllImageMetaDataWIthID(ctx, user.ID)
	if err != nil {
		return "", err
	}
	filePath, err := i.GetAllDownLoadLink(ctx, imageMetaData)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func (i *ImageRepositoryImpl) UpdateImageMetaData(ctx context.Context, data *models.ImageMetaData) error {
	panic("implement me")
}

func (i *ImageRepositoryImpl) DeleteProcessedObjext(ctx context.Context, data *models.ImageMetaData) error {
	objectKey := data.GetObjectKey()
	_, err := i.objectStorage.GetClient().DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(utils.PROCESSED_BUCKET),
		Key:    aws.String(objectKey),
	})
	return err
}

func (i *ImageRepositoryImpl) SaveImageMetaData(ctx context.Context, data *models.ImageMetaData) error {
	items, err := attributevalue.MarshalMap(data)
	if err != nil {
		log.Println("error while marshling user, err:", err)
		return err
	}
	_, err = i.dbClient.GetDB().PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(utils.IMAGE_TABLE),
		Item:      items,
	})
	if err != nil {
		log.Println("error while saving user, err:", err)
		return err
	}
	return nil
}

func (i *ImageRepositoryImpl) GetImageMetaDataWIthID(ctx context.Context, imageId, userId string) (*models.ImageMetaData, error) {
	keyCondition := "userId = :userId"

	filterExpression := "imageId = :imageId"

	// Define ExpressionAttributeValues
	expressionAttributeValues := map[string]types.AttributeValue{
		":userId":  &types.AttributeValueMemberS{Value: userId},
		":imageId": &types.AttributeValueMemberS{Value: imageId},
	}

	// Query the DynamoDB table
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(utils.IMAGE_TABLE),
		KeyConditionExpression:    aws.String(keyCondition),
		FilterExpression:          aws.String(filterExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	// Execute the query
	result, err := i.dbClient.GetDB().Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query table %s: %w", utils.IMAGE_TABLE, err)
	}

	// Check if no items are returned
	if len(result.Items) == 0 {
		return nil, fmt.Errorf("no metadata found for UserId: %s", userId)
	}

	// Unmarshal the results into a slice of ImageMetaData
	var metadata []*models.ImageMetaData
	for _, item := range result.Items {
		var imgMeta models.ImageMetaData
		err := attributevalue.UnmarshalMap(item, &imgMeta)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		metadata = append(metadata, &imgMeta)
	}

	return metadata[0], nil
}

func (i *ImageRepositoryImpl) GetAllImageMetaDataWIthID(ctx context.Context, userId string) ([]*models.ImageMetaData, error) {
	keyCondition := "userId = :userId"

	// Define ExpressionAttributeValues
	expressionAttributeValues := map[string]types.AttributeValue{
		":userId": &types.AttributeValueMemberS{Value: userId},
	}

	// Query the DynamoDB table
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(utils.IMAGE_TABLE),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	// Execute the query
	result, err := i.dbClient.GetDB().Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query table %s: %w", utils.IMAGE_TABLE, err)
	}

	// Check if no items are returned
	if len(result.Items) == 0 {
		return nil, fmt.Errorf("no metadata found for UserId: %s", userId)
	}

	// Unmarshal the results into a slice of ImageMetaData
	var metadata []*models.ImageMetaData
	for _, item := range result.Items {
		var imgMeta models.ImageMetaData
		err := attributevalue.UnmarshalMap(item, &imgMeta)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		metadata = append(metadata, &imgMeta)
	}

	return metadata, nil
}

// ideally we should fetch from utils.PROCESSED_BUCKET as the image will be moved there once the image is processed
func (i *ImageRepositoryImpl) GetDownLoadLink(ctx context.Context, imageMetaData *models.ImageMetaData) (*os.File, error) {
	objectKey := imageMetaData.GetObjectKey()
	result, err := i.objectStorage.GetClient().GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(utils.PROCESSED_BUCKET),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	file, err := os.Create(imageMetaData.ImageId + "-" + imageMetaData.FileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", imageMetaData.FileName, err)
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

func (i *ImageRepositoryImpl) GetAllDownLoadLink(ctx context.Context, imageMetaDatas []*models.ImageMetaData) (string, error) {
	var filePath string
	for _, imageMetaData := range imageMetaDatas {
		objectKey := imageMetaData.GetObjectKey()
		result, err := i.objectStorage.GetClient().GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(utils.PROCESSED_BUCKET),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			return "", err
		}
		fiLeName := imageMetaData.FileName
		filePath, err = apiutils.WriteFile(imageMetaData.UserId, fiLeName, result.Body)
		if err != nil {
			return "", err
		}
	}
	filePath, err := apiutils.ZipFolder(imageMetaDatas[0].UserId, imageMetaDatas[0].UserId+".zip")
	if err != nil {
		return "", err
	}
	return filePath, nil
}
