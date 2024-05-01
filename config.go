package config

import (
    "fmt"
    "log"
    "os"

    "github.com/spf13/viper"
)

type AppConfig struct {
    DatabaseURL       string `mapstructure:"DATABASE_URL"`
    APIKey            string `mapstructure:"API_KEY"`
    SecretKey         string `mapstructure:"SECRET_KEY"`
    AdditionalSecrets string `mapstructure:"OTHER_SENSITIVE_DATA"`
}

func SetupLogger() {
    log.SetOutput(os.Stdout) // Set the logger output to stdout
    log.SetPrefix("LOG: ")   // Every log message will be prefixed with "LOG: "
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // Log prepends the date, time, and file location
}

func LoadConfiguration(path string) (config AppConfig, err error) {
    SetupLogger()

    viper.AddConfigPath(path)
    viper.SetConfigName("app")
    viper.SetConfigType("env")

    viper.AutomaticEnv()

    viper.SetDefault("DATABASE_URL", "localhost")
    viper.SetDefault("API_KEY", "your-default-api-key")
    viper.SetDefault("SECRET_KEY", "your-default-secret-key")
    viper.SetDefault("OTHER_SENSITIVE_DATA", "default-sensitive-data")

    err = viper.ReadInConfig()
    if err != nil {
        log.Printf("Fatal error reading config file: %v\n", err)
        return config, fmt.Errorf("fatal error reading config file: %w", err) // Improved error wrapping
    }
    log.Println("Configuration loaded successfully.")

    err = viper.Unmarshal(&config)
    if err != nil {
        log.Printf("Error unmarshalling config: %v\n", err)
        return config, fmt.Errorf("error unmarshalling config: %w", err) // Improved error wrapping
    }
    log.Println("Configuration unmarshalled successfully.")
    return config, nil // Explicitly returning nil error on success
}