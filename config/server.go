package config

import (
	"fmt"

	utils "github.com/Prashansa-K/serviceCatalog/internal"
)

const (
	DEFAULT_HOST = "localhost"
	DEFAULT_PORT = "8080"
)

type ServerConfig struct {
	Host    string
	Port    string
	Address string
}

func GetServerConfig() *ServerConfig {
	host := utils.GetEnvWithDefault("SERVICE_HOST", DEFAULT_HOST)
	port := utils.GetEnvWithDefault("SERVICE_PORT", DEFAULT_PORT)

	return &ServerConfig{
		Host:    host,
		Port:    port,
		Address: fmt.Sprintf("%s:%s", host, port),
	}

}
