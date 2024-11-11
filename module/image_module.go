package module

import (
	"context"

	"github.com/gin-gonic/gin"
	routesV1 "github.com/rakeranjan/image-service/api/routes/v1"
	"github.com/rakeranjan/image-service/internal/configuration"
	"github.com/rakeranjan/image-service/internal/database/dynamodb"
	objectdatabase "github.com/rakeranjan/image-service/internal/database/object_database"
	"github.com/rakeranjan/image-service/internal/queue"
	"github.com/rakeranjan/image-service/utils"
)

type ImageService struct {
	ModuleName string
	conf       *configuration.Config
	ctx        context.Context
}

func NewImageService(ctx context.Context, conf *configuration.Config) *ImageService {
	return &ImageService{
		ModuleName: utils.IMAGE_SERVICE,
		conf:       conf,
		ctx:        ctx,
	}
}

func (u ImageService) StartProcess() error {
	setUpImageServiceInfra()
	router := gin.Default()
	routesV1.Router(router, u.ModuleName)
	port := u.conf.Port
	router.Run(port)
	return nil
}

func setUpImageServiceInfra() {
	db := dynamodb.NewDataBase()
	// db.DeleteTable(context.TODO(), utils.IMAGE_TABLE)
	db.CreateImageTable()
	objStorage := objectdatabase.NewObjectStorage()
	objStorage.CreateBucket(context.TODO(), utils.PROCESSED_BUCKET)
	objStorage.CreateBucket(context.TODO(), utils.PROCESSING_BUCKET)
	queue := queue.NewQueue()
	queue.CreateQueue(utils.IMAGE_PROCESSING_QUEUE)
}
