package main

import (
	"github.com/jinzhu/gorm"
)

// Up ...
func Up(tx *gorm.DB) error {
	type User struct {
		Name string
	}
	return tx.AutoMigrate(&User{}).Error
}
