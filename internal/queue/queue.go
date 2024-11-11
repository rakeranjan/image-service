package queue

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rakeranjan/image-service/internal/configuration"
	"github.com/rakeranjan/image-service/utils"
)

type Queue struct {
	client *sqs.Client
	conf   configuration.Config
}

var (
	instance *Queue
	once     sync.Once
)

func NewQueue() *Queue {
	configObj, _ := configuration.NewConfiguration()
	var err error
	once.Do(func() {
		config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(configObj.Region), config.WithBaseEndpoint(configObj.AwsBaseEndPoint))
		if err != nil {
			fmt.Println("error while loading config", err)
		}
		svc := sqs.NewFromConfig(config)
		instance = &Queue{
			client: svc,
			conf:   *configObj,
		}
	})
	if err != nil {
		return nil
	}
	return instance
}

func (q *Queue) GetClient() *sqs.Client {
	return q.client
}

func (q *Queue) CreateQueue(queueName string) *Queue {
	var queueUrl string
	queueAttributes := map[string]string{}
	queueAttributes["FifoQueue"] = "true"
	queue, err := q.client.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName:  aws.String(queueName),
		Attributes: queueAttributes,
	})
	if err != nil {
		log.Printf("Couldn't create queue %v. Here's why: %v\n", utils.IMAGE_PROCESSING_QUEUE, err)
	} else {
		queueUrl = *queue.QueueUrl
	}
	log.Println("queue created, url:", queueUrl)
	if err != nil {
		return nil
	}

	return instance
}
