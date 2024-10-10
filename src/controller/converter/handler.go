package converter

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"baliance.com/gooxml/document"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func ConvertFile(c *gin.Context) {
  type BindFile struct {
    File *multipart.FileHeader `form:"file" binding:"required"`
  }

	var bind_file BindFile

	available_types := []string{"pdf", "docx"}

	if err := c.ShouldBind(&bind_file); err != nil {
    error_message := fmt.Sprintf("Failed to process request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Error": error_message})
		return
	}

  file := bind_file.File

  file_type := strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, ".")) - 1]

  if !slices.Contains(available_types, file_type) {
    error_message := fmt.Sprintf("Invalid file type: %s", file_type)
    c.JSON(http.StatusBadRequest, gin.H{"Error": error_message})
    return
  }

  destination := filepath.Base(file.Filename)

  if err := c.SaveUploadedFile(file, destination); err != nil {
    error_message := fmt.Sprintf("Failed to upload file: %s", err.Error())
    c.JSON(http.StatusBadRequest, gin.H{"Error": error_message})
		return
  }

  if file_type == "docx" {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "", 12)

		doc, err := document.Open(destination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error opening .docx file"})
			return
		}

		for _, paragraph := range doc.Paragraphs() {
			var paragraph_text string
      for _, run := range paragraph.Runs() {
        paragraph_text += run.Text()
      }
      pdf.MultiCell(0, 10, paragraph_text, "", "", false)
		}

    output_file := "output.pdf"

		err = pdf.OutputFileAndClose("output.pdf")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error outputing file"})
      return
    }
    
    c.File(output_file)
    os.Remove(destination)
    os.Remove(output_file)

    return
  }

	c.JSON(http.StatusOK, gin.H{"Message": "File converted successfully"})
}