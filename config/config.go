package config

import (
	"os"
	"sync"
)

var instance *Config
var once sync.Once

type Config struct {
	Listen     string
	LogLevel   string
	TimeOut    string
	DbUser     string
	DbPassword string
	DbListen   string
	Database   string
}

func Load() *Config {
	once.Do(func() {
		var listen, logLevel, timeout, user, passwd, db, dbListen string
		if listen = os.Getenv("LISTEN"); listen == "" {
			listen = "localhost:8080"
		}
		if logLevel = os.Getenv("LOG_LEVEL"); logLevel == "" {
			logLevel = "info"
		}
		if timeout = os.Getenv("TIMEOUT"); timeout == "" {
			timeout = "10ms"
		}
		if user = os.Getenv("DB_USERNAME"); user == "" {
			user = "admin"
		}
		if db = os.Getenv("DATABASE"); db == "" {
			db = "test"
		}
		if passwd = os.Getenv("DB_PASSWORD"); passwd == "" {
			passwd = ""
		}
		if dbListen = os.Getenv("DB_INTERFACE"); dbListen == "" {
			dbListen = "127.0.0.1:3301"
		}
		instance = &Config{
			Listen:     listen,
			LogLevel:   logLevel,
			TimeOut:    timeout,
			DbUser:     user,
			DbPassword: passwd,
			DbListen:   dbListen,
			Database:   db,
		}
	})
	return instance
}
