package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDBInstance(dbPath string) (err error) {
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return
	}

	DB.AutoMigrate(&Message{})

	return
}

func Insert(value interface{}) error {
	return DB.Create(value).Error
}

func ReturnValues(where map[string]interface{}, dest interface{}) error {
	return DB.Where(where).Find(dest).Error
}
