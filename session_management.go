package sessionmanager

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

var appContext = context.Background()

var redisClient = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_ADDR"),
	Password: os.Getenv("REDIS_PASSWORD"),
	DB:       0,
})

func generateUniqueSessionToken() (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	sessionToken := hex.EncodeToString(randomBytes)
	return sessionToken, nil
}

func CreateNewSession(userID string, sessionTTL time.Duration) (string, error) {
	sessionToken, err := generateUniqueSessionToken()
	if err != nil {
		return "", err
	}

	err = redisClient.Set(appContext, sessionToken, userID, sessionTTL).Err()
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func ValidateSessionToken(receivedToken string) (bool, error) {
	userID, err := redisClient.Get(appContext, receivedToken).Result()
	if err == redis.Nil {
		return false, nil // Session token does not exist
	} else if err != nil {
		return false, err // Error while fetching from Redis
	}

	return userID != "", nil // True if session exists
}

func InvalidateSessionToken(sessionToken string) error {
	_, err := redisClient.Del(appContext, sessionToken).Result()
	return err
}