package store

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitTables() {
	if db, driverErr := gorm.Open(sqlite.Open("dev.db"), &gorm.Config{}); driverErr == nil {
		DB = db

		db.AutoMigrate(&Track{})
		db.AutoMigrate(&Hotlap{})
	} else {
		log.Fatalf("Couldn't open connection to driver: %v\n", driverErr)
	}
}
