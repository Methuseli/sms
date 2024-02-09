package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/Methuseli/sms/database"
	"github.com/Methuseli/sms/environment"
	"github.com/Methuseli/sms/models"
	"github.com/Methuseli/sms/utilities"
	"github.com/gin-gonic/gin"
	Log "github.com/sirupsen/logrus"
)

func GoogleOAuth(context *gin.Context) {
	code := context.Query("code")

	if code == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Authorization code not provided!"})
		return
	}

	tokenRes, err := utilities.GetGoogleOauthToken(code)

	Log.WithFields(Log.Fields{"message": "token retrieval", "error": "Error"}).Debugln(tokenRes)

	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	googleUser, err := utilities.GetGoogleUser(tokenRes.Access_token, tokenRes.Id_token)

	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	email := strings.ToLower(googleUser.Email)

	userData := models.User{
		Firstname: googleUser.Givenname,
		Lastname:  googleUser.Familyname,
		Email:     email,
		Password:  "",
		Provider:  "Google",
	}

	if database.Database.Model(&userData).Where("email = ?", email).Updates(&userData).RowsAffected == 0 {
		database.Database.Create(&userData)
	}

	var user models.User
	database.Database.First(&user, "email = ?", email)

	token, err := utilities.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	http.SetCookie(context.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24), // Adjust expiration as needed
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "Successfully logged in..."})
}

func GoogleOAuthRedirect(context *gin.Context) {

	config, _ := environment.LoadConfig(".")
	// Construct the Google authorization URL
	googleAuthURL := "https://accounts.google.com/o/oauth2/auth" +
		"?client_id=" + config.GoogleClientID +
		"&redirect_uri=" + config.GoogleOAuthRedirectUrl +
		"&response_type=code" +
		"&scope=email%20profile%20openid"

	// Redirect the user to the authorization URL
	context.Redirect(http.StatusTemporaryRedirect, googleAuthURL)
}
