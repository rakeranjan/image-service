package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rakeranjan/image-service/api/models"
	"github.com/rakeranjan/image-service/internal/configuration"
	db "github.com/rakeranjan/image-service/internal/database/dynamodb"
	objectdatabase "github.com/rakeranjan/image-service/internal/database/object_database"
	"github.com/rakeranjan/image-service/internal/queue"
	"github.com/rakeranjan/image-service/workers/repository"
)

type ImageProcessor struct {
	conf            *configuration.Config
	sqsQueue        *queue.Queue
	channel         chan []types.Message
	wg              sync.WaitGroup
	imageRepository ImageRepository
	dbClient        *db.DataBase
	s3Client        *objectdatabase.ObjectStorage
}

func NewImageProcessor(conf *configuration.Config, objectStorage *objectdatabase.ObjectStorage, client *queue.Queue, dbClient *db.DataBase) *ImageProcessor {
	if conf == nil {
		conf, _ = configuration.NewConfiguration()
	}
	channel := make(chan []types.Message, 100)
	imageRepository := repository.NewImageRepository(conf, objectStorage, dbClient)
	return &ImageProcessor{
		conf:            conf,
		sqsQueue:        client,
		channel:         channel,
		imageRepository: imageRepository,
	}
}

func (i *ImageProcessor) Process(ctx context.Context) {
	i.wg.Add(2)
	go func() {
		defer i.wg.Done()
		if err := i.FetchMessages(ctx); err != nil {
			log.Printf("FetchMessages terminated with error: %v", err)
		}
	}()
	go func() {
		defer i.wg.Done()
		if err := i.ProcessMessages(ctx); err != nil {
			log.Printf("ProcessMessages terminated with error: %v", err)
		}
	}()
	i.wg.Wait()
	log.Println("All workers finished")
}

func (i *ImageProcessor) ProcessMessages(ctx context.Context) error {
	defer i.wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("ProcessMessages: context canceled")
			return ctx.Err()
		case messages := <-i.channel:
			for _, msg := range messages {
				log.Printf("Processing message: %v\n", *msg.Body)
				if err := i.handleMessage(ctx, &msg); err != nil {
					log.Printf("Failed to process message: %v\n", err)
					continue
				}
				_, err := i.sqsQueue.GetClient().DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(i.conf.ProcessingQueueURL),
					ReceiptHandle: msg.ReceiptHandle,
				})
				if err != nil {
					log.Printf("Failed to delete message: %v\n", err)
				} else {
					log.Printf("Message deleted successfully: %v\n", *msg.MessageId)
				}
			}
		}
	}
}

func (i *ImageProcessor) FetchMessages(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("FetchMessages: context canceled")
			return ctx.Err()
		default:
			result, err := i.sqsQueue.GetClient().ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(i.conf.ProcessingQueueURL),
				MaxNumberOfMessages: 5,
				WaitTimeSeconds:     1,
			})
			if err != nil {
				log.Printf("Couldn't get messages from queue %v. Here's why: %v\n", i.conf.ProcessingQueueURL, err)
				return err
			}

			if len(result.Messages) > 0 {
				log.Printf("Fetched %d messages", len(result.Messages))
				i.channel <- result.Messages
			}
		}
	}
}

func (i *ImageProcessor) handleMessage(ctx context.Context, msg *types.Message) error {
	data := &models.ImageMetaData{}
	err := json.Unmarshal([]byte(*msg.Body), data)
	if err != nil {
		fmt.Printf("error wile parsing Data: %s, err: %s \n", *msg.Body, err.Error())
		return err
	}
	_, err = i.imageRepository.GetImage(ctx, data)
	if err != nil {
		fmt.Printf("error wile fetchingImage err: %s \n", err.Error())
		return err
	}
	filePath := data.ImageId + "-" + data.FileName
	height, width := decodeImage(filePath)
	defer os.Remove("./" + filePath)
	data.IsProcessed = true
	data.Height = height
	data.Width = width
	err = i.imageRepository.UpdateImageMetaData(ctx, data)
	if err != nil {
		fmt.Printf("error wile updating imageMetadata err: %s \n", err.Error())
		return err
	}
	err = i.imageRepository.MoveImageToProcessed(ctx, data, filePath)
	if err != nil {
		fmt.Printf("error wile moving image err: %s \n", err.Error())
		return err
	}
	log.Println("process completed")
	return nil
}

func decodeImage(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	defer file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file) // Image Struct
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Height, image.Width

}
