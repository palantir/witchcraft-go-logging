witchcraft-go-logging
=====================
`witchcraft-go-logging` is a Go implementation of the Witchcraft logging specification. It provides an API that can be
used for logging and some default implementations of the logging API using different existing popular Go logging
libraries. `witchcraft-go-logging` includes implementations that use [zap](https://github.com/uber-go/zap), 
[zerolog](https://github.com/rs/zerolog) and [glog](https://github.com/golang/glog).

Architecture
------------
`witchcraft-go-logging` defines versioned logger interfaces for specific logger types (service.1 loggers, request.2
loggers, etc.) and provides implementations of those interfaces. The logger packages define functions for instantiating
these loggers and often provide functions for creating parameters for the loggers and for storing and retrieving loggers 
from a `context.Context`.

The loggers are implemented using abstractions defined in the `wlog` package -- specifically, `wlog.LogEntry`,
`wlog.Logger` and `wlog.LeveledLogger`.

`wlog.LogEntry` is an interface that represents a log entry, and offers functions for appending typed key-value pairs to
an entry. It also provides an `ObjectValue` function which can append values of arbitrary types.

The `wlog.Logger` interface defines the `Log(params ...Param)` function, where a `Param` is a functional parameter that
typically operates on a `wlog.LogEntry` by calling set functions on it. Conceptually, the `Log` function applies all of
the append operations specified by the `Param` arguments to an internal `wlog.LogEntry` object and outputs the result.

The `wlog.LeveledLogger` is similar to the `wlog.Logger` interface, but rather than having a `Log` function it declares
functions for logging at specific levels with a message (`Debug(msg string, params ...Param)`, 
`Info(msg string, params ...Param)`, etc.) and defines a `SetLevel(level LogLevel)` function that can be used to 
configure the level of the logger. 

The `wlog.LoggerCreator` and `wlog.LevelLoggerCreator` function types (defined as `type LoggerCreator func(w io.Writer) Logger` 
and `type LeveledLoggerCreator func(w io.Writer, level LogLevel) LeveledLogger`, respectively) define signatures for
creating a `wlog.Logger` or `wlog.LeveledLogger` given the required parameters.

The specific logger types define an instantiation function that takes in one of the logger creator types defined above 
as an argument and instantiates a typed logger that is backed by the logger implementation returned by the creator. For 
example, `metric1log` defines `func NewFromCreator(w io.Writer, creator wlog.LoggerCreator) Logger`, which creates a new
`metric.1` logger using the provided creator function that writes to the specified output. Logger types also define an
instantiation function that does not require specifying a creator -- these functions use the logger creator supplied by
the globally defined default logger instead. For example, `metric1log` defines `func New(w io.Writer) Logger`.

Set the default logger provider
-------------------------------
In the canonical usage pattern for loggers, loggers are instantiated using the version of the function that does not
specify a logger implementation -- for example, `metric1log.New(w io.Writer) Logger`.

These functions use the `wlog.DefaultLoggerProvider()` function to get the logger creator required to instantiate the
logger. This function returns a `wlog.LoggerProvider`, which is defined as:

```go
type LoggerProvider interface {
	NewLogger(w io.Writer) Logger
	NewLeveledLogger(w io.Writer, level LogLevel) LeveledLogger
}
```

The default implementation of `wlog.DefaultLoggerProvider()` returns a no-op logger provider, meaning that all logging
operations performed by loggers that use this implementation will be no-ops.

The `wlog.SetDefaultLoggerProvider(provider LoggerProvider)` can be used to set the default logger provider to be a 
specific logger provider. For example, calling `wlog.SetDefaultLoggerProvider(wlogzap.LoggerProvider())` will set the
default logger provider to be the logger provider provided by the `wlogzap` package, which is an implementation of the
`wlog` logger interfaces backed by the `zap` logger.

Logger provider implementations generally declare an `init()` function in their top-level package that calls
`wlog.SetDefaultLoggerProvider` and provides the implemented logger provider. This means that underscore imports can be
used to set the default logger provider: for example, `import _ "github.palantir.build/deployability/witchcraft-go-logging/witchcraft-go-logging-zap"`
will set the default logger provider to be `witchcraft-go-logging-zap`.   

Active TODOs
------------
* Add section on usage to README
* Improve testing loggers that produce non-JSON output (glog)
* Port over more tests for audit logs

License
-------
This project is made available under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0).
