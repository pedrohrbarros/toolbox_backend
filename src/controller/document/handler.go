package document

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func convertFile(c *gin.Context) {
  type BindFile struct {
    Name string `form: "name" binding: required`
    File *multipart.FileHeader `form: "file" binding: required`
  }

	var bind_file BindFile

	if err := c.ShouldBind(&bind_file); err != nil {
    error_message := fmt.Sprintf("Failed to process request: %s", err.Error())
		c.JSON(http.StatusBadRequest, error_message)
	}

  file := bind_file.File
  destination := filepath.Base(bind_file.Name)

  if err := c.SaveUploadedFile(file, destination); err != nil {
    error_message := fmt.Sprintf("Failed to upload file: %s", err.Error())
    c.JSON(http.StatusBadRequest, error_message)
  }

  c.JSON(http.StatusOK, bind_file)
}