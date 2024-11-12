package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakeranjan/image-service/api/models"
	servicesV1 "github.com/rakeranjan/image-service/api/services/v1"
)

type UserHandlerImpl struct {
	UserService UserServiceV1
}

func NewUserHandlerImpl() *UserHandlerImpl {
	return &UserHandlerImpl{
		UserService: servicesV1.NewUserServiceV1(),
	}
}

func (u *UserHandlerImpl) Create(c *gin.Context) {
	var json models.User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := u.UserService.Create(c.Request.Context(), &json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (u *UserHandlerImpl) Login(c *gin.Context) {
	panic("not implemented")
}

func (u *UserHandlerImpl) Get(c *gin.Context) {
	panic("not implemented")
}
func (u *UserHandlerImpl) Update(c *gin.Context) {
	panic("not implemented")
}
func (u *UserHandlerImpl) Delete(c *gin.Context) {
	panic("not implemented")
}
