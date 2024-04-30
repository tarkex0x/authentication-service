package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	DatabaseURL        string `mapstructure:"DATABASE_URL"`
	APIKey             string `mapstructure:"API_KEY"`
	SecretKey          string `mapstructure:"SECRET_KEY"`
	AdditionalSecrets  string `mapstructure:"OTHER_SENSITIVE_DATA"`
}

func LoadConfiguration(path string) (config AppConfig, err error) {
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
		fmt.Println(fmt.Errorf("Fatal error reading config file: %s \n", err))
		os.Exit(1)
	}

	err = viper.Unmarshal(&config)
	return
}