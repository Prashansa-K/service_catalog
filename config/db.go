package config

import (
	"fmt"

	utils "github.com/Prashansa-K/serviceCatalog/internal"
)

const (
	DEFAULT_DB_HOST     = "localhost"
	DEFAULT_DB_PORT     = "5432"
	DEFAULT_DB_USER     = "postgres"
	DEFAULT_DB_PASSWORD = ""
	DEFAULT_DB_NAME     = "servicecatalog"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	DSN      string
}

func GetDBConfig() *DBConfig {
	host := utils.GetEnvWithDefault("DB_HOST", DEFAULT_DB_HOST)
	port := utils.GetEnvWithDefault("DB_PORT", DEFAULT_DB_PORT)
	user := utils.GetEnvWithDefault("DB_USER", DEFAULT_DB_USER)
	password := utils.GetEnvWithDefault("DB_PASSWORD", DEFAULT_DB_PASSWORD)
	dbName := utils.GetEnvWithDefault("DB_NAME", DEFAULT_DB_NAME)

	return &DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		DSN: fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbName,
		),
	}
}
