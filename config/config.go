package config

import (
	"os"
	"github.com/spf13/viper"
	"log"
)

var (
	DB_HOST string
	DB_PORT string
	DB_NAME string
	DB_USERNAME string
	DB_PASSWORD string
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
}
