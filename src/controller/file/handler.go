package file

import (
	"errors"
	"fmt"
	"image/jpeg"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	gooxml "baliance.com/gooxml/document"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/nfnt/resize"
	api_error_middleware "github.com/pedrohrbarros/toolbox_backend/src/middleware/error"
	"github.com/rs/zerolog/log"
)

// @Summary Document converter
// @Description Convert a word file into pdf
// @Tags File
// @Accept mpfd
// @Produce mpfd
// @Param file formData file true "File that will be converted"
// @Param expected_type query string true "Expected type that will be converted"
// @Success 200 {file} File converted
// @Failure 400 {object} error.ApiError
// @Failure 500 {object} error.ApiError
// @Router /file/converter [post]
func ConvertFile(c *gin.Context) {
  type BindFile struct {
    File *multipart.FileHeader `form:"file" binding:"required"`
  }

	type Params struct {
		ExpectedType string `form:"expected_type"`
	}

	var bind_file BindFile
	var query_params Params 

	if err := c.ShouldBind(&bind_file); err != nil {
		api_error := api_error_middleware.NewBadRequestError(fmt.Sprintf("Failed to process request=%s\n", err.Error()))
		log.Error().Msg(api_error.Message)
		c.JSON(api_error.Code, api_error)
		return
	}

	if err := c.ShouldBindQuery(&query_params); err != nil {
    api_error := api_error_middleware.NewBadRequestError(fmt.Sprintf("Invalid query parameters: %s", err.Error()))
    log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
    return
  }

  file := bind_file.File

  file_type := strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, ".")) - 1]

  destination := fmt.Sprintf("src/assets/temp/%s", filepath.Base(file.Filename))

  if err := c.SaveUploadedFile(file, destination); err != nil {
		api_error := api_error_middleware.NewInternalServerError(fmt.Sprintf("Failed to upload file: %s", err.Error()))
		log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
		return
  }

  if file_type == "docx" && query_params.ExpectedType == "pdf" {
    if err := ConvertDocxToPDF(c, destination); err != nil {
			api_error := api_error_middleware.NewInternalServerError(fmt.Sprintf("Failed to convert file: %s", err.Error()))
			c.JSON(api_error.Code, api_error)
			return
		}
		return
  } else {
		os.Remove(destination)
		api_error := api_error_middleware.NewBadRequestError("Still not allowed type")
		log.Error().Msg(api_error.Message)
		c.JSON(api_error.Code, api_error)
		return
	}
}

func ConvertDocxToPDF(c *gin.Context, destination string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	doc, err := gooxml.Open(destination)
	if err != nil {
		error_message := fmt.Sprintf("Error opening .docx file: %s", err.Error())
		log.Error().Msg(error_message)
		return errors.New("error while opening .docx file")
	}

	for _, paragraph := range doc.Paragraphs() {
		var paragraph_text string
		for _, run := range paragraph.Runs() {
			paragraph_text += run.Text()
		}
		pdf.MultiCell(0, 10, paragraph_text, "", "", false)
	}

	output_file := "output.pdf"

	err = pdf.OutputFileAndClose(output_file)
	if err != nil {
		error_message := fmt.Sprintf("Error outputing file: %s", err.Error())
		log.Error().Msg(error_message)
		return errors.New("error while outputing file")
	}
	c.File(output_file)
	
	if err := os.Remove(destination); err != nil {
		error_message := fmt.Sprintf("Error removing initial file: %s", err.Error())
		log.Error().Msg(error_message)
		return errors.New("error while removing initial file")
	}
	if err := os.Remove(output_file); err != nil {
		error_message := fmt.Sprintf("Error removing output file: %s", err.Error())
		log.Error().Msg(error_message)
		return errors.New("error while removing output file")
	}

	return nil
}

// @Summary Image Editor
// @Description Edit an image based on the parameters in the request
// @Tags File
// @Accept multipart/form-data
// @Produce image/jpeg
// @Param image formData file true "JPEG image file to edit (only accepts .jpeg files)"
// @Param width query int false "Image width"
// @Param height query int false "Image height"
// @Success 200 {file} File "Converted file"
// @Failure 400 {object} error.ApiError
// @Failure 500 {object} error.ApiError
// @Router /file/image/resizer [post]
func ResizeImage(c *gin.Context) {
  type Params struct {
    Width            int  `form:"width"`
    Height           int  `form:"height"`
  }

  var query_params Params

  if err := c.ShouldBindQuery(&query_params); err != nil {
    api_error := api_error_middleware.NewBadRequestError(fmt.Sprintf("Invalid query parameters: %s", err.Error()))
    log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
    return
  }

  file, err := c.FormFile("image")
  if err != nil {
    api_error := api_error_middleware.NewBadRequestError(fmt.Sprintf("File upload error: %s", err.Error()))
    log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
    return
  }

  file_reader, err := file.Open()
  if err != nil {
    api_error := api_error_middleware.NewInternalServerError(fmt.Sprintf("Failed to open file: %s", err.Error()))
    log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
    return
  }
  
  image, err := jpeg.Decode(file_reader)
	if err != nil {
    api_error := api_error_middleware.NewBadRequestError(fmt.Sprintf("Only JPEG image types are allowed: %s", err.Error()))
    log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
    return
  }
  defer file_reader.Close()

  resized_image := resize.Resize(uint(query_params.Width), uint(query_params.Height), image, resize.Lanczos3)

  output, err := os.Create("src/assets/temp/resized_image.jpg")
	if err != nil {
		api_error := api_error_middleware.NewInternalServerError(fmt.Sprintf("Failed to create file: %s", err.Error()))
    log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
	}
	defer output.Close()

  jpeg.Encode(output, resized_image, nil)

  c.File("src/assets/temp/resized_image.jpg")
}