package log

import (
	"bytes"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
	"testing"
)

var (
	consoleConf = &Config{
		Name:     "test",
		LogLevel: "debug",
		Format:   "logfmt",
		Output: struct {
			Type     string
			Rollback lumberjack.Logger
		}{"console", lumberjack.Logger{}},
	}

	fileConf = &Config{
		Name:     "file-test",
		LogLevel: "error",
		Format:   "json",
		Output: struct {
			Type     string
			Rollback lumberjack.Logger
		}{"file", lumberjack.Logger{
			Filename:   "test-log.log",
			MaxSize:    10,
			MaxAge:     1,
			MaxBackups: 1,
			LocalTime:  false,
			Compress:   true,
		}},
	}
)

func TestGetLogger(t *testing.T) {
	once = *new(sync.Once)
	clogger := &customLogger{loggers: make([]*log.Logger, 0, 0),
		consoleLoggers: make([]*log.Logger, 0, 1)}

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowInfo())

	clogger.consoleLoggers = append(clogger.consoleLoggers, &logger)
	logger = log.With(logger, "log-name", "default")

	require.Equal(t, clogger, GetLogger())
}

func TestGetLogger2(t *testing.T) {
	once = *new(sync.Once)
	require.Equal(t, 1, len(GetLogger(*consoleConf, *fileConf).consoleLoggers))
	require.Equal(t, 1, len(GetLogger(*consoleConf, *fileConf).loggers))

}

func TestGetLogger3(t *testing.T) {
	once = *new(sync.Once)

	clogger := &customLogger{loggers: make([]*log.Logger, 0, 0),
		consoleLoggers: make([]*log.Logger, 0, 1)}

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowInfo())

	clogger.consoleLoggers = append(clogger.consoleLoggers, &logger)
	logger = log.With(logger, "log-name", "default")

	consoleConf.Output.Type = "incorrect"
	require.Equal(t, clogger, GetLogger(*consoleConf))

	once = *new(sync.Once)
	consoleConf.Format = "incorrect"
	require.Equal(t, clogger, GetLogger(*consoleConf))
}

type logTestSuite struct {
	suite.Suite
	logger *customLogger
	output *bytes.Buffer
}

func (i *logTestSuite) SetupSuite() {
	once.Do(func() {
	})
	i.output = &bytes.Buffer{}
}

func (i *logTestSuite) SetupTest() {
	clogger := &customLogger{loggers: make([]*log.Logger, 0, 0),
		consoleLoggers: make([]*log.Logger, 0, 1)}

	logger := log.NewLogfmtLogger(log.NewSyncWriter(i.output))
	logger = level.NewFilter(logger, level.AllowDebug())

	clogger.consoleLoggers = append(clogger.consoleLoggers, &logger)
	logger = log.With(logger, "log-name", "test")

	i.logger = clogger
	instance = clogger
}

func (i *logTestSuite) TearDownTest() {
	i.logger = nil
	instance = nil
	i.output.Reset()
}
func (i *logTestSuite) TearDownSuite() {
	once = *new(sync.Once)
}

func TestSearchSuitStart(t *testing.T) {
	suite.Run(t, new(logTestSuite))
}

func (i *logTestSuite) TestDebug() {
	Debug("msg", "hello")
	require.Contains(i.T(), i.output.String(), "level=debug")
	require.Contains(i.T(), i.output.String(), "msg=hello")
	*instance.consoleLoggers[0] = level.NewFilter(*instance.consoleLoggers[0], level.AllowInfo())
	i.output.Reset()
	Debug("msg", "hello")
	require.Empty(i.T(), i.output.Bytes())
}

func (i *logTestSuite) TestInfo() {
	Info("msg", "hello")
	require.Contains(i.T(), i.output.String(), "level=info")
	require.Contains(i.T(), i.output.String(), "msg=hello")
	require.NotContains(i.T(), i.output.String(), "test-key=test-val")
	*instance.consoleLoggers[0] = log.With(*instance.consoleLoggers[0], "test-key", "test-val")
	i.output.Reset()
	Info("msg", "hello")
	require.Contains(i.T(), i.output.String(), "test-key=test-val")
}

func (i *logTestSuite) TestWarn() {
	Warn("msg", "hello")
	require.Contains(i.T(), i.output.String(), "level=warn")
	require.Contains(i.T(), i.output.String(), "msg=hello")

}

func (i *logTestSuite) TestError() {
	Error("msg", "hello")
	require.Contains(i.T(), i.output.String(), "level=error")
	require.Contains(i.T(), i.output.String(), "msg=hello")
	i.output.Reset()
	Error("err", errors.New("test error"))
	require.Contains(i.T(), i.output.String(), "err=\"test error\"")
}

func (i *logTestSuite) TestManyLoggers() {

	output2 := &bytes.Buffer{}

	logger := log.NewLogfmtLogger(log.NewSyncWriter(output2))
	logger = level.NewFilter(logger, level.AllowDebug())

	i.logger.consoleLoggers = append(i.logger.consoleLoggers, &logger)
	instance = i.logger
	logger = log.With(logger, "log-name", "test-2")

	output3 := &bytes.Buffer{}

	logger3 := log.NewLogfmtLogger(log.NewSyncWriter(output3))
	logger3 = level.NewFilter(logger3, level.AllowDebug())

	i.logger.consoleLoggers = append(i.logger.consoleLoggers, &logger3)

	Info("msg", "hello")
	require.Contains(i.T(), i.output.String(), "level=info")
	require.Contains(i.T(), i.output.String(), "msg=hello")
	require.Contains(i.T(), output2.String(), "level=info")
	require.Contains(i.T(), output2.String(), "msg=hello")
	require.Contains(i.T(), output3.String(), "msg=hello")
	require.Contains(i.T(), output3.String(), "msg=hello")
}
