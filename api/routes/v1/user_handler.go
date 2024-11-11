package v1

import "github.com/gin-gonic/gin"

type UserHandlerV1 interface {
	Create(c *gin.Context)
	Login(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}
