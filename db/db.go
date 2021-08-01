package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() {
	DB, _ = gorm.Open(sqlite.Open("data/gorm.db"), &gorm.Config{})
}
