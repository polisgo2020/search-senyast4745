package log

import (
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
)

const (
	json   = "json"
	logfmt = "logfmt"

	console = "console"
	file    = "file"

	DebugLvl   = "debug"
	InfoLvl    = "info"
	WarningLvl = "warning"
	ErrorLvl   = "error"
)

type Config struct {
	Name     string `yaml:"name"`
	LogLevel string `yaml:"log-level"`
	Format   string `yaml:"format"`
	Output   struct {
		Type     string
		Rollback lumberjack.Logger
	}
}

type customLogger struct {
	loggers        []*kitlog.Logger
	consoleLoggers []*kitlog.Logger
}

var instance *customLogger
var once sync.Once

func initDefaultLogger() *customLogger {
	logger := &customLogger{loggers: make([]*kitlog.Logger, 0, 0),
		consoleLoggers: make([]*kitlog.Logger, 0, 1)}

	log := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	log = level.NewFilter(log, level.AllowInfo())

	logger.consoleLoggers = append(logger.consoleLoggers, &log)

	return instance
}

func GetLogger(configs ...Config) *customLogger {
	once.Do(func() {

		if len(configs) < 1 {
			instance = initDefaultLogger()
			return
		}

		logger := &customLogger{loggers: make([]*kitlog.Logger, 0, len(configs)),
			consoleLoggers: make([]*kitlog.Logger, 0, len(configs))}
		for _, conf := range configs {

			var log kitlog.Logger
			var w io.Writer
			switch conf.Output.Type {
			case console:
				w = os.Stdout
				logger.consoleLoggers = append(logger.consoleLoggers, &log)
			case file:
				w = &conf.Output.Rollback
				logger.loggers = append(logger.loggers, &log)
			default:
				instance = initDefaultLogger()
				return
			}

			switch conf.Format {
			case json:
				log = kitlog.NewJSONLogger(kitlog.NewSyncWriter(w))
			case logfmt:
				log = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(w))
			default:
				instance = initDefaultLogger()
				return
			}

			switch conf.LogLevel {
			case DebugLvl:
				log = level.NewFilter(log, level.AllowDebug())
			case InfoLvl:
				log = level.NewFilter(log, level.AllowInfo())
			case WarningLvl:
				log = level.NewFilter(log, level.AllowWarn())
			case ErrorLvl:
				log = level.NewFilter(log, level.AllowError())
			default:
				fmt.Printf("incorrect log-level %s, using info \n", conf.LogLevel)
				log = level.NewFilter(log, level.AllowInfo())
			}

			log = kitlog.With(log, "log-name", conf.Name)

		}
		instance = logger
	})

	return instance
}

type logLevel func(kitlog.Logger) kitlog.Logger

func Debug(v ...interface{}) {
	GetLogger().loggingWithLevel(level.Debug, v)
}

func Info(v ...interface{}) {
	GetLogger().loggingWithLevel(level.Info, v)
}

func Warn(v ...interface{}) {
	GetLogger().loggingWithLevel(level.Warn, v)
}

func Error(v ...interface{}) {
	GetLogger().loggingWithLevel(level.Error, v)
}

func (log *customLogger) loggingWithLevel(levelFun logLevel, v interface{}) {
	readSyn := sync.RWMutex{}

	info := v.([]interface{})
	wg := sync.WaitGroup{}
	conLoggers := log.consoleLoggers

	for i := range conLoggers {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			readSyn.RLock()
			defer readSyn.RUnlock()
			levelFun(*conLoggers[i]).Log(info...)
		}(&wg)
	}

	loggers := log.loggers
	for i := range loggers {
		go func() {
			readSyn.RLock()
			defer readSyn.RUnlock()
			levelFun(*loggers[i]).Log(info...)
		}()
	}
	wg.Wait()
}
