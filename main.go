package main

import (
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "log"
    "net/http"
    "os"
)

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Set up Gin router
    r := gin.Default()
    r.Use(gin.Logger())
    r.Use(errorHandlingMiddleware())
    setUpRoutes(r)

    // Get port from .env, default to a predetermined port if not found
    port := os.Getenv("PORT")
    if port == "" {
        log.Fatalf("Port not defined in .env file")
    }

    // Start the HTTP server
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func setUpRoutes(router *gin.Engine) {
    // Define routes and associate handlers
    router.POST("/register", registerUser)
    router.POST("/login", loginUser)
    router.GET("/validateSession", validateSession)
}

func registerUser(c *gin.Context) {
    // Placeholder for user registration logic
    // TODO: Add database interaction and proper error handling
    if success := mockDatabaseInsert(); !success {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "User registration"})
}

func loginUser(c *gin.Context) {
    // Placeholder for user login logic
    // TODO: Add authentication logic and proper error handling
    if authenticated := mockUserAuthentication(); !authenticated {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "User login"})
}

func validateSession(c *gin.Context) {
    // Placeholder for session validation logic
    // TODO: Implement session validation and proper error handling
    if valid := mockSessionValidation(); !valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Session validation"})
}

func errorHandlingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Process request
        c.Next()

        // Check if there are any errors
        if len(c.Errors) > 0 {
            // Log the error
            log.Printf("Error encountered: %v", c.Errors[0])
            // Respond with the first error encountered
            c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors[0].Error()})
        }
    }
}

// Mock functions to simulate database operations and logic, for demonstration purposes
func mockDatabaseInsert() bool {
    // Simulate database insert operation
    return true // Change to false to simulate an insert error
}

func mockUserAuthentication() bool {
    // Simulate user authentication
    return true // Change to false to simulate authentication failure
}

func mockSessionValidation() bool {
    // Simulate session validation
    return true // Change to false to simulate session validation failure
}