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
		instance = &Config{
			Listen:   os.Getenv("LISTEN"),
			LogLevel: os.Getenv("LOG_LEVEL"),
			TimeOut:  os.Getenv("TIMEOUT"),
		}
	})
	return instance
}
