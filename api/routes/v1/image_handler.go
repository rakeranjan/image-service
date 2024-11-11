package v1

import (
	"github.com/gin-gonic/gin"
)

type ImageHandlerV1 interface {
	Upload(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}
