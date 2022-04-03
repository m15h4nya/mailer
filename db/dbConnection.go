package db

import (
	"apitask/config"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.DB), &gorm.Config{})
	if err != nil {
		fmt.Printf("[ERROR] GetConnection: %v", err)
	}
	return db
}
