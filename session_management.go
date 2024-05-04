package sessionmanager

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "os"
    "sync"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/joho/godotenv"
)

func init() {
    if err := godotenv.Load(); err != nil {
        panic("No .env file found")
    }
}

var (
    appContext = context.Background()

    redisClient = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),    // Redis server address
        Password: os.Getenv("REDIS_PASSWORD"), // Redis server password
        DB:       0,                           // Default Redis database
    })

    // Example of a cache map and its mutex.
    cache        = make(map[string]string)
    cacheMutex   = &sync.RWMutex{}
)

// generateUniqueSessionToken internally generates a new, random token for a session.
func generateUniqueSessionToken() (string, error) {
    randomBytes := make([]byte, 16) // using 128 bits
    if _, err := rand.Read(randomBytes); err != nil {
        return "", err
    }
    sessionToken := hex.EncodeToString(randomBytes)
    return sessionToken, nil
}

// CreateNewSession generates a unique session token for a user, stores it in Redis, and returns the token.
func CreateNewSession(userID string, sessionTTL time.Duration) (string, error) {
    sessionToken, err := generateUniqueSessionToken()
    if err != nil {
        return "", err
    }

    // Store the session token with the user id as value in Redis
    if err = redisClient.Set(appContext, sessionToken, userID, sessionTTL).Err(); err != nil {
        return "", err
    }

    return sessionToken, nil
}

// ValidateSessionToken checks if a given session token is valid and returns a boolean value accordingly.
func ValidateSessionToken(receivedToken string) (bool, error) {
    // Utilizing caching mechanism
    cacheMutex.RLock()
    if userID, found := cache[receivedToken]; found {
        cacheMutex.RUnlock()
        return userID != "", nil
    }
    cacheMutex.RUnlock()

    userID, err := redisClient.Get(appContext, receivedToken).Result()
    if err == redis.Nil {
        return false, nil // Token is not present in Redis
    } else if err != nil {
        return false, err // An error occurred while fetching token details
    }

    // Caching the token result after fetching from Redis
    cacheMutex.Lock()
    cache[receivedToken] = userID
    cacheMutex.Unlock()

    // Token is considered valid if userID is not empty
    return userID != "", nil
}

// InvalidateSessionToken removes a session token from Redis, effectively ending the session.
func InvalidateSessionToken(sessionToken string) error {
    _, err := redisClient.Del(appContext, sessionToken).Result()
    // Assuming invalidation should also clear from cache
    cacheMutex.Lock()
    delete(cache, sessionToken)
    cacheMutex.Unlock()
    return err
}