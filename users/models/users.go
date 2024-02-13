package models

import (
	"github.com/Methuseli/sms/users/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	// Log "github.com/sirupsen/logrus"
)

type User struct {
	gorm.Model
	ID          string  `gorm:"primary_key" json:"id"`
	Username    string  `gorm:"size:255;not null;unique" json:"username"`
	Password    string  `gorm:"size:255;not null" json:"password"`
	Firstname   string  `gorm:"size:255;not null" json:"firstname"`
	Middlename  *string `gorm:"size:255" json:"middlename"`
	Lastname    string  `gorm:"size:255;not null" json:"lastname"`
	Email       string  `gorm:"size:255;not null;unique" json:"email"`
	Phonenumber string  `gorm:"size:255;not null;unique" json:"phonenumber"`
	IsStudent   bool    `gorm:"default:false" json:"is_student"`
	Role        string  `gorm:"default:'user'" json:"role"`
	Provider    string  `gorm:"default: 'local'" json:"provider"`
}

func (user *User) Save() (*User, error) {
	err := database.Database.Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) BeforeCreate(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.ID = uuid.New().String()
	return nil
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserByUsername(username string) (User, error) {
	var user User
	err := database.Database.Where("username=?", username).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func FindUserByEmail(email string) (User, error) {
	var user User
	err := database.Database.Where("email=?", email).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}
