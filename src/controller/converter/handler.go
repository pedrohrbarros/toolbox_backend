package converter

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	gooxml "baliance.com/gooxml/document"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog/log"
)

func ConvertFile(c *gin.Context) {
  type BindFile struct {
    File *multipart.FileHeader `form:"file" binding:"required"`
  }

	var bind_file BindFile

	available_types := []string{"pdf", "docx"}

	if err := c.ShouldBind(&bind_file); err != nil {
    error_message := fmt.Sprintf("Failed to process request: %s", err.Error())
		log.Error().Msg(error_message)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Error while binding request"})
		return
	}

  file := bind_file.File

  file_type := strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, ".")) - 1]

  if !slices.Contains(available_types, file_type) {
    error_message := fmt.Sprintf("Invalid file type: %s", file_type)
		log.Error().Msg(error_message)
    c.JSON(http.StatusBadRequest, gin.H{"Error": error_message})
    return
  }

  destination := filepath.Base(file.Filename)

  if err := c.SaveUploadedFile(file, destination); err != nil {
    error_message := fmt.Sprintf("Failed to upload file: %s", err.Error())
		log.Error().Msg(error_message)
    c.JSON(http.StatusBadRequest, gin.H{"Error": "Error while upload file"})
		return
  }

  if file_type == "docx" {
    if err := ConvertDocxToPDF(c, destination); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		return
  } else {
		error_message := "Not allowed type"
		log.Error().Msg(error_message)
		c.JSON(http.StatusNotAcceptable, gin.H{"Message": error_message})
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