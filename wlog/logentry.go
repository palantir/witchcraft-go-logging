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
	"io"
	"maps"
	"reflect"
)

// LogEntry is an interface that represents a log entry, which is an entry on which a key-value pair can be set.
type LogEntry interface {
	// StringValue sets the value for the specified key to be the specified value.
	// Overwrites any previous value associated with the key.
	StringValue(k, v string)

	// OptionalStringValue sets the value for the specified key to be the specified value *if* the provided value
	// is non-empty. If the provided value is empty, this operation removes the key.
	// Overwrites any previous value associated with the key if the provided value is non-empty. If the provided value
	// is empty, this operation removes the key.
	//
	// Note that these semantics mean that it is not possible to use this function to associate the empty string value
	// ("") with a key: if this behavior is desired, call the StringValue function with an empty string instead.
	OptionalStringValue(k, v string)

	// SafeLongValue sets the value for the specified key to be the specified value.
	// Overwrites any previous value associated with the key.
	SafeLongValue(k string, v int64)

	// IntValue sets the value for the specified key to be the specified value.
	// Overwrites any previous value associated with the key.
	IntValue(k string, v int32)

	// StringListValue sets the value for the specified key to be the specified value. If the provided value is empty or
	// nil, sets the value for the specified key to be an empty array.
	// Overwrites any previous value associated with the key.
	StringListValue(k string, v []string)

	// StringMapValue sets the value for the specified key to be a map that contains all the entries in the provided
	// map.
	//
	// - If there is no value for the specified key, the value is set to be a map that contains all the entries in the
	//   provided map (if the provided map is nil or empty, the result is an empty map).
	// - If the existing value for the specified key is a map constructed from previous calls to StringMapValue, the
	//   entries in the provided map are merged with the existing entries. In this case, if there are any entries with
	//   matching keys, the values in the provided map will overwrite the existing values.
	// - If the existing value for the specified key is any other value, this call overwrites the value with a map that
	//   contains the entries in the provided map (if the provided map is nil or empty, this is an empty map).
	//
	// Note that these semantics mean that this function cannot remove any keys in the map specified by the provided
	// key. If a caller wants to ensure that the map for the specified key contains only the elements provided by this
	// call, a roundabout way of doing so is to call another function like "StringValue" for the same key first, since
	// that call will overwrite any existing StringMapValue and then the subsequent call to StringMapValue will
	// overwrite that value.
	StringMapValue(k string, v map[string]string)

	// AnyMapValue sets the value for the specified key to be a map that contains all the entries in the provided map.
	// If there is no value for the specified key, the value is set to be a map that contains all the entries in the
	// provided map (if the provided map is nil or empty, the result is an empty map).
	//
	// - If there is no value for the specified key, the value is set to be a map that contains all the entries in the
	//   provided map (if the provided map is nil or empty, the result is an empty map).
	// - If the existing value for the specified key is a map constructed from previous calls to AnyMapValue, the
	//   entries in the provided map are merged with the existing entries. In this case, if there are any entries with
	//   matching keys, the values in the provided map will overwrite the existing values.
	// - If the existing value for the specified key is any other value, this call overwrites the value with a map that
	//   contains the entries in the provided map (if the provided map is nil or empty, this is an empty map).
	//
	// Note that these semantics mean that this function cannot remove any keys in the map specified by the provided
	// key. If a caller wants to ensure that the map for the specified key contains only the elements provided by this
	// call, a roundabout way of doing so is to call another function like "StringValue" for the same key first, since
	// that call will overwrite any existing AnyMapValue and then the subsequent call to AnyMapValue will overwrite that
	// value.
	AnyMapValue(k string, v map[string]any)

	// ObjectValue sets the value for the specified key to be the specified value. If marshalerType is non-nil, then if
	// a custom marshaler is registered for that type, it may be used to log the entry. If marshalerType is nil or no
	// marshaler is registered for the provided type, the entry is logged using reflection.
	// Overwrites any previous value associated with the key.
	ObjectValue(k string, v any, marshalerType reflect.Type)
}

type Logger interface {
	Log(params ...Param)
}

type LoggerCreator func(w io.Writer) Logger

type LeveledLoggerCreator func(w io.Writer, level LogLevel) LeveledLogger

type LeveledLogger interface {
	Debug(msg string, params ...Param)
	Info(msg string, params ...Param)
	Warn(msg string, params ...Param)
	Error(msg string, params ...Param)
	SetLevel(level LogLevel)
}

type LevelChecker interface {
	// Enabled determines whether the provided level should be logged.
	// If implemented with LeveledLogger or SetLevel, they must remain consistent with Enabled.
	Enabled(level LogLevel) bool
}

type MapValueEntries struct {
	stringMapValues map[string]map[string]string
	anyMapValues    map[string]map[string]interface{}
}

func (m *MapValueEntries) StringMapValue(k string, v map[string]string) {
	mapValueEntriesAddValuesToMap(m, &m.stringMapValues, k, v)
}

func (m *MapValueEntries) AnyMapValue(k string, v map[string]interface{}) {
	mapValueEntriesAddValuesToMap(m, &m.anyMapValues, k, v)
}

// Delete deletes the specified key from all the maps in the data structure.
func (m *MapValueEntries) Delete(k string) {
	delete(m.stringMapValues, k)
	delete(m.anyMapValues, k)
}

// mapValueEntriesAddValuesToMap adds the entries in the provided "values" map to the map associated with the provided
// key in the map of maps referenced by mapValuesPtr.
//
// If the value of the provided mapValuesPtr is nil, this function allocates a new map for it.
//
// If the map of maps referenced by mapValuesPtr has a map associated with the provided key, the entries in the provided
// map are merged with the existing values. If there are any entries with matching keys, the values in the provided map
// will overwrite the existing values.
//
// If the map of maps referenced by mapValuesPtr does not have a map associated with the provided key, sets the value of
// the key to be a new map and then adds all the entries from the provided map. This means that, if the provided map is
// nil or empty, the result will be that the key will be associated with an empty map.
func mapValueEntriesAddValuesToMap[ValT any](m *MapValueEntries, mapValuesPtr *map[string]map[string]ValT, k string, v map[string]ValT) {
	if *mapValuesPtr == nil {
		*mapValuesPtr = make(map[string]map[string]ValT)
	}

	entryMapVals, ok := (*mapValuesPtr)[k]
	if !ok {
		// if entry does not exist, initialize with an empty map
		entryMapVals = make(map[string]ValT)
		(*mapValuesPtr)[k] = entryMapVals
	}
	// add all provided elements to map
	maps.Copy(entryMapVals, v)

	// clear key from all maps
	m.Delete(k)

	// set key on target map
	(*mapValuesPtr)[k] = entryMapVals
}

func (m *MapValueEntries) StringMapValues() map[string]map[string]string {
	return m.stringMapValues
}

func (m *MapValueEntries) AnyMapValues() map[string]map[string]interface{} {
	return m.anyMapValues
}
