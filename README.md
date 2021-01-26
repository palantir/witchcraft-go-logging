<p align="right">
<a href="https://autorelease.general.dmz.palantir.tech/palantir/witchcraft-go-logging"><img src="https://img.shields.io/badge/Perform%20an-Autorelease-success.svg" alt="Autorelease"></a>
</p>

witchcraft-go-logging
=====================
[![](https://godoc.org/github.com/palantir/witchcraft-go-logging?status.svg)](http://godoc.org/github.com/palantir/witchcraft-go-logging)

`witchcraft-go-logging` is a Go implementation of the Witchcraft logging specification. It provides an API that can be
used for logging and some default implementations of the logging API using different existing popular Go logging
libraries. `witchcraft-go-logging` includes implementations that use [zap](https://github.com/uber-go/zap), 
[zerolog](https://github.com/rs/zerolog) and [glog](https://github.com/golang/glog). We also provide an implementation
for [go-logr](https://github.com/go-logr/logr) that uses svc1log internally.

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

The default implementation of `wlog.DefaultLoggerProvider()` returns a logger that outputs a warning that states that
the default logger provider has not been set. For example, running the program:

```go
package main

import (
	"os"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

func main() {
	logger := svc1log.New(os.Stdout, wlog.InfoLevel)
	logger.Info("Hello")
}
```

Results in the following output to STDOUT:

```
[WARNING] Logging operation that uses the default logger provider was performed without specifying a logger provider implementation. To see logger output, set the global logger provider implementation using wlog.SetDefaultLoggerProvider or by importing an implementation. This warning can be disabled by setting the global logger provider to be the noop logger provider using wlog.SetDefaultLoggerProvider(wlog.NewNoopLoggerProvider()).
```

The global logger provider should always be set by the top-level program (the `main` package), so if this warning is
output it indicates that the top-level program should set a global logger implementation.

Most logger implementations have a package that contains an `init()` function that sets the global logger provider to be
the implementation's provider. This allows an underscore import to set the logger provider. For example, the following
sets the default logger provider to be a provider backed by `zap`:

```go
package main

import (
	"os"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	// import wlog-zap to set zap as the default logger provider
	_ "github.com/palantir/witchcraft-go-logging/wlog-zap"
)

func main() {
	logger := svc1log.New(os.Stdout, wlog.InfoLevel)
	logger.Info("Hello")
}
```

Running this program results in the following output to STDOUT:

```
{"level":"INFO","time":"2018-12-01T05:25:28.856348Z","message":"Hello","type":"service.1"}
```

It is also possible to set the default logger provider explicitly using the `SetDefaultLoggerProvider(provider LoggerProvider)`
function in `wlog`. For example, the following program also uses `wlog-zap` as the default logger provider, but does so
by calling `wlog.SetDefaultLoggerProvider(wlogzap.LoggerProvider())` rather than using an import: 

```go
package main

import (
	"os"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog-zap"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

func main() {
	wlog.SetDefaultLoggerProvider(wlogzap.LoggerProvider())
	logger := svc1log.New(os.Stdout, wlog.InfoLevel)
	logger.Info("Hello")
}
```

Setting the default logger provider to a no-op logger disables all output of loggers created using the default logger
provider. This can be done by calling `wlog.SetDefaultLoggerProvider(wlog.NewNoopLoggerProvider())`. 

Using loggers in code
---------------------
### Creating loggers and making them available in code
Each logger defines creation functions that typically take the `io.Writer` to which the logger output should be written
as an argument. There are typically 2 versions of a logger creation function: one that is explicitly provided with the
logger implementation that should be used to create the logger, and another which uses the default logger provider (as
determined at runtime) to create the logger. In most cases, loggers are created using the function that uses the default
logger provider (this makes it easier to set the default logger provider once to change all implementation).

Loggers are typically created by the top-level program (the program with a `main` package) and made available to other
code using some mechanism such as setting it in an exported package variable, passing it as an argument or setting it in
the contexts provided to program logic. This is a general issue common to all logging frameworks, and 
`witchcraft-go-logging` does not take an explicit stance on the correct approach.

Packages that are written as libraries typically do not instantiate loggers themselves -- they either accept the 
required loggers as arguments or have a context parameter and require that the expected loggers be set in the context.

### Using contexts to propagate loggers
Most logger packages define functions that can be used to set and retrieve the logger from a context. For example,
the `svc1log` package defines `WithLogger` and `FromContext` functions that can be used to set a logger on a context and
retrieve a logger from a context, respectively.

If a `FromContext` function is called on a context that does not have the logger set, it creates a default logger that
is returned instead. This ensures that the function will not return `nil`. However, this situation is usually indicative
of a programming error -- the consuming API expected a logger to be set on a context, but it was not (this implicit API
dependency is a [commonly expressed concern](https://dave.cheney.net/2017/01/26/context-is-for-cancelation) about using
storing loggers in contexts). As such, the default implementation of the logger returned in this situation is configured
to write to STDERR, and writes a warning about this situation (followed by the actual logger output). The logger 
returned by the `FromContext` function when no logger is present in the context is configurable, so if this default
behavior is not desirable it can be changed -- for example, one may return a noop logger to quietly suppress output or
return `nil` to force a panic in this situation.  

One advantage of using loggers stored in contexts is the ability to decorate them with parameters so that subsequent
calls use the provided parameters. For example, consider the following series of calls starting with `UpdateService`:

```go
func UpdateService(ctx context.Context, serviceID string) {
	ctx = svc1log.WithLoggerParams(ctx, svc1log.SafeParam("serviceId", serviceID))
	for _, currProcessID := range processIDs {
		updateProcess(ctx, currProcessID)
	}
}

func updateProcess(ctx context.Context, processID string) {
	ctx = svc1log.WithLoggerParams(ctx, svc1log.SafeParam("processId", processID))
	updateValue(ctx, processVals[processID], "timestamp", time.Now().String())
}

func updateValue(ctx context.Context, vals map[string]string, key, newValue string) {
	prevValue := vals[key]
	vals[key] = newValue
	svc1log.FromContext(ctx).Debug("Updating value", svc1log.SafeParam("prevValue", prevValue), svc1log.SafeParam("newValue", newValue))
}
```

In this series of calls, each function creates a new context that decorates its service logger with the provided
parameter. This has the result that, when `updateValue` performs its debug logging, the `serviceId` and `processId`
parameters that were added in the previous calls will be included in the logger output.

Active TODOs
------------
* Improve testing loggers that produce non-JSON output (glog)
* Port over more tests for audit logs

License
-------
This project is made available under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0).
