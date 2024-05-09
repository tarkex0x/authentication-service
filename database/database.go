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

var dbInstance *gorm.DB
var dbConnectionError error

func main() {
	initDB()
}

func initDB() {
	loadEnvVars()
	connectToDB()
	migrateSchemas()
}

func loadEnvVars() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func connectToDB() {
	dbInstance, dbConnectionError = gorm.Open(getDBDriver(), getDBConnectionString())
	if dbConnectionError != nil {
		log.Fatalf("Failed to connect to database: %v", dbConnectionError)
	}
}

func getDBDriver() string {
	return os.Getenv("DB_DRIVER")
}

func getDBConnectionString() string {
	driver := getDBDriver()
	if driver == "postgres" {
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
			os.Getenv("DB_NAME"), os.Getenv("DB_PASS"))
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}

func migrateSchemas() {
	dbInstance.AutoMigrate(&User{}, &Session{})
}

func createUser(user *User) {
	dbInstance.Create(user)
}

func findUserByID(userID uint) User {
	var user User
	dbInstance.First(&user, userID)
	return user
}

func updateUserAttributes(userID uint, newAttributes map[string]interface{}) {
	dbInstance.Model(&User{}).Where("id = ?", userID).Updates(newAttributes)
}

func deleteUser(userID uint) {
	var user User
	dbInstance.First(&user, userID)
	dbInstance.Delete(&user)
}

func createSession(session *Session) {
	dbInstance.Create(session)
}

func findSessionByID(sessionID uint) Session {
	var session Session
	dbInstance.First(&session, sessionID)
	return session
}

func updateSessionAttributes(sessionID uint, newAttributes map[string]interface{}) {
	dbInstance.Model(&Session{}).Where("id = ?", sessionID).Updates(newAttributes)
}

func deleteSession(sessionID uint) {
	var session Session
	dbInstance.First(&session, sessionID)
	dbInstance.Delete(&session)
}