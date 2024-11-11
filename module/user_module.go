package module

import (
	"context"

	"github.com/gin-gonic/gin"
	routesV1 "github.com/rakeranjan/image-service/api/routes/v1"
	"github.com/rakeranjan/image-service/internal/configuration"
	"github.com/rakeranjan/image-service/internal/database/dynamodb"
	"github.com/rakeranjan/image-service/utils"
)

type UserService struct {
	ModuleName string
	conf       *configuration.Config
	ctx        context.Context
}

func NewUserService(ctx context.Context, conf *configuration.Config) *UserService {
	return &UserService{
		ModuleName: utils.USER_SERVICE,
		conf:       conf,
		ctx:        ctx,
	}
}

func (u UserService) StartProcess() error {
	setUpUserServiceInfra()
	router := gin.Default()
	routesV1.Router(router, u.ModuleName)
	port := u.conf.Port
	router.Run(port)
	return nil
}

func setUpUserServiceInfra() {
	db := dynamodb.NewDataBase()
	db.CreateUserTable()
}
