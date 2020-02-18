package model

import (
	"fmt"
	"github.com/DavidHuie/gomigrate"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/gotag/config"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func InitMysql() *gorm.DB {
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

	//db.AutoMigrate(&File{})
	//db.AutoMigrate(&FileTag{})
	//db.AutoMigrate(&Tag{})

	logrus.Info("Connection to database seems OK!")

	logrus.Info("Starting DB migrations")
	//change to https://github.com/mattes/migrate
	migrator, _ := gomigrate.NewMigrator(db.DB(), gomigrate.Mysql{}, "./migrations")
	err = migrator.Migrate()

	DB = db

	return db
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
}
