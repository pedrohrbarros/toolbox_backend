package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pedrohrbarros/toolbox_backend/src/controller/url"
)

func InitRoutes(r *gin.RouterGroup) {
	r.POST("/shortener-url", url.ShortUrl)
	r.POST("/convert-doc/:desired_type")
	r.POST("/password-generator")
}