package module

import (
	"context"

	"github.com/rakeranjan/image-service/internal/configuration"
	objectdatabase "github.com/rakeranjan/image-service/internal/database/object_database"
	"github.com/rakeranjan/image-service/internal/queue"
	"github.com/rakeranjan/image-service/utils"
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
	// start listening to queue
	// STEP 1 recieve imageMetadat from queue
	// STEP 2 using imageMetadat fetch image from S3
	// STEP 3 analyse image & store info to imageMetaData, mark inProcessed to true
	// STEP 4 send mesage imageMetaDat as a message to Processed Queue
	return nil
}

func setUpImageProcessorInfra() {
	objStorage := objectdatabase.NewObjectStorage()
	objStorage.CreateBucket(context.TODO(), utils.PROCESSED_BUCKET)
	queue := queue.NewQueue()
	queue.CreateQueue(utils.IMAGE_PROCESSED_QUEUE)
}
