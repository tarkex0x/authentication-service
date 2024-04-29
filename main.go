package main

import (
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "log"
    "net/http"
    "os"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    r := gin.Default()
    r.Use(gin.Logger())
    r.Use(errorHandlingMiddleware())
    setUpRoutes(r)
    port := os.Getenv("PORT")
    if port == "" {
        log.Fatal("Port not defined in .env file")
    }
    r.Run(":" + port)
}

func setUpRoutes(router *gin.Engine) {
    router.POST("/register", registerUser)
    router.POST("/login", loginUser)
    router.GET("/validateSession", validateSession)
}

func registerUser(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "User registration"})
}

func loginUser(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "User login"})
}

func validateSession(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Session validation"})
}

func errorHandlingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors[0].Error()})
        }
    }
}