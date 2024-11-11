package main

import (
	"context"
	"log"

	"github.com/rakeranjan/image-service/internal/configuration"
	"github.com/rakeranjan/image-service/internal/database/dynamodb"
	objectdatabase "github.com/rakeranjan/image-service/internal/database/object_database"
	"github.com/rakeranjan/image-service/internal/queue"
	"github.com/rakeranjan/image-service/utils"
)

func main() {
	ctx := context.TODO()
	conf, err := configuration.NewConfiguration()
	if err != nil {
		log.Panicln("error while loading config, err: ", err)
		return
	}
	module := NewModule(ctx, conf)
	if module == nil {
		log.Panicln("error while creating module , err: ", err)
		return
	}
	module.StartProcess()
}

func setUpInfra() {
	db := dynamodb.NewDataBase()
	db.CreateUserTable()
	db.CreateImageTable()
	objStorage := objectdatabase.NewObjectStorage()
	objStorage.CreateBucket(context.TODO(), utils.PROCESSED_BUCKET)
	objStorage.CreateBucket(context.TODO(), utils.PROCESSING_BUCKET)
	queue := queue.NewQueue()
	queue.CreateQueue(utils.IMAGE_PROCESSED_QUEUE)
	queue.CreateQueue(utils.IMAGE_PROCESSING_QUEUE)
}
