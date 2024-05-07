package utils

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
)

func InitializeEnvironmentVariables() {
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: Error loading .env file: %v", err)
    }
}

func EnvironmentVariableOrDefault(key, fallbackValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallbackValue
}

func RespondWithError(writer http.ResponseWriter, statusCode int, errorMessage string) {
    if !RespondWithJSON(writer, statusCode, map[string]string{"error": errorMessage}) {
        log.Println("Failed to send error response")
    }
}

func RespondWithJSON(writer http.ResponseWriter, statusCode int, payload interface{}) bool {
    jsonResponse, err := json.Marshal(payload)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        if _, writeErr := writer.Write([]byte("HTTP 500: Internal Server Error")); writeErr != nil {
            log.Printf("Error writing failure response: %v", writeErr)
            return false
        }
        return true
    }
    writer.Header().Set("Content-Type", "application/json")
    writer.WriteHeader(statusCode)
    if _, writeErr := writer.Write(jsonResponse); writeErr != nil {
        log.Printf("Error writing JSON response: %v", writeErr)
        return false
    }
    return true
}

func ValidateRequestPayload(payload interface{}) (isValid bool, errorMessage string) {
    return true, ""
}