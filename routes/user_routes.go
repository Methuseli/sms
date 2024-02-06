package routes

import (
	"github.com/Methuseli/sms/controller"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	users.POST("/signup", controller.Register)
	users.POST("/login", controller.Login)
}
