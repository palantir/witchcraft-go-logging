package svc1log

import (
	"github.com/palantir/witchcraft-go-logging/wlog"
)

type multiLogger struct {
	loggers []Logger
}

// NewMultiLogger returns a Logger that logs to all the provided loggers.
func NewMultiLogger(loggers ...Logger) Logger {
	return &multiLogger{
		loggers: loggers,
	}
}

func (m multiLogger) Debug(msg string, params ...Param) {
	for _, logger := range m.loggers {
		logger.Debug(msg, params...)
	}
}

func (m multiLogger) Info(msg string, params ...Param) {
	for _, logger := range m.loggers {
		logger.Info(msg, params...)
	}
}

func (m multiLogger) Warn(msg string, params ...Param) {
	for _, logger := range m.loggers {
		logger.Warn(msg, params...)
	}
}

func (m multiLogger) Error(msg string, params ...Param) {
	for _, logger := range m.loggers {
		logger.Error(msg, params...)
	}
}

func (m multiLogger) SetLevel(level wlog.LogLevel) {
	for _, logger := range m.loggers {
		logger.SetLevel(level)
	}
}

func (m multiLogger) Enabled(level wlog.LogLevel) bool {
	for _, logger := range m.loggers {
		if l, ok := logger.(wlog.LevelChecker); ok {
			if l.Enabled(level) {
				return true
			}
		} else {
			return true
		}
	}
	// All loggers implement LevelChecker and none of them are enabled
	return false
}
