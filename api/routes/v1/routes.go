package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	handlerv1 "github.com/rakeranjan/image-service/api/handler/v1"
	"github.com/rakeranjan/image-service/api/middleware/auth"
	"github.com/rakeranjan/image-service/api/middleware/validators"
	"github.com/rakeranjan/image-service/utils"
)

func Router(engine *gin.Engine, module string) {
	switch module {
	case utils.USER_SERVICE:
		{
			UserRouter(engine, handlerv1.NewUserHandlerImpl())
		}
	case utils.IMAGE_SERVICE:
		{
			ImageRouter(engine, handlerv1.NewImageHandlerImpl())
		}
	}

}

func ImageRouter(engine *gin.Engine, imageHandler ImageHandlerV1) {
	engine.Use(auth.Authoriser())
	engine.MaxMultipartMemory = 1 << 5
	v1 := engine.Group("v1")
	{
		v1.POST("/image", imageHandler.Upload)
		v1.GET("/image/:id", imageHandler.GetByID)
		v1.GET("/images", imageHandler.List)
		v1.PUT("/image/:id", imageHandler.Update)
		v1.DELETE("/image/:id", imageHandler.Delete)
	}
}

// Can be moved to user service, putting it here for now to save my time for assignment
func UserRouter(engine *gin.Engine, userHandler UserHandlerV1) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("userNameFormat", validators.ValidateUserName)
		v.RegisterValidation("userNameFormat", validators.ValidateUserName)
		v.RegisterValidation("phoneNumberFormat", validators.ValidatePhoneNumber)
	}
	// validator
	v1 := engine.Group("v1")
	{
		v1.POST("/user", userHandler.Create)
		v1.GET("/user/login", userHandler.Login)
		v1.GET("/user/:id", userHandler.Get)
		v1.PUT("/image/:id", userHandler.Update)
		v1.DELETE("/image/:id", userHandler.Delete)
	}
}
