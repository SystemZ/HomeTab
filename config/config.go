package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	DB_HOST                   string
	DB_PORT                   string
	DB_NAME                   string
	DB_USERNAME               string
	DB_PASSWORD               string
	REDIS_HOST                string
	REDIS_PASSWORD            string
	SESSION_VALID_S           int
	REGISTER_ON               bool
	REGISTER_WHITELIST        bool
	REGISTER_TOKEN            string
	MQTT_VHOST                string
	MQTT_EXTERNAL_SERVER_HOST string
	MQTT_EXTERNAL_SERVER_PORT int
)

func init() {
	if os.Getenv("ENVIRONMENT") == "DEV" {
		log.Println("Dev env detected")
	} else {
		log.Println("Production env detected")
	}
	viper.AutomaticEnv()
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
	// TaskTab specific
	viper.SetDefault("SESSION_VALID_S", 2592000)
	SESSION_VALID_S = viper.GetInt("SESSION_VALID_S")
	viper.SetDefault("REGISTER_WHITELIST", true)
	REGISTER_WHITELIST = viper.GetBool("REGISTER_WHITELIST")
	viper.SetDefault("REGISTER_ON", false)
	REGISTER_ON = viper.GetBool("REGISTER_ON")
	viper.SetDefault("REGISTER_TOKEN", "unknown")
	REGISTER_TOKEN = viper.GetString("REGISTER_TOKEN")
	// MQTT
	viper.SetDefault("MQTT_VHOST", "tasktab")
	MQTT_VHOST = viper.GetString("MQTT_VHOST")
	viper.SetDefault("MQTT_EXTERNAL_SERVER_HOST", "127.0.0.1")
	MQTT_EXTERNAL_SERVER_HOST = viper.GetString("MQTT_EXTERNAL_SERVER_HOST")
	viper.SetDefault("MQTT_EXTERNAL_SERVER_PORT", 1883)
	MQTT_EXTERNAL_SERVER_PORT = viper.GetInt("MQTT_EXTERNAL_SERVER_PORT")
}
