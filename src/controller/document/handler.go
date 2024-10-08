package document

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"

	"github.com/gin-gonic/gin"
)

func ConvertFile(c *gin.Context) {
  type BindFile struct {
    Name string `form:"name" binding:"required"`
    File *multipart.FileHeader `form:"file" binding:"required"`
  }

	var bind_file BindFile

	if err := c.ShouldBind(&bind_file); err != nil {
    error_message := fmt.Sprintf("Failed to process request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Error": error_message})
		return
	}

  file := bind_file.File
  destination := filepath.Base(file.Filename)

  if err := c.SaveUploadedFile(file, destination); err != nil {
    error_message := fmt.Sprintf("Failed to upload file: %s", err.Error())
    c.JSON(http.StatusBadRequest, gin.H{"Error": error_message})
		return
  }

	desired_type := c.Params.ByName("desired_type")
	available_types := []string{"pdf", "docx", "txt"}

	if !slices.Contains(available_types, desired_type) {
		c.JSON(http.StatusConflict, gin.H{"Error": "Invalid file type"})
		return
	} else {
		c.JSON(http.StatusOK, bind_file)
	}

}