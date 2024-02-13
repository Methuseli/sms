package routes

import (
	"github.com/Methuseli/sms/users/controller"
	"github.com/Methuseli/sms/users/middleware"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, adapter *gormadapter.Adapter) {
	users := router.Group("/users")
	users.POST("/signup", middleware.Authorize(adapter), controller.Register)
	users.POST("/login", controller.Login)
	users.GET("", middleware.Authorize(adapter), controller.GetAllUsers)
	users.PATCH("/:id", middleware.Authorize(adapter), controller.EditUser)
	users.GET("/:id", middleware.Authorize(adapter), controller.GetUser)
	users.POST("/forgot-password/:id", controller.ForgotPassword)
}
