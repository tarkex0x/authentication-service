package utils

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
)

func InitializeEnvironmentVariables() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
}

func EnvironmentVariableOrDefault(key, fallbackValue string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        return fallbackValue
    }
    return value
}

func RespondWithError(writer http.ResponseWriter, statusCode int, errorMessage string) {
    RespondWithJSON(writer, statusCode, map[string]string{"error": errorMessage})
}

func RespondWithJSON(writer http.ResponseWriter, statusCode int, payload interface{}) {
    jsonResponse, err := json.Marshal(payload)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte("HTTP 500: Internal Server Error"))
        return
    }
    writer.Header().Set("Content-Type", "application/json")
    writer.WriteHeader(statusCode)
    writer.Write(jsonResponse)
}

func ValidateRequestPayload(payload interface{}) (isValid bool, errorMessage string) {
    return true, ""
}