package gologrwlog

import (
	"fmt"
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type goLogrWrapper struct {
	origin string
	logger svc1log.Logger
}

// NewGoLogrWrapper returns a go-logr Logger implementation that uses svc1log to write logs.
func NewGoLogrWrapper(logger svc1log.Logger, origin string) logr.Logger {
	logger = svc1log.WithParams(logger, svc1log.Origin(origin))
	return &goLogrWrapper{
		origin: origin,
		logger: logger,
	}
}

func (s *goLogrWrapper) Info(msg string, keysAndValues ...interface{}) {
	s.logger.Info(msg, toSafeParams(s.logger, keysAndValues))
}

func (s *goLogrWrapper) Enabled() bool {
	return true
}

func (s *goLogrWrapper) Error(err error, msg string, keysAndValues ...interface{}) {
	s.logger.Error(msg, svc1log.Stacktrace(err), toSafeParams(s.logger, keysAndValues))
}

func (s *goLogrWrapper) V(level int) logr.InfoLogger {
	return NewGoLogrWrapper(s.logger, s.origin)
}

func (s *goLogrWrapper) WithValues(keysAndValues ...interface{}) logr.Logger {
	logger := svc1log.WithParams(s.logger, toSafeParams(s.logger, keysAndValues))
	return NewGoLogrWrapper(logger, s.origin)
}

func (s *goLogrWrapper) WithName(name string) logr.Logger {
	return NewGoLogrWrapper(s.logger, filepath.Join(s.origin, name))
}

func toSafeParams(logger svc1log.Logger, keysAndValues []interface{}) svc1log.Param {
	if len(keysAndValues)%2 == 1 {
		logger.Error("KeysAndValues pair slice has an odd number of arguments; ignoring all",
			svc1log.SafeParam("keysAndValuesLen", len(keysAndValues)))
		return svc1log.SafeParams(map[string]interface{}{})
	}

	params := map[string]interface{}{}
	for i := 0; i < len(keysAndValues); i = i + 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			logger.Error("Key type is not string",
				svc1log.SafeParam("actualType", fmt.Sprintf("%T", keysAndValues[i])),
				svc1log.SafeParam("key", key))
			continue
		}
		params[key] = keysAndValues[i+1]
	}
	return svc1log.SafeParams(params)
}
