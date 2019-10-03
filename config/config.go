package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	DB_HOST         string
	DB_PORT         string
	DB_NAME         string
	DB_USERNAME     string
	DB_PASSWORD     string
	REDIS_HOST      string
	REDIS_PASSWORD  string
	SESSION_VALID_S int
)

func init() {
	if os.Getenv("ENVIRONMENT") == "DEV" {
		log.Println("Dev env detected")
	} else {
		log.Println("Production env detected")
	}
	viper.AutomaticEnv()

	viper.SetDefault("DB_HOST", "localhost")
	DB_HOST = viper.GetString("DB_HOST")
	viper.SetDefault("DB_PORT", "3306")
	DB_PORT = viper.GetString("DB_PORT")
	viper.SetDefault("DB_NAME", "dev")
	DB_NAME = viper.GetString("DB_NAME")
	viper.SetDefault("DB_USERNAME", "dev")
	DB_USERNAME = viper.GetString("DB_USERNAME")
	viper.SetDefault("DB_PASSWORD", "dev")
	DB_PASSWORD = viper.GetString("DB_PASSWORD")
	viper.SetDefault("REDIS_HOST", "localhost")
	REDIS_HOST = viper.GetString("REDIS_HOST")
	viper.SetDefault("REDIS_PASSWORD", "")
	REDIS_PASSWORD = viper.GetString("REDIS_PASSWORD")
	viper.SetDefault("SESSION_VALID_S", 3600)
	SESSION_VALID_S = viper.GetInt("SESSION_VALID_S")
}
