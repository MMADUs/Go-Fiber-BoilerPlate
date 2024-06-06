package config

import (
	"fmt"
	"log"
	"os"

	"app2/model"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDB() error {
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
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("connected to database")
	db.Logger = logger.Default.LogMode(logger.Info)

	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		return fmt.Errorf("failed to create extension: %v", err)
	}

	if migrate := os.Getenv("MIGRATE"); migrate == "TRUE" {
		log.Println("running migrations")
		db.AutoMigrate(&model.Category{}, &model.Product{}, &model.User{})
	}

	DB = Dbinstance{
		Db: db,
	}

	return nil
}

func GetDB() *gorm.DB {
	return DB.Db
}

func CloseDB() {
	db := GetDB()

	if db == nil {
		log.Println("No database connection to close")
		return
	}

	DbConnection, err := db.DB()
	if err != nil {
		log.Println("Error getting DB from gorm:", err)
		return
	}

	err = DbConnection.Close()
	if err != nil {
		log.Println("Error closing DB:", err)
		return
	}
}