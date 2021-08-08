package util

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

//BcryptPassword hash密码
func BcryptPassword(raw string) string {
	password, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(password)
}

//BcryptCompare 比较密码
func BcryptCompare(input, hashedPassword string) error {
	if len(hashedPassword) == 0 {
		return errors.New("bcrypt-hash不能为空")
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input))
}
