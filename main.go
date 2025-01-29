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

	router := gin.Default()

  router.Use(cors.Default())

	routes.InitRoutes(&router.RouterGroup)

	router.Use(gin.Recovery())

	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Description = "Toolbox API documentation for any day-to-day function"
	docs.SwaggerInfo.Title = "Toolbox Swagger UI"

	port := cmp.Or(os.Getenv("PORT"), "8000")

	log.Print("Server running on port " + port)

	router.Run((":" + port))
}
