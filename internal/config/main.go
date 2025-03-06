package config

import (
	"os"
	"time"
)

type Config struct {
	PollingInterval time.Duration
	HostPath        string
}

var configInstance *Config

func GetConfig() *Config {
	if configInstance == nil {
		pollingInterval := 60 * time.Second
		if interval := os.Getenv("POLLING_INTERVAL"); interval != "" {
			if d, err := time.ParseDuration(interval); err == nil {
				pollingInterval = d
			}
		}

		hostPath := "/"
		if path := os.Getenv("HOST_PATH"); path != "" {
			hostPath = path
		}

		configInstance = &Config{
			PollingInterval: pollingInterval,
			HostPath:        hostPath,
		}
	}

	return configInstance
}
