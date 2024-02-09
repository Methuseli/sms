package routes

import (
	"github.com/Methuseli/sms/controller"
	"github.com/Methuseli/sms/middleware"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, adapter *gormadapter.Adapter) {
	users := router.Group("/users")
	users.POST("/signup", middleware.Authorize("/api/v1/users/signup", "POST", adapter), controller.Register)
	users.POST("/login", controller.Login)
	users.GET("", middleware.Authorize("/api/v1/users", "GET", adapter), controller.GetAllUsers)
}
