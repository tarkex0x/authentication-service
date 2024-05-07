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

var db *gorm.DB
var dbConnectError error

func main() {
    initializeDB()
}

func initializeDB() {
    loadEnvironmentVariables()
    establishDBConnection()
    migrateDBSchemas()
}

func loadEnvironmentVariables() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
}

func establishDBConnection() {
    db, dbConnectError = gorm.Open(determineDBDriver(), assembleDBConnectionString())
    if dbConnectError != nil {
        log.Fatalf("Could not establish connection with the database: %v", dbConnectError)
    }
}

func determineDBDriver() string {
    return os.Getenv("DB_DRIVER")
}

func assembleDBConnectionString() string {
    driver := determineDBDriver()
    if driver == "postgres" {
        return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", 
            os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), 
            os.Getenv("DB_NAME"), os.Getenv("DB_PASS"))
    }
    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", 
        os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), 
        os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}

func migrateDBSchemas() {
    db.AutoMigrate(&User{}, &Session{})
}

func addUser(user *User) {
    db.Create(user)
}

func getUserByID(userID uint) User {
    var user User
    db.First(&user, userID)
    return user
}

func updateUser(userID uint, newAttributes map[string]interface{}) {
    db.Model(&User{}).Where("id = ?", userID).Updates(newAttributes)
}

func removeUser(userID uint) {
    var user User
    db.First(&user, userID)
    db.Delete(&user)
}

func addSession(session *Session) {
    db.Create(session)
}

func getSessionByID(sessionID uint) Session {
    var session Session
    db.First(&session, sessionID)
    return session
}

func updateSession(sessionID uint, newAttributes map[string]interface{}) {
    db.Model(&Session{}).Where("id = ?", sessionID).Updates(newAttributes)
}

func removeSession(sessionID uint) {
    var session Session
    db.First(&session, sessionID)
    db.Delete(&session)
}