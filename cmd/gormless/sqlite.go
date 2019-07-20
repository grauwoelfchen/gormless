// +build !mssql,!mysql,!postgres,sqlite

package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func connect(databaseURL string) (*gorm.DB, error) {
	return gorm.Open("sqlite3", databaseURL)
}
