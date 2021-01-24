package config

import (
	"github.com/spf13/viper"
)

var (
	DB_HOST     string
	DB_PORT     string
	DB_USERNAME string
	DB_PASSWORD string
	DB_NAME     string

	REDIS_HOST     string
	REDIS_PASSWORD string

	GIT_COMMIT string

	DEV_MODE        bool
	CACHE_DIR       string
	LIVE_VID_THUMB  bool
	SESSION_VALID_S int
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("DB_HOST", "localhost")
	DB_HOST = viper.GetString("DB_HOST")

	viper.SetDefault("DB_PORT", "3306")
	DB_PORT = viper.GetString("DB_PORT")

	viper.SetDefault("DB_USERNAME", "dev")
	DB_USERNAME = viper.GetString("DB_USERNAME")

	viper.SetDefault("DB_PASSWORD", "dev")
	DB_PASSWORD = viper.GetString("DB_PASSWORD")

	viper.SetDefault("DB_NAME", "dev")
	DB_NAME = viper.GetString("DB_NAME")

	viper.SetDefault("REDIS_HOST", "localhost")
	REDIS_HOST = viper.GetString("REDIS_HOST")

	viper.SetDefault("REDIS_PASSWORD", "")
	REDIS_PASSWORD = viper.GetString("REDIS_PASSWORD")

	viper.SetDefault("DEV_MODE", false)
	DEV_MODE = viper.GetBool("DEV_MODE")

	viper.SetDefault("CACHE_DIR", "")
	CACHE_DIR = viper.GetString("CACHE_DIR")

	viper.SetDefault("LIVE_VID_THUMB", true)
	LIVE_VID_THUMB = viper.GetBool("LIVE_VID_THUMB")

	viper.SetDefault("SESSION_VALID_S", 2592000)
	SESSION_VALID_S = viper.GetInt("SESSION_VALID_S")
}