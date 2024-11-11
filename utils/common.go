package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/rakeranjan/image-service/api/models"
)

func GetUSer(c *gin.Context) *models.User {
	data, ok := c.Get(USER)
	if !ok {
		return nil
	}
	user, ok := data.(*models.User)
	if !ok {
		return nil
	}
	return user
}
