package config

import (
	"os"
	"sync"
)

var instance *Config
var once sync.Once

type Config struct {
	Listen   string
	LogLevel string
	TimeOut  string
}

func Load() *Config {
	once.Do(func() {
		var listen, logLevel, timeout string
		if listen = os.Getenv("LISTEN"); listen == "" {
			listen = "localhost:8080"
		}
		if logLevel = os.Getenv("LOG_LEVEL"); logLevel == "" {
			logLevel = "info"
		}
		if timeout = os.Getenv("TIMEOUT"); timeout == "" {
			timeout = "1"
		}
		instance = &Config{
			Listen:   listen,
			LogLevel: logLevel,
			TimeOut:  timeout,
		}
	})
	return instance
}
