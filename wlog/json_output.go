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

package wlog

import (
	"bytes"
	"encoding"
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"io"
	"math"
	"slices"
	"sync"
)

type jsonOutput struct {
	mu sync.Mutex
	w  io.Writer
}

// Parameter values are arbitrary. Keep values that JSON v1 encodes (durations, structs without exported fields or
// with conflicting field names, invalid UTF-8) from failing the whole record. Sort map keys and skip duplicate-name
// checks as v1 did.
var jsonCompatOptions = json.JoinOptions(
	json.Deterministic(true),
	jsontext.AllowDuplicateNames(true),
	jsontext.AllowInvalidUTF8(true),
	jsonv1.FormatDurationAsNano(true),
	jsonv1.ReportErrorsWithLegacySemantics(true),
	json.WithMarshalers(jsonMarshalers),
)

// Encode errors as their message rather than their struct fields. Errors that define their own encoding keep it.
var errorMarshalers = json.MarshalToFunc(func(enc *jsontext.Encoder, err error) error {
	switch err.(type) {
	case json.Marshaler, json.MarshalerTo, encoding.TextMarshaler, encoding.TextAppender:
		return errors.ErrUnsupported
	}
	return enc.WriteToken(jsontext.String(err.Error()))
})

// JSON v2 allocates to make map entries and interface values addressable before encoding them. Write parameter maps
// and common value types directly and leave the rest to JSON v2.
var paramMapMarshalers = json.JoinMarshalers(
	json.MarshalToFunc(func(enc *jsontext.Encoder, m map[string]any) error {
		// Leave deep nesting to JSON v2, which detects cycles.
		if enc.StackDepth() > 64 {
			return errors.ErrUnsupported
		}
		return marshalAnyMap(enc, m)
	}),
	json.MarshalToFunc(func(enc *jsontext.Encoder, m map[string]string) error {
		return marshalSortedMap(enc, m, func(enc *jsontext.Encoder, v string) error {
			return enc.WriteToken(jsontext.String(v))
		})
	}),
)

var mapKeys = sync.Pool{New: func() any { return new([]string) }}

func marshalAnyMap(enc *jsontext.Encoder, m map[string]any) error {
	return marshalSortedMap(enc, m, marshalAnyValue)
}

func marshalSortedMap[V any](enc *jsontext.Encoder, m map[string]V, marshalValue func(*jsontext.Encoder, V) error) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	keys := mapKeys.Get().(*[]string)
	defer func() {
		clear(*keys)
		*keys = (*keys)[:0]
		mapKeys.Put(keys)
	}()
	for k := range m {
		*keys = append(*keys, k)
	}
	// Match JSON v2's deterministic order.
	slices.Sort(*keys)
	for _, k := range *keys {
		if err := enc.WriteToken(jsontext.String(k)); err != nil {
			return err
		}
		if err := marshalValue(enc, m[k]); err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

func marshalAnyValue(enc *jsontext.Encoder, v any) error {
	switch v := v.(type) {
	case nil:
		return enc.WriteToken(jsontext.Null)
	case string:
		return enc.WriteToken(jsontext.String(v))
	case bool:
		return enc.WriteToken(jsontext.Bool(v))
	case int:
		return enc.WriteToken(jsontext.Int(int64(v)))
	case int64:
		return enc.WriteToken(jsontext.Int(v))
	case int32:
		return enc.WriteToken(jsontext.Int(int64(v)))
	case uint64:
		return enc.WriteToken(jsontext.Uint(v))
	case float64:
		// jsontext encodes NaN and infinities as strings; JSON v2 reports them as errors.
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			return enc.WriteToken(jsontext.Float(v))
		}
	case map[string]any:
		if enc.StackDepth() <= 64 {
			return marshalAnyMap(enc, v)
		}
	}
	return json.MarshalEncode(enc, v)
}

var jsonBuffers = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func putJSONBuffer(buf *bytes.Buffer) {
	// Do not let an unusually large record inflate the shared pool.
	if buf.Cap() > 64<<10 {
		return
	}
	buf.Reset()
	jsonBuffers.Put(buf)
}

func newJSONOutput(w io.Writer) *jsonOutput {
	return &jsonOutput{w: w}
}

func (o *jsonOutput) write(value any) error {
	buf := jsonBuffers.Get().(*bytes.Buffer)
	defer putJSONBuffer(buf)

	// Marshal a complete record before writing so a bad value cannot corrupt the stream.
	if err := json.MarshalWrite(buf, value, jsonCompatOptions); err != nil {
		return err
	}
	buf.WriteByte('\n')

	o.mu.Lock()
	defer o.mu.Unlock()
	n, err := o.w.Write(buf.Bytes())
	if err == nil && n != buf.Len() {
		return io.ErrShortWrite
	}
	return err
}
