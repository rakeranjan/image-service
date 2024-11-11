package v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rakeranjan/image-service/api/models"
)

type UserServiceV1 interface {
	Create(ctx context.Context, user *models.User) (*models.UserResponse, error)
	Login(ctx context.Context, user *models.User)
	Get(ctx context.Context, c *gin.Context)
	Update(ctx context.Context, c *gin.Context)
	Delete(ctx context.Context, c *gin.Context)
}
