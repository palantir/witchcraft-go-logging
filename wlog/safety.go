// Copyright (c) 2024 Palantir Technologies. All rights reserved.
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
	"fmt"
	"maps"
	"reflect"
	"sync"
)

const (
	safetyTag = "safety"

	safeValue   = "safe"
	unsafeValue = "unsafe"
)

type LogSafety struct {
	Safe      bool
	Message   string
	cacheable bool
}

type SafetyChecker interface {
	ParamsSafe(safeParams map[string]interface{}) map[string]LogSafety
	OmitUnsafeParams(safeParams map[string]interface{}) map[string]interface{}
}

type defaultSafetyChecker struct {
	cache map[string]struct{}
	lock  sync.RWMutex
}

func NewSafetyChecker() SafetyChecker {
	return &defaultSafetyChecker{
		cache: make(map[string]struct{}),
	}
}

func (d *defaultSafetyChecker) OmitUnsafeParams(safeParams map[string]interface{}) map[string]interface{} {
	safetyMap := d.ParamsSafe(safeParams)
	newParams := make(map[string]interface{})
	for key, val := range safeParams {
		safety, _ := safetyMap[key]
		if !safety.Safe {
			newParams[key] = safety.Message
		} else {
			newParams[key] = val
		}
	}
	return newParams
}

func (d *defaultSafetyChecker) ParamsSafe(safeParams map[string]interface{}) map[string]LogSafety {
	safetyMap := make(map[string]LogSafety)
	for key, val := range safeParams {
		cache := d.getCachedSafeStructs()
		safe, message, safeStructs := isSafeRecursive(val, cache)
		safetyMap[key] = LogSafety{
			Safe:    safe,
			Message: message,
		}
		// TODO(awerner): Consider moving this out of the for loop.
		// Advantages:
		//  Repopulating and copying the map each loop means that structs that appear in multiple params will be O(1)
		//  in computing if it is safe when looking in a subsequent struct.
		// Disadvantages:
		//  We recopy + lock the map to write each iteration of the loop.
		// Should benchmark to decide which is actually faster on realistically sized param maps.
		d.putSafeStructsInCache(safeStructs)
	}

	return safetyMap
}

func (d *defaultSafetyChecker) getCachedSafeStructs() map[string]struct{} {
	d.lock.RLock()
	defer d.lock.RUnlock()

	// TODO(awerner): Also reconsider this... might be quite memory intensive. Not sure if the tradeoffs are really
	//   worth it. Other option is to pass the real cache and lock it every time. Likely solved by benchmarking.
	// Return a copy of the map so that other threads can interact with the map at the same time, no need to
	//   grab a lock on every map read.
	cacheCopy := make(map[string]struct{})
	maps.Copy(d.cache, cacheCopy)
	return cacheCopy
}

func (d *defaultSafetyChecker) putSafeStructsInCache(structNames []string) {
	if len(structNames) == 0 {
		return
	}

	d.lock.Lock()
	defer d.lock.Unlock()
	for _, structName := range structNames {
		d.cache[structName] = struct{}{}
	}
}

func IsParamSafe(paramsMap map[string]interface{}) map[string]LogSafety {
	safetyMap := make(map[string]LogSafety)
	for key, val := range paramsMap {
		safe, message, _ := isSafeRecursive(val, map[string]struct{}{})
		safetyMap[key] = LogSafety{
			Safe:    safe,
			Message: message,
		}
	}
	return safetyMap
}

func isSafeRecursive(val interface{}, cachedSafeStructs map[string]struct{}) (bool, string, []string) {
	if val == nil {
		// Nil vals are safe
		return true, "", []string{}
	}

	valT := reflect.TypeOf(val)
	valV := reflect.ValueOf(val)

	if isPrimitiveType(valT.Kind()) {
		return true, "", []string{}
	}
	// For now
	if valT.Kind() == reflect.Interface {
		return false, "", []string{}
	}

	// one inner type - array, slice, chan, or pointer
	if valT.Kind() == reflect.Array || valT.Kind() == reflect.Slice || valT.Kind() == reflect.Chan || valT.Kind() == reflect.Pointer {
		if isPrimitiveType(valT.Elem().Kind()) {
			return true, "", []string{}
		}
		newVal := reflect.New(valT.Elem())
		return isSafeRecursive(newVal.Elem().Interface(), cachedSafeStructs)
	}

	// two inner types - map
	if valT.Kind() == reflect.Map {
		// need to check key and values
		mapSafe := true
		message := ""
		safeStructs := make([]string, 0)
		if !isPrimitiveType(valT.Key().Kind()) {
			newVal := reflect.New(valT.Key())
			mapSafe, message, safeStructs = isSafeRecursive(newVal.Elem().Interface(), cachedSafeStructs)
		}
		if mapSafe && !isPrimitiveType(valT.Elem().Kind()) {
			newVal := reflect.New(valT.Elem())
			mapSafe, message, safeStructs = isSafeRecursive(newVal.Elem().Interface(), cachedSafeStructs)
		}
		return mapSafe, message, safeStructs
	}

	// struct
	if valT.Kind() == reflect.Struct {
		if _, present := cachedSafeStructs[valT.Name()]; present {
			return true, "", []string{}
		}

		safe := true
		message := ""
		safeStructs := make([]string, 0)
		for i := 0; i < valT.NumField(); i++ {
			structFieldSafe, msg, structs := structFieldIsSafe(valT.Field(i), valV.Field(i), cachedSafeStructs)

			safe = safe && structFieldSafe
			if !structFieldSafe {
				message = msg
			}

			safeStructs = append(safeStructs, structs...)
		}

		// TODO(awerner): add this to the cached safe structs for usage if encounter duplicate structs in same parent obj
		if safe {
			safeStructs = append(safeStructs, valT.Name())
		}

		return safe, message, safeStructs
	}

	// This is a base case that should never get hit. Should remove this once it is no longer possible to hit...
	// Currently Kind() == Interface hits it, which I'm still not sure what actually has that type.
	return true, "", []string{}
}

func structFieldIsSafe(field reflect.StructField, fieldVal reflect.Value, cachedSafeStructs map[string]struct{}) (bool, string, []string) {
	tagVal, ok := field.Tag.Lookup(safetyTag)
	if !ok {
		// If no tag is set, set it to safe.
		tagVal = safeValue
	}

	if tagVal == unsafeValue {
		return false, unsafeArgMessage(field), []string{}
	}

	fieldValIsSafe := true
	message := ""
	safeStructs := make([]string, 0)
	if !isPrimitiveType(fieldVal.Kind()) {
		// If cannot interface (non-exported field), don't dive further.
		// Default marshalling won't include this field either.
		if fieldVal.CanInterface() {
			fieldValIsSafe, message, safeStructs = isSafeRecursive(fieldVal.Interface(), cachedSafeStructs)
		}
	}

	return fieldValIsSafe, message, safeStructs
}

func isPrimitiveType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Interface:
		return false
	case reflect.Array:
		return false
	case reflect.Slice:
		return false
	case reflect.Chan:
		return false
	case reflect.Pointer:
		return false
	case reflect.Map:
		return false
	case reflect.Struct:
		return false
	default:
		return true
	}
}

func unsafeArgMessage(field reflect.StructField) string {
	return fmt.Sprintf("'%s' was passed as a safe arg, but is actually tagged as unsafe.", field.Name)
}
