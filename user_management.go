package main

import (
    "errors"
    "fmt"
    "os"

    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string
    Password string
}

func CreateUser(db *gorm.DB, username, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    return db.Create(&User{Username: username, Password: string(hashedPassword)}).Error
}

func AuthenticateUser(db *gorm.DB, username, password string) (bool, error) {
    var user User
    if result := db.First(&user, "username = ?", username); errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return false, errors.New("user not found")
    } else if result.Error != nil {
        return false, result.Error
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return false, err
    }
    return true, nil
}

func InitializeDB() (*gorm.DB, error) {
    dbURL := os.Getenv("DATABASE_URL")
    db, err := gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    if err = db.AutoMigrate(&User{}); err != nil {
        return nil, err
    }
    return db, nil
}

func main() {
    db, err := InitializeDB()
    if err != nil {
        fmt.Println("Failed to connect to database:", err)
        return
    }

    username, password := "testUser", "testPassword"

    if err := CreateUser(db, username, password); err != nil {
        fmt.Println("Error creating user:", err)
        return
    }
    fmt.Println("User created successfully")

    authenticated, err := AuthenticateUser(db, username, password)
    if err != nil {
        fmt.Println("Authentication failed:", err)
        return
    }
    if authenticated {
        fmt.Println("User authenticated successfully")
    } else {
        fmt.Println("Invalid credentials")
    }
}