package v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rakeranjan/image-service/api/models"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*models.UserResponse), args.Error(1)
}
func (m *MockUserService) Login(ctx context.Context, user *models.User) {
	panic("not implemented")
}
func (m *MockUserService) Get(ctx context.Context, c *gin.Context) {
	panic("not implemented")
}
func (m *MockUserService) Update(ctx context.Context, c *gin.Context) {
	panic("not implemented")
}
func (m *MockUserService) Delete(ctx context.Context, c *gin.Context) {
	panic("not implemented")
}
