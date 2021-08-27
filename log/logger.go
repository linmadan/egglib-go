package log

import "io"

type Logger interface {
	SetServiceName(serviceName string)
	SetLevel(level string)
	AddHook(write io.Writer)
	Trace(msg string, appends ...map[string]interface{})
	Debug(msg string, appends ...map[string]interface{})
	Info(msg string, appends ...map[string]interface{})
	Warn(msg string, appends ...map[string]interface{})
	Error(msg string, appends ...map[string]interface{})
	Fatal(msg string, appends ...map[string]interface{})
	Panic(msg string, appends ...map[string]interface{})
}
