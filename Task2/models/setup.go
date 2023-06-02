package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Output to console
		logger.Config{
			SlowThreshold: time.Second, // Customize the slow threshold if desired
			LogLevel:      logger.Info, // Set log level as desired
			Colorful:      true,        // Enable colored output (optional)
		},
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("Failed to connect to database")
	}


	database.AutoMigrate(&Spot{})


	DB = database
}
