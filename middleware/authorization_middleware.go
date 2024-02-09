package middleware

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	Log "github.com/sirupsen/logrus"
)

func Authorize(obj string, act string, adapter *gormadapter.Adapter) gin.HandlerFunc {
    return func(context *gin.Context) {
                // Get current user/subject
        val, existed := context.Get("current_subject")
        if !existed {
            context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "user hasn't logged in yet"})
            return
        }
        // Casbin enforces policy
        ok, err := enforce(val.(string), obj, act, adapter)
        if err != nil {
            Log.Println(err)
            context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "error occurred when authorizing user"})
            return
        }
        if !ok {
            context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "forbidden"})
            return
        }
        context.Next()
    }
}

func enforce(sub string, obj string, act string, adapter *gormadapter.Adapter) (bool, error) {
        // Load model configuration file and policy store adapter
    enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
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