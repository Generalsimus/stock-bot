package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// type MarketPosition struct {
// 	afterPositionId int64
// 	open            float32
// 	close           float32
// 	gorm.Model
// }

var Database *gorm.DB

func Init() {
	log.Println("Connecting Database...")

	db, err := gorm.Open(sqlite.Open("marketDB.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err)
		os.Exit(2)
	}

	log.Println("Connected Successfully to Database")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")

	// db.AutoMigrate(&MarketPosition{})
	db.AutoMigrate(&Order{})
	Database = db

}