package image

import (
	"fmt"
	"image/jpeg"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/pedrohrbarros/toolbox_backend/src/middleware/error"
	"github.com/rs/zerolog/log"
)

// @Summary Image Editor
// @Description Edit an image based on the parameters in the request
// @Tags Image
// @Accept json
// @Produce json
// @Accept multipart/form-data
// @Param width query int false "Image width"
// @Param height query int false "Image height"
// @Param image formData file true "Image to edit"
// @Success 200 {file} File "Converted file"
// @Failure 400 {object} error.ApiError
// @Failure 500 {object} error.ApiError
// @Router /edit-image [post]
func ResizeImage(c *gin.Context) {
  type Params struct {
    Width            int  `form:"width"`
    Height           int  `form:"height"`
  }

  var query_params Params

  if err := c.ShouldBindQuery(&query_params); err != nil {
    api_error := error.NewBadRequestError(fmt.Sprintf("Invalid query parameters: %s", err.Error()))
    log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
    return
  }

  file, err := c.FormFile("image")
  if err != nil {
    api_error := error.NewBadRequestError(fmt.Sprintf("File upload error: %s", err.Error()))
    log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
    return
  }

  file_reader, err := file.Open()
  if err != nil {
    api_error := error.NewInternalServerError(fmt.Sprintf("Failed to open file: %s", err.Error()))
    log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
    return
  }
  
  image, err := jpeg.Decode(file_reader)
	if err != nil {
    api_error := error.NewInternalServerError(fmt.Sprintf("Failed to open image: %s", err.Error()))
    log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
    return
  }
  defer file_reader.Close()

  resized_image := resize.Resize(uint(query_params.Width), uint(query_params.Height), image, resize.Lanczos3)

  output, err := os.Create("src/assets/temp/resized_image.jpg")
	if err != nil {
		api_error := error.NewInternalServerError(fmt.Sprintf("Failed to create file: %s", err.Error()))
    log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
	}
	defer output.Close()

  jpeg.Encode(output, resized_image, nil)

  c.File("src/assets/temp/resized_image.jpg")
}
