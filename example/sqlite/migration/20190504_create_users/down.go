package main

import (
	"github.com/jinzhu/gorm"
)

// Down ...
func Down(tx *gorm.DB) error {
	type User struct {
		Name string
	}
	return tx.DropTableIfExists(&User{}, "users").Error
}
