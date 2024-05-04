package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"type:varchar(100);unique_index"`
	Sessions []Session
}

type Session struct {
	gorm.Model
	UserID uint
	Token  string `gorm:"type:varchar(100);unique_index"`
}

var databaseConnection *gorm.DB
var connectionError error

func initializeDatabase() {
	connectionError = godotenv.Load()
	if connectionError != nil {
		log.Fatalf("Error loading .env file")
	}

	databaseDriver := os.Getenv("DB_DRIVER")
	databaseUsername := os.Getenv("DB_USER")
	databasePassword := os.Getenv("DB_PASS")
	databaseName := os.Getenv("DB_NAME")
	databaseHost := os.Getenv("DB_HOST")
	databasePort := os.Getenv("DB_PORT")

	databaseSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", databaseUsername, databasePassword, databaseHost, databasePort, databaseName)
	if databaseDriver == "postgres" {
		databaseSourceName = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", databaseHost, databasePort, databaseUsername, databaseName, databasePassword)
	}

	databaseConnection, connectionError = gorm.Open(databaseDriver, databaseSourceName)
	if connectionError != nil {
		log.Fatalf("Could not connect to database: %v", connectionError)
	}

	databaseConnection.AutoMigrate(&User{}, &Session{})
}

func createUserInDatabase(user *User) {
	databaseConnection.Create(user)
}

func fetchUserByID(userID uint) User {
	var user User
	databaseConnection.First(&user, userID)
	return user
}

func updateUserByID(userID uint, newData map[string]interface{}) {
	databaseConnection.Model(&User{}).Where("id = ?", userID).Updates(newData)
}

func deleteUserByID(userID uint) {
	var user User
	databaseConnection.First(&user, userID)
	databaseConnection.Delete(&user)
}

func createSessionInDatabase(session *Session) {
	databaseConnection.Create(session)
}

func fetchSessionByID(sessionID uint) Session {
	var session Session
	databaseConnection.First(&session, sessionID)
	return session
}

func updateSessionByID(sessionID uint, newData map[string]interface{}) {
	databaseConnection.Model(&Session{}).Where("id = ?", sessionID).Updates(newData)
}

func deleteSessionByID(sessionID uint) {
	var session Session
	databaseConnection.First(&session, sessionID)
	databaseConnection.Delete(&session)
}

func main() {
	initializeDatabase()
}