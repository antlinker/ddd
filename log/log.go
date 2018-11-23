package log

import (
	"fmt"
	"strings"
)

// ILogger 日志接口
type ILogger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(fomat string, v ...interface{})
	InputLevel() Level
	SetInputLevel(l Level)
}

// Level 日志级别
type Level uint32

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	}
	return "unknown"
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid logrus Level: %q", lvl)
}

const (
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel Level = iota
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
)

var (
	logger ILogger
	// logger ILogger = &logrus.Logger{
	// 	Out:       os.Stderr,
	// 	Formatter: new(logrus.JSONFormatter),
	// 	Hooks:     make(logrus.LevelHooks),
	// 	Level:     logrus.DebugLevel,
	// }
)

// Logger 日志
func Logger() ILogger {
	return logger
}

// SetLogger 设置日志处理器
func SetLogger(log ILogger) {
	logger = log
}

// Debug 输出日志Debug级别
func Debug(v ...interface{}) {
	logger.Debug(v...)
}

// Debugf 输出日志Debugf级别
func Debugf(fomat string, v ...interface{}) {
	logger.Debugf(fomat, v...)
}

// Info 输出日志Info级别
func Info(v ...interface{}) {
	logger.Info(v...)
}

// Infof 输出日志Info级别
func Infof(fomat string, v ...interface{}) {
	logger.Infof(fomat, v...)
}

// Error 输出日志Error级别
func Error(v ...interface{}) {
	logger.Error(v...)
}

// Errorf 输出日志Error级别
func Errorf(fomat string, v ...interface{}) {
	logger.Errorf(fomat, v...)
}

// Warn 输出日志Warn级别
func Warn(v ...interface{}) {
	logger.Warn(v...)
}

// Warnf 输出日志Warn级别
func Warnf(fomat string, v ...interface{}) {
	logger.Warnf(fomat, v...)
}

// Fatal 输出日志Fatal级别
func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

// Fatalf 输出日志Fatal级别
func Fatalf(fomat string, v ...interface{}) {
	logger.Fatalf(fomat, v...)
}
