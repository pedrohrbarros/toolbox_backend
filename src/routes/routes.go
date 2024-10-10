package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pedrohrbarros/toolbox_backend/src/controller/url"
	"github.com/pedrohrbarros/toolbox_backend/src/controller/converter"
)

func InitRoutes(r *gin.RouterGroup) {
	r.POST("/shortener-url", url.ShortUrl)
	r.POST("/convert/", converter.ConvertFile)
	r.POST("/secret-generator")
}