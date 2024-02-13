package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func SessionMiddleware(context *gin.Context) {
    // Check for session cookie
    sessionCookie, err := context.Cookie("session")
    if err != nil {
        // If cookie doesn't exist, create one
        sessionCookie = CreateSessionCookie(context)
        context.SetCookie("session", sessionCookie, 3600*24, "/", "", false, true) // Set appropriate options
    }

    jwtToken, err := jwt.ParseWithClaims(sessionCookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Use your secret key for verification
        return []byte(os.Getenv("SECRET_KEY")), nil
    })

    if err != nil || !jwtToken.Valid {
        context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "invalid session"})
        return
    }

    context.Next()
}

func CreateSessionCookie(context *gin.Context) string {
    // Generate a unique session ID or use a library
    sessionID, err := GenerateRandomSessionID()
    if err != nil {
        context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "error authentication user"})
        return ""
    }

    // Create a JWT token with session ID as claim
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sessionID": sessionID})
    tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

    return tokenString
}

// func IsPrivateEndpoint(endpoint string) bool {
//     // Define your logic for determining private endpoints here
//     // For example, you can check if it starts with "/api/private"
//     return strings.HasPrefix(endpoint, "/api/v1/private")
// }

func GenerateRandomSessionID() (string, error) {
    // Generate random bytes
    randomBytes := make([]byte, 16)
    _, err := rand.Read(randomBytes)
    if err != nil {
        return "", err
    }

    // Convert bytes to hexadecimal string
    return hex.EncodeToString(randomBytes), nil
}
