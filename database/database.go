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
var err error

func initDB() {
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	if dbDriver == "postgres" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPass)
	}

	db, err = gorm.Open(dbDriver, dsn)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	db.AutoMigrate(&User{}, &Session{})
}

func createUser(user *User) {
	db.Create(user)
}

func getUser(id uint) User {
	var user User
	db.First(&user, id)
	return user
}

func updateUser(id uint, newData map[string]interface{}) {
	db.Model(&User{}).Where("id = ?", id).Updates(newData)
}

func deleteUser(id uint) {
	var user User
	db.First(&user, id)
	db.Delete(&user)
}

func createSession(session *Session) {
	db.Create(session)
}

func getSession(id uint) Session {
	var session Session
	db.First(&session, id)
	return session
}

func updateSession(id uint, newData map[string]interface{}) {
	db.Model(&Session{}).Where("id = ?", id).Updates(newData)
}

func deleteSession(id uint) {
	var session Session
	db.First(&session, id)
	db.Delete(&session)
}

func main() {
	initDB()
}