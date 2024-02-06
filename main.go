package main

import (
	"log"

	"github.com/Methuseli/sms/database"
	"github.com/Methuseli/sms/environment"
	"github.com/Methuseli/sms/middleware"
	"github.com/Methuseli/sms/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	routes "github.com/Methuseli/sms/routes"
)

var (
	server *gin.Engine
)

func init() {
	config, err := environment.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	database.ConnectDB(&config)
	database.Database.AutoMigrate(&models.User{})

	server = gin.Default()
}

func main() {
	config, err := environment.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))
	server.Use(middleware.SessionMiddleware)

	router := server.Group("/api/v1")

	routes.DefaultRoutes(router)
	routes.UserRoutes(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
