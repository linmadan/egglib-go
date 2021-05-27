package logrus

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"io"
	"os"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.WarnLevel)
}

type Logger struct {
	serviceName string
	logrus      *logrus.Logger
}

func (logger *Logger) SetServiceName(serviceName string) {
	logger.serviceName = serviceName
}

func (logger *Logger) SetLevel(level string) {
	switch level {
	case "trace":
		logger.logrus.Level = logrus.TraceLevel
	case "debug":
		logger.logrus.Level = logrus.DebugLevel
	case "info":
		logger.logrus.Level = logrus.InfoLevel
	case "warn":
		logger.logrus.Level = logrus.WarnLevel
	case "error":
		logger.logrus.Level = logrus.ErrorLevel
	case "fatal":
		logger.logrus.Level = logrus.FatalLevel
	case "panic":
		logger.logrus.Level = logrus.PanicLevel
	default:
		logger.logrus.Level = logrus.DebugLevel
	}
}

func (logger *Logger) AddHook(w io.Writer) {
	level := logger.logrus.Level
	var levels []logrus.Level
	// 默认已经添加了一个当前log level的hook,所以此处 level+1
	for i := level + 1; i <= logrus.TraceLevel; i++ {
		levels = append(levels, level)
	}
	logger.logrus.AddHook(&writer.Hook{
		Writer:    w,
		LogLevels: levels,
	})
}

func (logger *Logger) Trace(msg string, appends ...map[string]interface{}) {
	contextLogger := logger.logrus.WithFields(logrus.Fields{"serviceName": logger.serviceName})
	for _, append := range appends {
		contextLogger = contextLogger.WithFields(append)
	}
	contextLogger.Trace(msg)
}

func (logger *Logger) Debug(msg string, appends ...map[string]interface{}) {
	contextLogger := logger.logrus.WithFields(logrus.Fields{"serviceName": logger.serviceName})
	for _, append := range appends {
		contextLogger = contextLogger.WithFields(append)
	}
	contextLogger.Debug(msg)
}

func (logger *Logger) Info(msg string, appends ...map[string]interface{}) {
	contextLogger := logger.logrus.WithFields(logrus.Fields{"serviceName": logger.serviceName})
	for _, append := range appends {
		contextLogger = contextLogger.WithFields(append)
	}
	contextLogger.Info(msg)
}

func (logger *Logger) Warn(msg string, appends ...map[string]interface{}) {
	contextLogger := logger.logrus.WithFields(logrus.Fields{"serviceName": logger.serviceName})
	for _, append := range appends {
		contextLogger = contextLogger.WithFields(append)
	}
	contextLogger.Warn(msg)
}

func (logger *Logger) Error(msg string, appends ...map[string]interface{}) {
	contextLogger := logger.logrus.WithFields(logrus.Fields{"serviceName": logger.serviceName})
	for _, append := range appends {
		contextLogger = contextLogger.WithFields(append)
	}
	contextLogger.Error(msg)
}

func (logger *Logger) Fatal(msg string, appends ...map[string]interface{}) {
	contextLogger := logger.logrus.WithFields(logrus.Fields{"serviceName": logger.serviceName})
	for _, append := range appends {
		contextLogger = contextLogger.WithFields(append)
	}
	contextLogger.Fatal(msg)
}

func (logger *Logger) Panic(msg string, appends ...map[string]interface{}) {
	contextLogger := logger.logrus.WithFields(logrus.Fields{"serviceName": logger.serviceName})
	for _, append := range appends {
		contextLogger = contextLogger.WithFields(append)
	}
	contextLogger.Panic(msg)
}

func NewLogrusLogger() *Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout
	return &Logger{
		logrus: logger,
	}
}
