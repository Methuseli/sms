package database

import (
	"fmt"
	"log"

	"github.com/Methuseli/sms/users/environment"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDB(config *environment.Config) {
	var err error
	// fmt.Printf("Password: %s \n", config.DBUserPassword)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Johannesburg", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// fmt.Println(err)
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("🚀 Connected Successfully to the Database")
}