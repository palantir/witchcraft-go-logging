package wloginternal

import (
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

type Param[T logging.LogTypes] interface {
	apply(*T)
}

// ParamFunc is a function that implements Param and modifies a Conjure log object.
type ParamFunc[T logging.LogTypes] func(log *T)

func (f ParamFunc[T]) apply(log *T) {
	f(log)
}

func ApplyParams[T logging.LogTypes](log *T, params ...Param[T]) {
	for _, p := range params {
		if p != nil {
			p.apply(log)
		}
	}
}

func LogParams[T logging.LogTypes](logger func(*T), params ...Param[T]) {
	log, release := GetPooled[T]()
	ApplyParams(log, params...)
	logger(log)
	release(log)
}
