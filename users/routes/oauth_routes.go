package routes

import (
	"github.com/Methuseli/sms/users/controller"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

func OAuthRoutes(router *gin.RouterGroup, adapter *gormadapter.Adapter) {
	oauth := router.Group("oauth2/")
	oauth.GET("login/google", controller.GoogleOAuthRedirect)
	oauth.GET("callback/google", controller.GoogleOAuth)
}