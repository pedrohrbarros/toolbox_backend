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

	if err := router.Run(os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}
