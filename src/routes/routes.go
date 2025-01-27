package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pedrohrbarros/toolbox_backend/src/controller/url"
	"github.com/pedrohrbarros/toolbox_backend/src/controller/file"
	"github.com/pedrohrbarros/toolbox_backend/src/controller/secret"
	"github.com/pedrohrbarros/toolbox_backend/src/controller/image"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(r *gin.RouterGroup) {
	r.POST("/url/shortener", url.ShortUrl)
	r.POST("/file/convert", file.ConvertFile)
	r.POST("/secret/generator", secret.GenerateSecret)
	r.POST("/image/resize", image.ResizeImage)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}