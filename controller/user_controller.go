package controller

import (
	"net/http"
	"time"

	"github.com/Methuseli/sms/models"
	"github.com/Methuseli/sms/utilities"
	"github.com/gin-gonic/gin"
)

func Register(context *gin.Context) {
	var input models.User

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username:    input.Username,
		Password:    input.Password,
		Firstname:   input.Firstname,
		Middlename:  input.Middlename,
		Lastname:    input.Lastname,
		Email:       input.Email,
		Phonenumber: input.Phonenumber,
		IsStudent:   input.IsStudent,
	}

	savedUser, err := user.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"user": savedUser})
}

func Login(context *gin.Context) {
	var input models.AuthenticationInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.FindUserByUsername(input.Username)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.ValidatePassword(input.Password)

	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	jwt, err := utilities.GenerateJWT(user)
    http.SetCookie(context.Writer, &http.Cookie{
        Name:     "token",
        Value:    jwt,
        Expires:  time.Now().Add(time.Hour * 24), // Adjust expiration as needed
        Path:     "/",
        Secure:   true,
        HttpOnly: true,
    })
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Successfully logged in..."})
}
