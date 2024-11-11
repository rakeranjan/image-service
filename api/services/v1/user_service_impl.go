package v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rakeranjan/image-service/api/models"
	"github.com/rakeranjan/image-service/api/repository"
	"github.com/rakeranjan/image-service/internal/configuration"
)

type UserServiceImplV1 struct {
	userRepository UserRepository
	conf           *configuration.Config
}

func NewUserServiceV1() *UserServiceImplV1 {
	conf, _ := configuration.NewConfiguration()
	return &UserServiceImplV1{
		userRepository: repository.NewUserRepository(),
		conf:           conf,
	}
}

func (i UserServiceImplV1) Create(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	user, err := i.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	token, err := user.Encrypt(i.conf.SecretValue)
	if err != nil {
		return nil, err
	}
	response := &models.UserResponse{
		User:        *user,
		AccessToken: token,
	}
	return response, nil
}
func (i UserServiceImplV1) Login(ctx context.Context, user *models.User) {}
func (i UserServiceImplV1) Get(ctx context.Context, c *gin.Context)      {}
func (i UserServiceImplV1) Update(ctx context.Context, c *gin.Context)   {}
func (i UserServiceImplV1) Delete(ctx context.Context, c *gin.Context)   {}
