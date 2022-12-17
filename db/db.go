package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// var Database *gorm.DB

func GetDb() *gorm.DB {
	fmt.Println("Connecting Database...")

	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err)
		os.Exit(2)
	}

	fmt.Println("Connected Successfully to Database")
	fmt.Println("Running Migrations")

	db.AutoMigrate(&AlpacaOrder{})
	db.AutoMigrate(&Bar{})
	return db
}
