package routes

import (
	"github.com/Methuseli/sms/controller"
	"github.com/gin-gonic/gin"
)

func OAuthRoutes(router *gin.RouterGroup) {
	oauth := router.Group("oauth2/")
	oauth.GET("login/google", controller.GoogleOAuthRedirect)
	oauth.GET("callback/google", controller.GoogleOAuth)
}