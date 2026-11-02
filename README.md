<p align="right">
<a href="https://autorelease.general.dmz.palantir.tech/palantir/witchcraft-go-logging"><img src="https://img.shields.io/badge/Perform%20an-Autorelease-success.svg" alt="Autorelease"></a>
</p>

witchcraft-go-logging
=====================
[![](https://godoc.org/github.com/palantir/witchcraft-go-logging?status.svg)](http://godoc.org/github.com/palantir/witchcraft-go-logging)

`witchcraft-go-logging` is a Go implementation of the Witchcraft logging specification. It provides an API that can be
used for logging along with several implementations and adapters.

**Implementations** of the wlog interfaces include:

- The built-in object JSON provider, which encodes generated log objects with JSON v2 without zap or zerolog
- [zap](https://github.com/uber-go/zap) via [wlog-zap](wlog-zap)
- [zerolog](https://github.com/rs/zerolog) via [wlog-zerolog](wlog-zerolog)
- [glog](https://github.com/golang/glog) via [wlog-glog](wlog-glog)
- [wlog-tmpl](wlog-tmpl) for rendering structured logging using human-friendly templates.

**Adapters** wrap the witchcraft-go-logging logger implementations (svc1log, ev2log, req2log, etc) to allow interoperability with other Go logging interfaces. We currently provide
- [svc1zap](adapters/svc1zap) wraps a svc1log.Logger to provide a [zap](https://github.com/uber-go/zap) Logger.

Architecture
------------
`witchcraft-go-logging` defines versioned logger interfaces for specific logger types (service.1 loggers, request.2
loggers, etc.) and provides implementations of those interfaces. The logger packages define functions for instantiating
these loggers and often provide functions for creating parameters for the loggers and for storing and retrieving loggers 
from a `context.Context`.

Each typed logger applies its parameters to a generated Conjure log object. For example, `svc1log.SafeParam` and
`svc1log.UnsafeParam` accept values of arbitrary Go types. A provider with a method such as
`LogServiceLogV1(*logging.ServiceLogV1)` can encode that object directly. The built-in object JSON provider supports
typed methods for the log types and encodes them with JSON v2 so struct fields appear in declaration order. Optional
empty fields are omitted, and required empty collections are encoded as `[]` or `{}`. Audit.3 falls back to entry
emission when its generated object cannot preserve a missing or invalid event ID or a noncanonical log-entry ID.

Typed provider methods receive borrowed objects. Implementations must consume them synchronously or copy any object
and nested data retained after returning; typed loggers reset and reuse their objects immediately after each call.

Providers without a typed method use the existing `wlog.LogEntry` path. The typed logger converts its object to a
`wlog.Param`, which applies fields to the provider's entry. Direct calls to `wlog.Logger.Log(params ...wlog.Param)` also
use this path. `wlog.LogEntry`, `wlog.Logger`, and `wlog.LeveledLogger` remain the provider contract for zap, zerolog,
glog, templates, and custom providers.

`wlog.LoggerCreator` and `wlog.LeveledLoggerCreator` create provider loggers for an `io.Writer` (and a level for the
latter). Typed packages accept these creators in `NewFromCreator` and provide `New` to use the global default provider.
Both forms pass the caller's writer to the provider.

Default logger provider
-----------------------
Typed packages use `wlog.DefaultLoggerProvider()` when constructed with `New`. The built-in default is
`wlog.NewJSONMarshalLoggerProvider()`, which writes JSON to the `io.Writer` passed by the caller. Its interface is:

```go
type LoggerProvider interface {
	NewLogger(w io.Writer) Logger
	NewLeveledLogger(w io.Writer, level LogLevel) LeveledLogger
}
```

For example:

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

This writes a line like the following to STDOUT:

```
{"type":"service.1","time":"2018-12-01T05:25:28.856348Z","level":"INFO","message":"Hello"}
```

Programs can override the default with `wlog.SetDefaultLoggerProvider`, pass a provider to `NewFromCreator`, or import a
provider package for its `init()` registration. The built-in provider encodes typed log objects directly and uses maps
for generic entry calls. To disable output, use `wlog.SetDefaultLoggerProvider(wlog.NewNoopLoggerProvider())`.

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

Updating audit logger category definitions
------------------------------------------
Audit logger category definitions come from 2 different places:
* Conjure objects generated from the `audit-api` Conjure definition
* Generated Go code that defines the request and result fields and translates an `AuditCategoryV2` object to an audit
  parameter

The Conjure objects come from the Conjure definitions, while the generated Go code is based on Java code that defines
audit categories.

In order to update audit logger category definitions, it is necessary to update the Conjure definition and to generate
the Go code from the latest Java code.

### Updating the Conjure definition
Audit logger categories are defined by the "audit-api" Conjure definition. The Conjure definition is published
internally. The definition can be updated in this repository by following these steps:
  1. Obtain the latest version of the Conjure definition (the `.conjure.json` file)
  2. Remove the existing definition from the `conjure/defs` directory
  3. Place the new definition in the `conjure/defs` directory
  4. Update the Conjure definition in the `godel/config/conjure-plugin.yml` file to point to the new definition
  5. Run `./godelw conjure` to generate the code for the new definition

### Updating the generated Go code
The generated Go code is based on the Java code in the `AuditCategoryConfiguration.java` in the `foundry-audit` project.
This is an internal project. The generated go code can be updated as follows:
  1. Obtain the latest version of the Java code and note the path to the "AuditCategoryConfiguration.java" file
  2. Run the `[category_transform.go](wlog/auditlog/audit3log/internal/category_transform.go)` file as a main program
     with the path to the "AuditCategoryConfiguration.java" as the argument. The working directory should be the
     `wlog/auditlog/audit3log/internal` directory.

This will generate the "[categories.go](wlog/auditlog/audit3log/categories.go)" file. Note that the generator logic
assumes that the definitions in the Java code are written in a particular manner -- if the format or definition of the
definitions in the Java code changes, the generator must be updated to handle the new format.

Active TODOs
------------
* Improve testing loggers that produce non-JSON output (glog)
* Port over more tests for audit logs

License
-------
This project is made available under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0).
