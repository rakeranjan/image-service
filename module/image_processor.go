package module

import (
	"context"

	"github.com/rakeranjan/image-service/internal/configuration"
	db "github.com/rakeranjan/image-service/internal/database/dynamodb"
	objectdatabase "github.com/rakeranjan/image-service/internal/database/object_database"
	"github.com/rakeranjan/image-service/internal/queue"
	"github.com/rakeranjan/image-service/utils"
	"github.com/rakeranjan/image-service/workers"
)

// IMAGE_PROCESSOR

type ImageProcessor struct {
	ModuleName string
	conf       *configuration.Config
	ctx        context.Context
}

func NewImageProcessor(ctx context.Context, conf *configuration.Config) *ImageProcessor {
	return &ImageProcessor{
		ModuleName: utils.IMAGE_PROCESSOR,
		conf:       conf,
		ctx:        ctx,
	}
}

func (u ImageProcessor) StartProcess() error {
	setUpImageProcessorInfra()
	objectStorage := objectdatabase.NewObjectStorage()
	sqsClient := queue.NewQueue()
	dbClient := db.NewDataBase()
	worker := workers.NewImageProcessor(u.conf, objectStorage, sqsClient, dbClient)
	worker.Process(u.ctx)
	return nil
}

func setUpImageProcessorInfra() {
	objStorage := objectdatabase.NewObjectStorage()
	objStorage.CreateBucket(context.TODO(), utils.PROCESSED_BUCKET)
	queue := queue.NewQueue()
	queue.CreateQueue(utils.IMAGE_PROCESSING_QUEUE)
	queue.CreateQueue(utils.IMAGE_PROCESSED_QUEUE)
}
