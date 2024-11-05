package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	docs "github.com/pedrohrbarros/toolbox_backend/docs"
	"github.com/pedrohrbarros/toolbox_backend/src/routes"
)

func main() {

	router := gin.Default()

	routes.InitRoutes(&router.RouterGroup)

	docs.SwaggerInfo.BasePath = "/swagger/"

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	if err := router.Run(port); err != nil {
		log.Fatal(err)
	}
}
