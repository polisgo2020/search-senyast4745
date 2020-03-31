package util

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/polisgo2020/search-senyast4745/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

type appYamlConfig struct {
	Logger []log.Config `yaml:"logger"`
	Server struct {
		Port    string `yaml:"port"`
		Timeout int    `yaml:"timeout"`
	}
}

type App struct {
	Mux  *chi.Mux
	Port string
}

func NewAppFromConfig(configFile string) (*App, error) {
	dat, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	conf := appYamlConfig{}

	if err = yaml.Unmarshal(dat, &conf); err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	r.Use(middleware.Timeout(time.Duration(conf.Server.Timeout) * time.Millisecond))

	log.GetLogger(conf.Logger...)
	return &App{Mux: r, Port: conf.Server.Port}, nil
}

func NewDefaultApp() (*App, error) {

	r := chi.NewMux()
	log.GetLogger()
	return &App{Mux: r, Port: "8080"}, nil
}

func NewApp() (*App, error) {
	if _, err := os.Stat("config.yml"); err == nil {
		return NewAppFromConfig("config.yml")
	} else {
		return NewDefaultApp()
	}
}
