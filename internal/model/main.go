package model

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/go-redis/redis_rate/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"github.com/systemz/hometab/internal/config"
	"time"
)

var (
	DB          *gorm.DB
	Redis       *redis.Client
	AuthLimiter *redis_rate.Limiter
)

func InitMysql() {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", config.DB_USERNAME, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME))
	if err != nil {
		logrus.Error(err.Error())
		panic("Failed to connect to database")
	}

	err = db.DB().Ping()
	if err != nil {
		logrus.Panic("Ping to db failed")
	}

	//https://github.com/go-sql-driver/mysql/issues/257
	db.DB().SetMaxIdleConns(0)
	db.LogMode(config.DEV_MODE)
	logrus.Info("Connection to database seems OK!")

	DB = db

	logrus.Info("Starting DB migrations")
	m, err := migrate.New(
		"file:///migrations",
		fmt.Sprintf(
			"mysql://%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
			config.DB_USERNAME,
			config.DB_PASSWORD,
			config.DB_HOST,
			config.DB_PORT,
			config.DB_NAME,
		))
	m.Up()
}

func InitRedis() {
	if len(config.REDIS_PASSWORD) > 1 {
		Redis = redis.NewClient(&redis.Options{
			Addr:     config.REDIS_HOST + ":6379",
			Password: config.REDIS_PASSWORD,
		})
	} else {
		Redis = redis.NewClient(&redis.Options{
			Addr: config.REDIS_HOST + ":6379",
		})
	}

	_, err := Redis.Ping().Result()
	if err != nil {
		logrus.Error(err.Error())
		logrus.Panic("Ping to Redis failed")
	}

	logrus.Info("Connection to Redis seems OK!")
	// configure rate limiting
	AuthLimiter = redis_rate.NewLimiter(Redis, &redis_rate.Limit{
		Burst:  5,
		Rate:   5,
		Period: time.Second * 10,
	})
}
