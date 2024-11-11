package v1

import (
	"context"

	"github.com/rakeranjan/image-service/api/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
}
