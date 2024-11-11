package module

import (
	"context"

	"github.com/rakeranjan/image-service/internal/configuration"
	"github.com/rakeranjan/image-service/utils"
)

// NOTIFICATION_SERVICE
type NotificationService struct {
	ModuleName string
	conf       *configuration.Config
	ctx        context.Context
}

func NewNotificationService(ctx context.Context, conf *configuration.Config) *NotificationService {
	return &NotificationService{
		ModuleName: utils.NOTIFICATION_SERVICE,
		conf:       conf,
		ctx:        ctx,
	}
}

func (u NotificationService) StartProcess() error {
	setUpImageProcessorInfra()
	// start listening to queue
	// STEP 1 recieve imageMetadata from Processed Queue
	// STEP 2 using imageMetadata fetch image from  S3 Processed bucket
	// STEP 4 send mesage to user using phoneNumber present in imageMetadata with the s3 link
	return nil
}

func setUpNotificationInfra() {
	// Setup SNS
}
