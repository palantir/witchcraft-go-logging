// Copyright (c) 2026 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package benchmarks

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type callerResult struct {
	Operation string `json:"operation"`
	Count     int    `json:"count"`
	Success   bool   `json:"success"`
}

// Include context lookup and per-call parameter construction, which the fixture benchmarks exclude.
func BenchmarkSvc1LogCaller(b *testing.B) {
	for _, payload := range []string{"primitive", "structured"} {
		b.Run(payload, func(b *testing.B) {
			for _, construction := range []string{"reused", "constructed"} {
				b.Run(construction, func(b *testing.B) {
					RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
						// Runtime values include the boxing costs hidden by literal-only parameters.
						requestID := strings.Repeat("r", 32)
						attempt := len(requestID) + 300
						elapsedMs := float64(attempt) / 10
						user := "user-" + requestID
						logger := svc1log.NewFromCreator(io.Discard, wlog.InfoLevel, provider.NewLeveledLogger)
						ctx := svc1log.WithLoggerParams(svc1log.WithLogger(context.Background(), logger),
							svc1log.Origin("worker"),
							svc1log.SafeParam("component", "scheduler"),
							svc1log.Tag("deployment", "production"),
						)
						var reused []svc1log.Param
						if payload == "primitive" {
							reused = []svc1log.Param{
								svc1log.SafeParam("requestId", requestID),
								svc1log.SafeParam("attempt", attempt),
								svc1log.SafeParam("success", true),
								svc1log.SafeParam("elapsedMs", elapsedMs),
								svc1log.UnsafeParam("user", user),
								svc1log.Tag("operation", "dispatch"),
							}
						} else {
							reused = []svc1log.Param{
								svc1log.SafeParam("requestId", requestID),
								svc1log.SafeParam("result", callerResult{Operation: "dispatch", Count: attempt, Success: true}),
								svc1log.SafeParam("targets", []string{requestID, "host-b"}),
								svc1log.SafeParam("counts", map[string]int{"queued": attempt, "running": 2}),
								svc1log.UnsafeParam("user", user),
								svc1log.Tag("operation", "dispatch"),
							}
						}
						var emit func()
						switch {
						case construction == "reused":
							emit = func() { svc1log.FromContext(ctx).Info("Dispatched work", reused...) }
						case payload == "primitive":
							emit = func() {
								svc1log.FromContext(ctx).Info("Dispatched work",
									svc1log.SafeParam("requestId", requestID),
									svc1log.SafeParam("attempt", attempt),
									svc1log.SafeParam("success", true),
									svc1log.SafeParam("elapsedMs", elapsedMs),
									svc1log.UnsafeParam("user", user),
									svc1log.Tag("operation", "dispatch"),
								)
							}
						default:
							emit = func() {
								svc1log.FromContext(ctx).Info("Dispatched work",
									svc1log.SafeParam("requestId", requestID),
									svc1log.SafeParam("result", callerResult{Operation: "dispatch", Count: attempt, Success: true}),
									svc1log.SafeParam("targets", []string{requestID, "host-b"}),
									svc1log.SafeParam("counts", map[string]int{"queued": attempt, "running": 2}),
									svc1log.UnsafeParam("user", user),
									svc1log.Tag("operation", "dispatch"),
								)
							}
						}
						benchmarkLogCalls(b, emit)
					})
				})
			}
		})
	}
}
