package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiutils "github.com/rakeranjan/image-service/api/api_utils"
	"github.com/rakeranjan/image-service/utils"
)

func Authoriser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(utils.AUTHORIZATION)
		user, err := apiutils.DecodeJWTToStruct(token)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not Authorized"})
			return
		}
		c.Set(utils.USER, user)
		c.Next()
	}
}
