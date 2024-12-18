package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pedrohrbarros/toolbox_backend/src/controller/url"
	"github.com/pedrohrbarros/toolbox_backend/src/controller/converter"
	"github.com/pedrohrbarros/toolbox_backend/src/controller/secret"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(r *gin.RouterGroup) {
	r.POST("/shortener-url", url.ShortUrl)
	r.POST("/convert/", converter.ConvertFile)
	r.POST("/secret-generator", secret.GenerateSecret)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}