package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	DEV_MODE           bool
	HTTP_PORT          string
	TEMPLATE_PATH      string
	DB_HOST            string
	DB_PORT            string
	DB_NAME            string
	DB_USERNAME        string
	DB_PASSWORD        string
	REDIS_HOST         string
	REDIS_PASSWORD     string
	SESSION_VALID_S    int
	REGISTER_ON        bool
	REGISTER_WHITELIST bool
	REGISTER_TOKEN     string
	PUSHY_ME_SECRET    string
	// gotag part
	CACHE_DIR      string
	LIVE_VID_THUMB bool
	//
	GIT_COMMIT string
)

func init() {
	if os.Getenv("ENVIRONMENT") == "DEV" {
		log.Println("Dev env detected")
	} else {
		log.Println("Production env detected")
	}
	viper.AutomaticEnv()
	// dev
	viper.SetDefault("DEV_MODE", false)
	DEV_MODE = viper.GetBool("DEV_MODE")
	// general
	viper.SetDefault("HTTP_PORT", "3000")
	HTTP_PORT = viper.GetString("HTTP_PORT")
	viper.SetDefault("TEMPLATE_PATH", "./templates/")
	TEMPLATE_PATH = viper.GetString("TEMPLATE_PATH")

	// DB stuff
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
	// app specific
	viper.SetDefault("SESSION_VALID_S", 2592000)
	SESSION_VALID_S = viper.GetInt("SESSION_VALID_S")
	viper.SetDefault("REGISTER_WHITELIST", true)
	REGISTER_WHITELIST = viper.GetBool("REGISTER_WHITELIST")
	viper.SetDefault("REGISTER_ON", false)
	REGISTER_ON = viper.GetBool("REGISTER_ON")
	viper.SetDefault("REGISTER_TOKEN", "unknown")
	REGISTER_TOKEN = viper.GetString("REGISTER_TOKEN")
	// pushy.me
	viper.SetDefault("PUSHY_ME_SECRET", "changeme")
	PUSHY_ME_SECRET = viper.GetString("PUSHY_ME_SECRET")

	// gotag part
	viper.SetDefault("CACHE_DIR", "")
	CACHE_DIR = viper.GetString("CACHE_DIR")
	viper.SetDefault("LIVE_VID_THUMB", true)
	LIVE_VID_THUMB = viper.GetBool("LIVE_VID_THUMB")
}
