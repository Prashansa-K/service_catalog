package db

import (
	"log"

	"github.com/Prashansa-K/serviceCatalog/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func connect() error {
	// Get the database configuration from the config.go file
	dbConfig := config.GetDBConfig()

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dbConfig.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)

		return err
	}

	DB = db

	return nil
}

func isConnected() bool {
	return DB != nil
}

func GetDB() (*gorm.DB, error) {
	if isConnected() {
		return DB, nil
	}

	err := connect()

	if err != nil {
		return nil, err
	}

	return DB, nil
}
