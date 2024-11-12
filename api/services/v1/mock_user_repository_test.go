package v1

import (
	"context"

	"github.com/rakeranjan/image-service/api/models"
	"github.com/stretchr/testify/mock"
)

type MockUserReposity struct {
	mock.Mock
}

func (m *MockUserReposity) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*models.User), args.Error(1)
}
