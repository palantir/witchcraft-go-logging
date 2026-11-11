// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

// Package objecttest checks that the entry path of a typed logger writes every field of its generated log object.
package objecttest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
	"github.com/palantir/witchcraft-go-logging/wlog/wlogtests"
	"github.com/stretchr/testify/require"
)

// Fill sets every exported field of the struct that ptr points to, recursively, to a non-zero value derived from its
// path. A type without exported fields, such as an enum or union, must be supplied in overrides.
func Fill(t testing.TB, ptr any, overrides ...any) {
	t.Helper()
	values := map[reflect.Type]reflect.Value{
		reflect.TypeFor[datetime.DateTime](): reflect.ValueOf(datetime.DateTime(time.Date(2025, 1, 2, 3, 4, 5, 6000, time.UTC))),
		reflect.TypeFor[logging.LogLevel]():  reflect.ValueOf(logging.New_LogLevel(logging.LogLevel_WARN)),
		reflect.TypeFor[logtypes.LogLevel](): reflect.ValueOf(logtypes.LogLevelWarn),
	}
	for _, o := range overrides {
		values[reflect.TypeOf(o)] = reflect.ValueOf(o)
	}
	f := filler{values: values}
	f.fill(reflect.ValueOf(ptr).Elem(), reflect.TypeOf(ptr).Elem().Name())
	require.NoError(t, f.err)
}

type filler struct {
	values map[reflect.Type]reflect.Value
	next   int
	err    error
}

func (f *filler) fill(v reflect.Value, path string) {
	if o, ok := f.values[v.Type()]; ok {
		v.Set(o)
		return
	}
	f.next++
	switch v.Kind() {
	case reflect.String:
		v.SetString(path)
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(f.next))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(uint64(f.next))
	case reflect.Float32, reflect.Float64:
		v.SetFloat(float64(f.next) + 0.5)
	case reflect.Interface:
		v.Set(reflect.ValueOf(path))
	case reflect.Pointer:
		v.Set(reflect.New(v.Type().Elem()))
		f.fill(v.Elem(), path)
	case reflect.Slice:
		v.Set(reflect.MakeSlice(v.Type(), 1, 1))
		f.fill(v.Index(0), path+"[0]")
	case reflect.Array:
		for i := range v.Len() {
			f.fill(v.Index(i), fmt.Sprintf("%s[%d]", path, i))
		}
	case reflect.Map:
		key := reflect.New(v.Type().Key()).Elem()
		f.fill(key, path+".key")
		elem := reflect.New(v.Type().Elem()).Elem()
		f.fill(elem, path+".value")
		v.Set(reflect.MakeMap(v.Type()))
		v.SetMapIndex(key, elem)
	case reflect.Struct:
		exported := 0
		for i := range v.NumField() {
			if field := v.Type().Field(i); field.IsExported() {
				exported++
				f.fill(v.Field(i), path+"."+field.Name)
			}
		}
		if exported == 0 && v.NumField() > 0 && f.err == nil {
			f.err = fmt.Errorf("%s: type %s has no exported fields; pass a value of it to Fill", path, v.Type())
		}
	default:
		if f.err == nil {
			f.err = fmt.Errorf("%s: unsupported kind %s", path, v.Kind())
		}
	}
}

// RequireParity requires that the filled and empty log objects, emitted to the built-in JSON logger's typed method and
// to its entry path, decode to the same JSON value. The empty case checks that both paths omit the same unset fields.
func RequireParity[T any](t *testing.T, filled, empty *T, direct func(wlog.LeveledLogger, *T), entry func(wlog.LeveledLogger, *T)) {
	t.Helper()
	for name, log := range map[string]*T{"filled": filled, "empty": empty} {
		t.Run(name, func(t *testing.T) {
			require.Equal(t,
				record(t, wlog.NewJSONMarshalLoggerProvider(), func(l wlog.LeveledLogger) { direct(l, log) }),
				record(t, wlogtests.MapJSONProvider(), func(l wlog.LeveledLogger) { entry(l, log) }),
			)
		})
	}
}

func record(t testing.TB, provider wlog.LoggerProvider, emit func(wlog.LeveledLogger)) map[string]any {
	t.Helper()
	var buf bytes.Buffer
	emit(provider.NewLeveledLogger(io.Writer(&buf), wlog.DebugLevel))
	var out map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &out), buf.String())
	return out
}
