// Copyright (c) 2018 Palantir Technologies. All rights reserved.
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

package wlog

import (
	"bytes"
	"encoding/json/v2"
	"io"
	"log"
	"sync"
)

// NewJSONMarshalLoggerProvider returns a new logger provider that uses a MapLogEntry as its log entry and encodes it
// with JSON v2.
func NewJSONMarshalLoggerProvider() LoggerProvider {
	return &jsonMarshalLoggerProvider{}
}

type jsonMarshalLogger struct {
	mu sync.Mutex
	w  io.Writer
	*AtomicLogLevel
}

func (l *jsonMarshalLogger) Log(params ...Param) {
	l.logOutput(params)
}

func (l *jsonMarshalLogger) Debug(msg string, params ...Param) {
	if l.Enabled(DebugLevel) {
		l.logOutput(ParamsWithMessage(msg, params))
	}
}

func (l *jsonMarshalLogger) Info(msg string, params ...Param) {
	if l.Enabled(InfoLevel) {
		l.logOutput(ParamsWithMessage(msg, params))
	}
}

func (l *jsonMarshalLogger) Warn(msg string, params ...Param) {
	if l.Enabled(WarnLevel) {
		l.logOutput(ParamsWithMessage(msg, params))
	}
}

func (l *jsonMarshalLogger) Error(msg string, params ...Param) {
	if l.Enabled(ErrorLevel) {
		l.logOutput(ParamsWithMessage(msg, params))
	}
}

func (l *jsonMarshalLogger) logOutput(params []Param) {
	entry := NewMapLogEntry()
	ApplyParams(entry, params)
	if err := l.write(entry.AllValues()); err != nil {
		log.Printf("json logger failed: %v", err)
	}
}

var jsonBuffers = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func (l *jsonMarshalLogger) write(value any) error {
	buf := jsonBuffers.Get().(*bytes.Buffer)
	defer func(buf *bytes.Buffer) {
		// Do not let an unusually large record inflate the shared pool.
		if buf.Cap() > 64<<10 {
			return
		}
		buf.Reset()
		jsonBuffers.Put(buf)
	}(buf)

	// Marshal a complete record before writing so a bad value cannot corrupt the stream.
	if err := json.MarshalWrite(buf, value, jsonCompatOptions); err != nil {
		return err
	}
	buf.WriteByte('\n')

	l.mu.Lock()
	defer l.mu.Unlock()
	n, err := l.w.Write(buf.Bytes())
	if err == nil && n != buf.Len() {
		return io.ErrShortWrite
	}
	return err
}

type jsonMarshalLoggerProvider struct{}

func (*jsonMarshalLoggerProvider) NewLogger(w io.Writer) Logger {
	return &jsonMarshalLogger{
		w: w,
	}
}

func (*jsonMarshalLoggerProvider) NewLeveledLogger(w io.Writer, level LogLevel) LeveledLogger {
	return &jsonMarshalLogger{
		w:              w,
		AtomicLogLevel: NewAtomicLogLevel(level),
	}
}

func (p *jsonMarshalLoggerProvider) NewLogEntry() LogEntry {
	return NewMapLogEntry()
}
