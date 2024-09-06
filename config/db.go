package config

import (
	"fmt"
	"log"
	"os"

	"app/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("connected to database")
	db.Logger = logger.Default.LogMode(logger.Info)

	if migrate := os.Getenv("MIGRATE"); migrate == "TRUE" {
		log.Println("running migrations")

		err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
		if err != nil {
			log.Fatalf("failed to create extension: %v", err)
		}

		if err := db.AutoMigrate(&model.User{}); err != nil {
			log.Fatalf("failed to perform auto migration: %v", err)
		}
	}

	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	sqlDB.Close()
}
