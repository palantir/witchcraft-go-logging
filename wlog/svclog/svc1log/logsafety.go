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

package svc1log

import (
	"fmt"
	"reflect"
)

// validateLogSafety validates that the given value is safe to be logged in a SafeParam.
// It recursively checks for any fields marked with safety:"unsafe" or safety:"do-not-log" tags.
// If unsafe fields are found, it returns a safe replacement value and an error describing the issue.
func validateLogSafety(v interface{}) (interface{}, error) {
	visited := make(map[uintptr]bool)
	if err := validateLogSafetyRecursive(reflect.ValueOf(v), visited); err != nil {
		return fmt.Sprintf("[REDACTED: Value contains fields marked as safety:\"unsafe\" or safety:\"do-not-log\" and cannot be logged safely. Please use UnsafeParam() instead or remove unsafe safety tags. Error: %v]", err), err
	}
	return v, nil
}

// validateLogSafetyRecursive recursively validates the log safety of a reflected value.
func validateLogSafetyRecursive(rv reflect.Value, visited map[uintptr]bool) error {
	// Handle nil values
	if !rv.IsValid() {
		return nil
	}

	// Handle pointers and interfaces
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil
		}
		// Prevent infinite recursion on circular references
		if rv.Kind() == reflect.Ptr {
			ptr := rv.Pointer()
			if visited[ptr] {
				return nil
			}
			visited[ptr] = true
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		return validateStructSafety(rv, visited)
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			if err := validateLogSafetyRecursive(rv.Index(i), visited); err != nil {
				return fmt.Errorf("element %d: %w", i, err)
			}
		}
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			if err := validateLogSafetyRecursive(rv.MapIndex(key), visited); err != nil {
				return fmt.Errorf("map value for key %v: %w", key.Interface(), err)
			}
		}
	}
	return nil
}

// validateStructSafety validates that a struct doesn't contain any fields marked as safety:"unsafe" or safety:"do-not-log".
func validateStructSafety(rv reflect.Value, visited map[uintptr]bool) error {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		// Skip unexported fields as they won't be serialized by JSON marshaling
		if !field.IsExported() {
			continue
		}

		// Check the safety tag
		safetyTag := field.Tag.Get("safety")
		switch safetyTag {
		case "unsafe", "do-not-log":
			return fmt.Errorf("field %s.%s is marked as safety:\"%s\" and cannot be used in SafeParam", rt.Name(), field.Name, safetyTag)
		case "safe":
			// Explicitly safe, continue with recursive validation
		case "":
			// No safety tag - default behavior, continue with recursive validation
		default:
			// Unknown safety tag value - treat conservatively as unsafe
			return fmt.Errorf("field %s.%s has unknown safety tag value \"%s\"", rt.Name(), field.Name, safetyTag)
		}

		// Recursively validate nested structures
		fieldValue := rv.Field(i)
		if err := validateLogSafetyRecursive(fieldValue, visited); err != nil {
			return fmt.Errorf("field %s.%s: %w", rt.Name(), field.Name, err)
		}
	}
	return nil
}
