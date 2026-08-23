package config

import (
	"os"
	"time"
)

type Config struct {
	Address         string
	ShutdownTimeout time.Duration
	LogLevel        string
	Driver          string
}

func Load() Config {
	c := Config{Address: ":8084", ShutdownTimeout: 10 * time.Second, LogLevel: "info", Driver: "memory"}
	if v := os.Getenv("API_QUOTA_ADDRESS"); v != "" {
		c.Address = v
	}
	if v := os.Getenv("API_QUOTA_SHUTDOWN_TIMEOUT"); v != "" {
		if d, e := time.ParseDuration(v); e == nil {
			c.ShutdownTimeout = d
		}
	}
	if v := os.Getenv("API_QUOTA_LOG_LEVEL"); v != "" {
		c.LogLevel = v
	}
	if v := os.Getenv("API_QUOTA_STORAGE_DRIVER"); v != "" {
		c.Driver = v
	}
	return c
}
