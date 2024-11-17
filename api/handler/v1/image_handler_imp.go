package v1

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	apiutils "github.com/rakeranjan/image-service/api/api_utils"
	servicesV1 "github.com/rakeranjan/image-service/api/services/v1"
	"github.com/rakeranjan/image-service/utils"
)

type ImageHandlerImpl struct {
	imageService ImageServiceV1
}

func NewImageHandlerImpl() *ImageHandlerImpl {
	return &ImageHandlerImpl{
		imageService: servicesV1.NewImageServiceV1(),
	}
}

func (i *ImageHandlerImpl) Upload(c *gin.Context) {
	user := utils.GetUSer(c)
	file, _ := c.FormFile("file")
	ok := apiutils.ValidImageFile(file.Filename)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file extension: " + file.Filename})
		return
	}
	imageMetaData, err := i.imageService.Upload(c.Request.Context(), user, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while uploading file"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "uploaded successfully", "imageMetaData": imageMetaData})
}
func (i *ImageHandlerImpl) GetByID(c *gin.Context) {
	user := utils.GetUSer(c)
	ctx := c.Request.Context()
	id := c.Param("id")
	response, err := i.imageService.GetByID(ctx, user, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filePath := response.ImageId + "-" + response.FileName
	c.File(filePath)
	defer os.Remove(filePath)
}
func (i *ImageHandlerImpl) List(c *gin.Context) {
	user := utils.GetUSer(c)
	ctx := c.Request.Context()
	response, err := i.imageService.List(ctx, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.File(response)
	fileName := strings.Split(response, ".")[0]
	defer os.Remove(response)
	defer os.RemoveAll(utils.FILE_SUFFIX + fileName)
}
func (i *ImageHandlerImpl) Update(c *gin.Context) {
	user := utils.GetUSer(c)
	ctx := c.Request.Context()
	id := c.Param("id")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := i.imageService.Update(ctx, user, id, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "updated successfully", "imageMetaData": data})
}
func (i *ImageHandlerImpl) Delete(c *gin.Context) {
	user := utils.GetUSer(c)
	ctx := c.Request.Context()
	id := c.Param("id")
	ok := i.imageService.Delete(ctx, user, id)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "deleted successfully"})
}
