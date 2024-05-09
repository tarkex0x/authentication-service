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

func (u *User) CreateUser(db *gorm.DB) error {
    var err error
    if u.Password, err = hashPassword(u.Password); err != nil {
        log.Printf("Error hashing password for user %s: %v", u.Username, err)
        return errors.New("failed to hash password")
    }

    if err := db.Create(u).Error; err != nil {
        log.Printf("Error creating user %s: %v", u.Username, err)
        return err
    }
    return nil
}

func hashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func (u *User) Authenticate(password string) (bool, error) {
    if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
        log.Printf("Authentication error for user %s: %v", u.Username, err)
        return false, errors.New("invalid credentials")
    }
    return true, nil
}

func InitializeDB() (*gorm.DB, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "test.db"
    }

    db, err := gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
    if err != nil {
        log.Printf("Database connection failed: %v", err)
        return nil, err
    }
    if err = db.AutoMigrate(&User{}); err != nil {
        log.Printf("Database migration failed: %v", err)
        return nil, err
    }
    return db, nil
}

func FindUserByUsername(db *gorm.DB, username string) (*User, error) {
    var user User
    result := db.First(&user, "username = ?", username)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            log.Printf("User %s not found", username)
            return nil, errors.New("user not found")
        }
        log.Printf("Error finding user %s: %v", username, result.Error)
        return nil, result.Error
    }
    return &user, nil
}

func main() {
    db, err := InitializeDB()
    if err != nil {
        fmt.Println("Error initializing database:", err)
        return
    }

    username, password := "testUser", "testPassword"
    user := User{Username: username, Password: password}

    if err := (&user).CreateUser(db); err != nil {
        fmt.Println("Error creating user:", err)
        return
    }
    fmt.Println("User successfully created")

    fetchedUser, err := FindUserByUsername(db, username)
    if err != nil {
        fmt.Println("Error finding user:", err)
        return
    }

    authenticated, err := fetchedUser.Authenticate(password)
    if err != nil {
        fmt.Println("Authentication error:", err)
        return
    }

    if authenticated {
        fmt.Println("User authenticated successfully")
    } else {
        fmt.Println("Invalid credentials")
    }
}