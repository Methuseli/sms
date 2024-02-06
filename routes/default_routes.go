package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DefaultRoutes(router *gin.RouterGroup) {
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to SMS"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
