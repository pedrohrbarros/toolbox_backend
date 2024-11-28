package converter

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	gooxml "baliance.com/gooxml/document"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	api_err "github.com/pedrohrbarros/toolbox_backend/src/middleware/error"
	"github.com/rs/zerolog/log"
)

// @Summary Document converter
// @Description Convert a word file into pdf
// @Tags Converter
// @Accept mpfd
// @Produce mpfd
// @Param file formData file true "File that will be converted"
// @Success 200 {file} File converted
// @Failure 400 {object} error.ApiError
// @Failure 500 {object} error.ApiError
// @Router /convert [post]
func ConvertFile(c *gin.Context) {
  type BindFile struct {
    File *multipart.FileHeader `form:"file" binding:"required"`
  }

	var bind_file BindFile

	if err := c.ShouldBind(&bind_file); err != nil {
		api_error := api_err.NewBadRequestError(fmt.Sprintf("Failed to process request=%s\n", err.Error()))
		log.Error().Msg(api_error.Message)
		c.JSON(api_error.Code, api_error)
		return
	}

  file := bind_file.File

  file_type := strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, ".")) - 1]

  destination := filepath.Base(file.Filename)

  if err := c.SaveUploadedFile(file, destination); err != nil {
		api_error := api_err.NewInternalServerError(fmt.Sprintf("Failed to upload file: %s", err.Error()))
		log.Error().Msg(api_error.Message)
    c.JSON(api_error.Code, api_error)
		return
  }

  if file_type == "docx" {
    if err := ConvertDocxToPDF(c, destination); err != nil {
			api_error := api_err.NewInternalServerError(fmt.Sprintf("Failed to convert file: %s", err.Error()))
			c.JSON(api_error.Code, api_error)
			return
		}
		return
  } else {
		api_error := api_err.NewBadRequestError("Not allowed type")
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