package main

import (
	"cmp"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	docs "github.com/pedrohrbarros/toolbox_backend/docs"
	"github.com/pedrohrbarros/toolbox_backend/src/routes"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

  router.Use(cors.Default())

	routes.InitRoutes(&router.RouterGroup)

	router.Use(gin.Recovery())

	docs.SwaggerInfo.BasePath = "/swagger/"

	port := cmp.Or(os.Getenv("PORT"), "8000")

	log.Print("Server running on port " + port)

	router.Run((":" + port))
}
