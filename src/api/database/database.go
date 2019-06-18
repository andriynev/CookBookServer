package database

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var dbConn *gorm.DB

func GetDB() (db *gorm.DB, err error) {
	if dbConn == nil {
		return nil, errors.New("database is not initialized")
	}
	return dbConn, nil
}

func InitDB(dsn string) (db *gorm.DB, err error) {
	if dbConn != nil {
		return dbConn, nil
	}
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}

	dbConn = db.LogMode(true)
	return
}
