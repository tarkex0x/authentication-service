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

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_ADDR"), 
	Password: os.Getenv("REDIS_PASSWORD"), 
	DB:       0,                          
})

func generateSessionToken() (string, error) {
	b := make([]byte, 16) 
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	return token, nil
}

func CreateSession(userID string, ttl time.Duration) (string, error) {
	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	err = rdb.Set(ctx, token, userID, ttl).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func CheckSession(token string) (bool, error) {
	result, err := rdb.Get(ctx, token).Result()
	if err == redis.Nil {
		return false, nil 
	} else if err != nil {
		return false, err 
	}

	return result != "", nil 
}

func DeleteSession(token string) error {
	_, err := rdb.Del(ctx, token).Result()
	return err
}