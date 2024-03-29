package main

import (
	"fmt"
	"os"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"

	Log "github.com/sirupsen/logrus"

	"github.com/Methuseli/sms/users/database"
	"github.com/Methuseli/sms/users/environment"
	"github.com/Methuseli/sms/users/middleware"
	"github.com/Methuseli/sms/users/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	routes "github.com/Methuseli/sms/users/routes"
)

var (
	server *gin.Engine
)

var enforcer *casbin.Enforcer

var config *environment.Config

func init() {
	// homeDir, _ := os.UserHomeDir()

    // Log.Printf(homeDir)

	Log.SetFormatter(&Log.JSONFormatter{})
	Log.SetLevel(Log.DebugLevel)
	Log.SetOutput(os.Stdout)

	config, err := environment.LoadConfig(".")
	if err != nil {
		Log.Fatal("🚀 Could not load environment variables ", err)
	}

	database.ConnectDB(&config)
	database.Database.AutoMigrate(&models.User{})

	adapter, err := gormadapter.NewAdapterByDB(database.Database)
	if err != nil {
		Log.Fatal(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	server = gin.New()
	server.Use(gin.LoggerWithWriter(Log.New().Writer()), gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))
	server.Use(middleware.SessionMiddleware)

	router := server.Group("/user-service/api/v1")
	defaultRouter := server.Group("")

	routes.DefaultRoutes(router)
	routes.UserRoutes(router, adapter)
	routes.OAuthRoutes(defaultRouter, adapter)
}

func main() {
	config, err := environment.LoadConfig(".")
	if err != nil {
		Log.Fatal("🚀 Could not load environment variables", err)
	}

	Log.Fatal(server.Run(":" + config.ServerPort))
}
