package main

import (
	"cmp"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	docs "github.com/pedrohrbarros/toolbox_backend/docs"
	"github.com/pedrohrbarros/toolbox_backend/src/routes"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	routes.InitRoutes(&router.RouterGroup)

	router.Use(gin.Recovery())

	docs.SwaggerInfo.BasePath = "/swagger/"

	port := cmp.Or(os.Getenv("PORT"), "8080")

	log.Print("Server running on port " + port)

	router.Run((":" + port))
}
