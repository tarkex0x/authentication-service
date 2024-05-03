package main

import (
    "errors"
    "fmt"
    "log"
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
        // Log the error for debugging purposes
        log.Printf("Error generating bcrypt hash for user %s: %v", username, err)
        return err
    }
    if err := db.Create(&User{Username: username, Password: string(hashedPassword)}).Error; err != nil {
        // More detailed logging on database operations
        log.Printf("Error creating user %s in the database: %v", username, err)
        return err
    }
    return nil
}

func AuthenticateUser(db *gorm.DB, username, password string) (bool, error) {
    var user User
    if result := db.First(&user, "username = ?", username); errors.Is(result.Error, gorm.ErrRecordNotFound) {
        log.Printf("User %s not found", username)
        return false, errors.New("user not found")
    } else if result.Error != nil {
        log.Printf("Error retrieving user %s: %v", username, result.Error)
        return false, result.Error
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        log.Printf("Error comparing hash for user %s: %v", username, err)
        // Consider not revealing too specific error details for security reasons
        return false, errors.New("authentication failed")
    }
    return true, nil
}

func InitializeDB() (*gorm.DB, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Println("DATABASE_URL is not set")
        return nil, errors.New("database URL is not set")
    }

    db, err := gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
    if err != nil {
        log.Printf("Failed to connect to database: %v", err)
        return nil, err
    }
    if err = db.AutoMigrate(&User{}); err != nil {
        log.Printf("Failed to migrate database: %v", err)
        return nil, err
    }
    return db, nil
}

func main() {
    db, err := InitializeDB()
    if err != nil {
        fmt.Println("Failed to initialize database:", err)
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