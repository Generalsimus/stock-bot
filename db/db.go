package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// var Database *gorm.DB

func GetDb() *gorm.DB {
	log.Println("Connecting Database...")

	db, err := gorm.Open(sqlite.Open("marketDB.db"), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err)
		os.Exit(2)
	}

	log.Println("Connected Successfully to Database")
	log.Println("Running Migrations")

	db.AutoMigrate(&Order{})
	db.AutoMigrate(&Bar{})
	return db
}
