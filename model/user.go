package model

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

// Host represents the host for this application
// swagger:model user
type User struct {
	// ID
	//
	// required: true
	Id uint `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`

	// Username
	//
	// required: true
	Username  string     `gorm:"column:username" json:"username"`
	Email     string     `gorm:"column:email" json:"email"`
	Hash      string     `gorm:"column:hash" json:"hash"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(username string, email string, password string) {
	var user User

	hash, err := HashPassword(password)
	if err != nil {
		log.Print(err.Error())
	}
	user.Username = strings.ToLower(username)
	user.Email = strings.ToLower(email)
	user.Hash = hash
	//FIXME timezones
	now := time.Now()
	user.CreatedAt = &now

	err = DB.Save(&user).Error
	if err != nil {
		log.Printf("%v", err)
	}
}

func IsPasswordOk(username string, password string) (ok bool, usr User) {
	DB.Where(&User{Username: username}).First(&usr)
	if CheckPasswordHash(password, usr.Hash) {
		return true, usr
	}
	return false, usr
}

func GetUserById(id uint) (ok bool, usr User) {
	res := DB.Where(&User{Id: id}).First(&usr)
	if res.RowsAffected > 0 {
		return true, usr
	}
	return false, usr
}
