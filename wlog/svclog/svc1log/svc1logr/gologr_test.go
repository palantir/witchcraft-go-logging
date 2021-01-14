// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package svc1logr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/assert"

	// Use zap as logger implementation
	_ "github.com/palantir/witchcraft-go-logging/wlog-zap"
)

func TestSvc1LogrWrapper(t *testing.T) {
	buf := new(bytes.Buffer)
	svcLogger := svc1log.New(buf, wlog.DebugLevel)

	logr1 := New(svcLogger, "foo")
	logr2 := logr1.WithValues("key2", "val2", "key3", "val3")
	logr3 := logr1.WithName("bar")

	logr1.Info("logr 1", "key1", "val1")
	assertLogLine(t, buf.Bytes(), objmatcher.MapMatcher{
		"level":   objmatcher.NewEqualsMatcher("INFO"),
		"time":    objmatcher.NewRegExpMatcher(".+"),
		"message": objmatcher.NewEqualsMatcher("logr 1"),
		"type":    objmatcher.NewEqualsMatcher(svc1log.TypeValue),
		"origin":  objmatcher.NewEqualsMatcher("foo"),
		"params": objmatcher.MapMatcher{
			"key1": objmatcher.NewEqualsMatcher("val1"),
		},
	})
	buf.Reset()

	logr2.Info("logr 2")
	assertLogLine(t, buf.Bytes(), objmatcher.MapMatcher{
		"level":   objmatcher.NewEqualsMatcher("INFO"),
		"time":    objmatcher.NewRegExpMatcher(".+"),
		"message": objmatcher.NewEqualsMatcher("logr 2"),
		"type":    objmatcher.NewEqualsMatcher(svc1log.TypeValue),
		"origin":  objmatcher.NewEqualsMatcher("foo"),
		"params": objmatcher.MapMatcher{
			"key2": objmatcher.NewEqualsMatcher("val2"),
			"key3": objmatcher.NewEqualsMatcher("val3"),
		},
	})
	buf.Reset()

	logr3.Error(fmt.Errorf("test error"), "logr 3")
	assertLogLine(t, buf.Bytes(), objmatcher.MapMatcher{
		"level":      objmatcher.NewEqualsMatcher("ERROR"),
		"time":       objmatcher.NewRegExpMatcher(".+"),
		"message":    objmatcher.NewEqualsMatcher("logr 3"),
		"type":       objmatcher.NewEqualsMatcher(svc1log.TypeValue),
		"origin":     objmatcher.NewEqualsMatcher("foo/bar"),
		"stacktrace": objmatcher.NewEqualsMatcher("test error"),
	})
	buf.Reset()
}

func assertLogLine(t *testing.T, logLine []byte, matcher objmatcher.MapMatcher) {
	logEntry := map[string]interface{}{}
	err := json.Unmarshal(logLine, &logEntry)
	assert.NoError(t, err)
	assert.NoError(t, matcher.Matches(logEntry))
}
