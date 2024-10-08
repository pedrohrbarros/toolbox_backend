package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pedrohrbarros/toolbox_backend/src/controller/url"
	"github.com/pedrohrbarros/toolbox_backend/src/controller/document"
)

func InitRoutes(r *gin.RouterGroup) {
	r.POST("/shortener-url", url.ShortUrl)
	r.POST("/convert-doc/:desired_type", document.convertFile)
	r.POST("/secret-generator")
}