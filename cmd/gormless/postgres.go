// +build !mssql,!mysql,postgres,!sqlite

package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func connect(databaseURL string) (*gorm.DB, error) {
	return gorm.Open("postgres", databaseURL)
}
