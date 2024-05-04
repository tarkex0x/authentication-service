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

// CreateUser processes the user creation logic.
func (u *User) CreateUser(db *gorm.DB) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Error generating bcrypt hash for user %s: %v", u.Username, err)
        return errors.New("failed to generate encrypted password")
    }
    u.Password = string(hashedPassword)

    if err := db.Create(u).Error; err != nil {
        log.Printf("Error creating user %s in the database: %v", u.Username, err)
        return err
    }
    return nil
}

// Authenticate checks if the user credentials are valid.
func (u *User) Authenticate(db *gorm.DB, password string) (bool, error) {
    if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
        log.Printf("Error comparing hash for user %s: %v", u.Username, err)
        return false, errors.New("authentication failed")
    }
    return true, nil
}

// InitializeDB establishes a connection to the database.
func InitializeDB() (*gorm.DB, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "test.db" // A sensible default for SQLite
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

// FindUserByUsername queries the database for a user by username.
func FindUserByUsername(db *gorm.DB, username string) (*User, error) {
    var user User
    result := db.First(&user, "username = ?", username)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            log.Printf("User %s not found", username)
            return nil, errors.New("user not found")
        }
        log.Printf("Error retrieving user %s: %v", username, result.Error)
        return nil, result.Error
    }
    return &user, nil
}

func main() {
    db, err := InitializeDB()
    if err != nil {
        fmt.Println("Failed to initialize database:", err)
        return
    }

    username, password := "testUser", "testPassword"
    user := &User{Username: username, Password: password}

    if err := user.CreateUser(db); err != nil {
        fmt.Println("Error creating user:", err)
        return
    }
    fmt.Println("User created successfully")

    user, err = FindUserByUsername(db, username)
    if err != nil {
        fmt.Println(err)
        return
    }

    authenticated, err := user.Authenticate(db, password)
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