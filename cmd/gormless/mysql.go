// +build !mssql,mysql,!postgres,!sqlite

package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func connect(databaseURL string) (*gorm.DB, error) {
	return gorm.Open("mysql", databaseURL)
}
