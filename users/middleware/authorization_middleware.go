package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Methuseli/sms/users/database"
	"github.com/Methuseli/sms/users/models"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	Log "github.com/sirupsen/logrus"
)

func Authorize(adapter *gormadapter.Adapter) gin.HandlerFunc {
	return func(context *gin.Context) {
		// Get current user/subject

		// Check if accessing a protected endpoint

		token, err := context.Cookie("token")
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "error validating authentication token"})
		}
		authenticationToken, err := jwt.Parse(token, func(accessToken *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_PRIVATE_KEY")), nil
		})

		if err != nil || !authenticationToken.Valid {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "user not authenticated"})
			return
		}

		claims, ok := authenticationToken.Claims.(jwt.MapClaims)
		Log.Printf("claims: %#v", claims)
		Log.Debug(ok)

		if !ok {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Invalid token"})
			return
		}

		userID := claims["id"].(string)
		context.Set("userId", userID)
		Log.Debug("User ID: ", userID)

		var user models.User
		if err := database.Database.Where("id = ?", userID).First(&user).Error; err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to fetch user"})
			return
		}

		role := user.Role
		resource := context.Request.URL.Path
		method := context.Request.Method
		context.Set("userRole", role)

		configFile := "environment/rbac_model.conf"

		// Casbin enforces policy
		ok, err = enforce(role, resource, method, adapter, configFile)
		if err != nil {
			Log.Println(err)
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "error occurred when authorizing user"})
			return
		}

		if ok {
			context.Next()
		}

		configFile = "environment/rbac_model_m2.conf"
		ok, err = enforce(role, resource, method, adapter, configFile)

		if !ok {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "forbidden"})
			return
		}

		context.Next()
	}
}

func enforce(sub string, obj string, act string, adapter *gormadapter.Adapter, configFile string) (bool, error) {
	// Load model configuration file and policy store adapter
	Log.WithFields(Log.Fields{"subject": sub, "resource": obj, "method": act}).Debug("Load model configuration file and policy store")
	enforcer, err := casbin.NewEnforcer(configFile, adapter)
	if err != nil {
		return false, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}
	// Load policies from DB dynamically
	err = enforcer.LoadPolicy()


	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}


	// Verify
	ok, err := enforcer.Enforce(sub, obj, act)
	return ok, err
}
