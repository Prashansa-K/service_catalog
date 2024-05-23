package config

import (
	"fmt"

	utils "github.com/Prashansa-K/serviceCatalog/internal"
)

const (
	DEFAULT_METRICS_HOST = "localhost"
	DEFAULT_METRICS_PORT = "8081"
)

type MetricsServerConfig struct {
	Host    string
	Port    string
	Address string
}

func GetMetricsServerConfig() *MetricsServerConfig {
	host := utils.GetEnvWithDefault("METRICS_HOST", DEFAULT_METRICS_HOST)
	port := utils.GetEnvWithDefault("METRICS_PORT", DEFAULT_METRICS_PORT)
	return &MetricsServerConfig{
		Host:    host,
		Port:    port,
		Address: fmt.Sprintf("%s:%s", host, port),
	}

}
