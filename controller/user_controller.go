package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/Methuseli/sms/database"
	"github.com/Methuseli/sms/models"
	"github.com/Methuseli/sms/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	Log "github.com/sirupsen/logrus"
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

	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_username_key\"") {
		context.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "user with the same username already exists"})
		return
	} else if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_email_key\"") {
		context.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "user with the same email already exists"})
		return
	} else if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_phonenumber_key\"") {
		context.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "user with the same phonenumber already exists"})
		return
	} else if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Something went wrong"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"status": "success", "user": savedUser, "message": "Successfully registered"})
}

func Login(context *gin.Context) {
	var input models.AuthenticationInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := models.FindUserByUsername(input.Username)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	err = user.ValidatePassword(input.Password)

	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
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
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "Successfully logged in..."})
}

func GetAllUsers(context *gin.Context) {
	var users []models.User

	query := database.Database.Model(&models.User{}).Find(&users)
	// Log.Println("query ", query.Error)

	if query.Error != nil {
		// Log.Println("Error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to fetch users"})
		return
	}
	context.JSON(http.StatusOK, users)
}

func checkUser(context *gin.Context, id string) bool {
	requestUserId, exists := context.Get("userId")
	if !exists {
		return false
	}
	userRole, exists := context.Get("userRole")
	if !exists {
		return false
	}

	if requestUserId == id {
		return true
	} else if userRole == "admin" {
		return true
	} else {
		return false
	}
}

func GetUser(context *gin.Context) {
	userId := context.Param("id")

	id, err := uuid.Parse(userId)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid user ID"})
		return
	}

	ok := checkUser(context, userId)
	if !ok {
		context.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Forbidden"})
		return
	}

	var user models.User
	result := database.Database.First(&user, id)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Failed to find user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "success", "user": user})
}

func EditUser(context *gin.Context) {
	userId := context.Param("id")
	id, err := uuid.Parse(userId)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid user ID"})
		return
	}

	ok := checkUser(context, userId)
	if !ok {
		context.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Forbidden"})
		return
	}

	var updateFields map[string]interface{}
	if err := context.BindJSON(&updateFields); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid JSON"})
		return
	}

	var user models.User
	result := database.Database.First(&user, id)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Failed to find user"})
		return
	}

	if err := database.Database.Model(user).Updates(updateFields).Error; err != nil {
		Log.Println(err)
		// Handle error (e.g., database error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "updated successfully"})
}
