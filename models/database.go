package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("cronic.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database file: %w", err)
	}
	return DB.AutoMigrate(&Job{})
}
