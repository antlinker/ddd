package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	_ ILogger = &logrusLog{}
)

func init() {
	logger = &logrusLog{
		Logger: logrus.Logger{
			Out:       os.Stderr,
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		},
	}

}

type logrusLog struct {
	logrus.Logger
}

func (l *logrusLog) InputLevel() Level {
	switch l.Logger.Level {
	case logrus.DebugLevel:
		return DebugLevel
	case logrus.InfoLevel:
		return InfoLevel
	case logrus.WarnLevel:
		return WarnLevel
	case logrus.ErrorLevel:
		return ErrorLevel
	case logrus.FatalLevel:
		return FatalLevel
	case logrus.PanicLevel:
		return FatalLevel
	}
	return DebugLevel
}
func (l *logrusLog) SetInputLevel(lvl Level) {
	switch lvl {
	case DebugLevel:
		l.Logger.Level = logrus.DebugLevel
	case InfoLevel:
		l.Logger.Level = logrus.InfoLevel
	case WarnLevel:
		l.Logger.Level = logrus.WarnLevel
	case ErrorLevel:
		l.Logger.Level = logrus.ErrorLevel
	case FatalLevel:
		l.Logger.Level = logrus.FatalLevel

	}

}
