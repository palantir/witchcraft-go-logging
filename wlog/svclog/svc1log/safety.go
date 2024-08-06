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

package svc1log

import (
	"fmt"
	"reflect"
)

const (
	safetyTag = "safety"

	safeValue   = "safe"
	unsafeValue = "unsafe"
)

type LogSafety struct {
	Safe    bool
	Message string
}

// ok, next thing is whether or not the value is cacheable
// if at any point a type is interface... the answer is no

func IsParamSafe(paramsMap map[string]interface{}) map[string]LogSafety {
	// This is probably not quite right. We should evaluate other types that could include structs
	safetyMap := make(map[string]LogSafety)
	for key, val := range paramsMap {
		safe, message := isSafeRecursive(val)
		safetyMap[key] = LogSafety{
			Safe:    safe,
			Message: message,
		}
	}
	return safetyMap
}

func IsSafe(val interface{}) (bool, string) {
	return isSafeRecursive(val)
}

func isSafeRecursive(val interface{}) (bool, string) {
	if val == nil {
		// Nil vals are safe
		return true, ""
	}

	valT := reflect.TypeOf(val)
	valV := reflect.ValueOf(val)

	if isPrimitiveType(valT.Kind()) {
		return true, ""
	}

	// This is a choice... personally I think that interface types can never be safe, due to
	// the unknown nature of what can be put in it. Always return false if a type contains one.
	//if valT.Kind() == reflect.Interface {
	//	return false, interfaceInSafeArgMessage(valT.Name())
	//}

	// one inner type - array, slice, chan, or pointer
	if valT.Kind() == reflect.Array || valT.Kind() == reflect.Slice || valT.Kind() == reflect.Chan || valT.Kind() == reflect.Pointer {
		if isPrimitiveType(valT.Elem().Kind()) {
			return true, ""
		}
		newVal := reflect.New(valT.Elem())
		return isSafeRecursive(newVal.Elem().Interface())
	}

	// two inner types - map
	if valT.Kind() == reflect.Map {
		// need to check key and values
		mapSafe := true
		message := ""
		if !isPrimitiveType(valT.Key().Kind()) {
			newVal := reflect.New(valT.Key())
			mapSafe, message = isSafeRecursive(newVal.Elem().Interface())
		}
		if mapSafe && !isPrimitiveType(valT.Elem().Kind()) {
			newVal := reflect.New(valT.Elem())
			mapSafe, message = isSafeRecursive(newVal.Elem().Interface())
		}
		// still idk
		return mapSafe, message
	}

	// struct
	if valT.Kind() == reflect.Struct {
		safe := true
		message := ""
		for i := 0; i < valT.NumField(); i++ {
			structFieldSafe, msg := structFieldIsSafe(valT.Field(i), valV.Field(i))

			safe = safe && structFieldSafe
			if !structFieldSafe {
				message = msg
			}
		}
		return safe, message
	}

	// This is a base case that should never get hit.
	return true, ""
}

func structFieldIsSafe(field reflect.StructField, fieldVal reflect.Value) (bool, string) {
	tagVal, ok := field.Tag.Lookup(safetyTag)
	if !ok {
		// If no tag is set, set it to safe.
		tagVal = safeValue
	}

	if tagVal == unsafeValue {
		return false, unsafeArgMessage(field)
	}

	fieldValIsSafe := true
	message := ""
	if !isPrimitiveType(fieldVal.Kind()) {
		// If cannot interface (non-exported field), don't dive further.
		// Default marshalling won't include this field either.
		if fieldVal.CanInterface() {
			fieldValIsSafe, message = isSafeRecursive(fieldVal.Interface())
		}
	}

	return tagVal == safeValue && fieldValIsSafe, message
}

func isPrimitiveType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Array:
		return false
	case reflect.Interface:
		return false
	case reflect.Slice:
		return false
	case reflect.Map:
		return false
	case reflect.Struct:
		return false
	case reflect.Chan:
		return false
	default:
		return true
	}
}

func unsafeArgMessage(field reflect.StructField) string {
	return fmt.Sprintf("'%s' was passed as a safe arg, but is actually tagged as unsafe.", field.Name)
}
