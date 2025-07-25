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
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/stretchr/testify/assert"
)

// Test structs for log safety validation
type SafeStruct struct {
	Name  string `json:"name" safety:"safe"`
	ID    int    `json:"id" safety:"safe"`
	Value string `json:"value"` // No tag defaults to safe validation
}

type UnsafeStruct struct {
	Name     string `json:"name" safety:"safe"`
	Password string `json:"password" safety:"unsafe"`
}

type DoNotLogStruct struct {
	Name   string `json:"name" safety:"safe"`
	Secret string `json:"secret" safety:"do-not-log"`
}

type NestedUnsafeStruct struct {
	SafeField   string       `json:"safe_field" safety:"safe"`
	UnsafeChild UnsafeStruct `json:"unsafe_child"`
}

type NestedDoNotLogStruct struct {
	SafeField     string         `json:"safe_field" safety:"safe"`
	DoNotLogChild DoNotLogStruct `json:"do_not_log_child"`
}

type MixedSafetyStruct struct {
	SafeField     string `json:"safe_field" safety:"safe"`
	UnsafeField   string `json:"unsafe_field" safety:"unsafe"`
	DoNotLogField string `json:"do_not_log_field" safety:"do-not-log"`
	NoTagField    string `json:"no_tag_field"` // Should be allowed
}

type UnknownSafetyTagStruct struct {
	Name         string `json:"name" safety:"safe"`
	UnknownField string `json:"unknown_field" safety:"unknown-value"`
}

type CircularStruct struct {
	Name string          `json:"name" safety:"safe"`
	Self *CircularStruct `json:"self"`
}

func TestValidateLogSafety_SafeValues(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{"nil", nil},
		{"string", "safe string"},
		{"int", 42},
		{"bool", true},
		{"safe struct", SafeStruct{Name: "test", ID: 1, Value: "value"}},
		{"slice of safe structs", []SafeStruct{{Name: "test1"}, {Name: "test2"}}},
		{"map of safe values", map[string]string{"key1": "value1", "key2": "value2"}},
		{"pointer to safe struct", &SafeStruct{Name: "test", ID: 1}},
		{"circular safe struct", &CircularStruct{Name: "test", Self: &CircularStruct{Name: "nested"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			safeValue, err := validateLogSafety(tt.value)
			assert.NoError(t, err, "Expected no error for safe value")
			assert.Equal(t, tt.value, safeValue, "Safe value should be returned unchanged")
		})
	}
}

func TestValidateLogSafety_UnsafeValues(t *testing.T) {
	tests := []struct {
		name        string
		value       interface{}
		expectedErr string
	}{
		{
			"unsafe struct",
			UnsafeStruct{Name: "test", Password: "secret"},
			"field UnsafeStruct.Password is marked as safety:\"unsafe\"",
		},
		{
			"do-not-log struct",
			DoNotLogStruct{Name: "test", Secret: "secret"},
			"field DoNotLogStruct.Secret is marked as safety:\"do-not-log\"",
		},
		{
			"mixed safety struct",
			MixedSafetyStruct{SafeField: "safe", UnsafeField: "unsafe", DoNotLogField: "secret", NoTagField: "ok"},
			"field MixedSafetyStruct.UnsafeField is marked as safety:\"unsafe\"",
		},
		{
			"unknown safety tag",
			UnknownSafetyTagStruct{Name: "test", UnknownField: "value"},
			"field UnknownSafetyTagStruct.UnknownField has unknown safety tag value \"unknown-value\"",
		},
		{
			"nested unsafe struct",
			NestedUnsafeStruct{SafeField: "safe", UnsafeChild: UnsafeStruct{Name: "test", Password: "secret"}},
			"field NestedUnsafeStruct.UnsafeChild: field UnsafeStruct.Password is marked as safety:\"unsafe\"",
		},
		{
			"nested do-not-log struct",
			NestedDoNotLogStruct{SafeField: "safe", DoNotLogChild: DoNotLogStruct{Name: "test", Secret: "secret"}},
			"field NestedDoNotLogStruct.DoNotLogChild: field DoNotLogStruct.Secret is marked as safety:\"do-not-log\"",
		},
		{
			"slice containing unsafe struct",
			[]UnsafeStruct{{Name: "test", Password: "secret"}},
			"element 0: field UnsafeStruct.Password is marked as safety:\"unsafe\"",
		},
		{
			"map containing do-not-log struct",
			map[string]DoNotLogStruct{"key": {Name: "test", Secret: "secret"}},
			"map value for key key: field DoNotLogStruct.Secret is marked as safety:\"do-not-log\"",
		},
		{
			"pointer to unsafe struct",
			&UnsafeStruct{Name: "test", Password: "secret"},
			"field UnsafeStruct.Password is marked as safety:\"unsafe\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			safeValue, err := validateLogSafety(tt.value)
			assert.Error(t, err, "Expected error for unsafe value")
			assert.Contains(t, err.Error(), tt.expectedErr, "Error message should contain expected text")

			// Verify that a redacted message is returned
			assert.IsType(t, "", safeValue, "Safe value should be a string")
			redactedMsg, ok := safeValue.(string)
			assert.True(t, ok, "Safe value should be a string")
			assert.Contains(t, redactedMsg, "[REDACTED:", "Should contain redaction message")
			assert.Contains(t, redactedMsg, "safety:", "Should explain why it was redacted")
		})
	}
}

func TestSafeParam_ValidatesInput(t *testing.T) {
	// Test that SafeParam accepts safe values
	safeStruct := SafeStruct{Name: "test", ID: 1}
	param := SafeParam("test", safeStruct)
	assert.NotNil(t, param, "SafeParam should return a parameter for safe input")

	// Test that SafeParam replaces unsafe values with redaction message
	unsafeStruct := UnsafeStruct{Name: "test", Password: "secret"}
	param = SafeParam("test", unsafeStruct)
	assert.NotNil(t, param, "SafeParam should return a parameter even for unsafe input")

	// Verify the redacted content by applying to a MapLogEntry
	logEntry := wlog.NewMapLogEntry()
	param.apply(logEntry)

	allValues := logEntry.AllValues()
	paramsValue, exists := allValues[ParamsKey]
	assert.True(t, exists, "Should have params field")

	paramsMap, ok := paramsValue.(map[string]interface{})
	assert.True(t, ok, "Params should be a map")

	testValue, exists := paramsMap["test"]
	assert.True(t, exists, "Should have test key")

	redactedStr, ok := testValue.(string)
	assert.True(t, ok, "Redacted value should be a string")
	assert.Contains(t, redactedStr, "[REDACTED:", "Should contain redaction message")
	assert.Contains(t, redactedStr, "safety:\"unsafe\"", "Should explain why it was redacted")
}

func TestSafeParams_ValidatesInput(t *testing.T) {
	// Test that SafeParams accepts safe values
	safeMap := map[string]interface{}{
		"safe1": SafeStruct{Name: "test1"},
		"safe2": "safe string",
	}
	param := SafeParams(safeMap)
	assert.NotNil(t, param, "SafeParams should return a parameter for safe input")

	// Test that SafeParams replaces unsafe values with redaction message
	unsafeMap := map[string]interface{}{
		"safe":   SafeStruct{Name: "test1"},
		"unsafe": UnsafeStruct{Name: "test", Password: "secret"},
	}
	param = SafeParams(unsafeMap)
	assert.NotNil(t, param, "SafeParams should return a parameter even for unsafe input")

	// Verify the redacted content by applying to a MapLogEntry
	logEntry := wlog.NewMapLogEntry()
	param.apply(logEntry)

	allValues := logEntry.AllValues()
	paramsValue, exists := allValues[ParamsKey]
	assert.True(t, exists, "Should have params field")

	paramsMap, ok := paramsValue.(map[string]interface{})
	assert.True(t, ok, "Params should be a map")

	// Check that the safe value is preserved
	safeValue, exists := paramsMap["safe"]
	assert.True(t, exists, "Should have safe key")
	assert.IsType(t, SafeStruct{}, safeValue, "Safe value should be original struct")

	// Check that the unsafe value is redacted
	unsafeValue, exists := paramsMap["unsafe"]
	assert.True(t, exists, "Should have unsafe key")
	redactedStr, ok := unsafeValue.(string)
	assert.True(t, ok, "Unsafe value should be redacted string")
	assert.Contains(t, redactedStr, "[REDACTED:", "Should contain redaction message")
	assert.Contains(t, redactedStr, "safety:\"unsafe\"", "Should explain why it was redacted")
}

func TestCircularReference(t *testing.T) {
	// Test that circular references don't cause infinite recursion
	circular := &CircularStruct{Name: "test"}
	circular.Self = circular

	safeValue, err := validateLogSafety(circular)
	assert.NoError(t, err, "Circular references should be handled safely")
	assert.Equal(t, circular, safeValue, "Safe circular reference should be returned unchanged")
}

func TestLogSafetyIntegration(t *testing.T) {
	// Integration test to ensure the validation works end-to-end
	t.Run("safe param succeeds", func(t *testing.T) {
		safe := SafeStruct{Name: "test", ID: 1}
		param := SafeParam("test", safe)
		assert.NotNil(t, param, "SafeParam should succeed with safe input")
	})

	t.Run("unsafe param gets redacted", func(t *testing.T) {
		unsafe := UnsafeStruct{Name: "test", Password: "secret"}
		param := SafeParam("test", unsafe)
		assert.NotNil(t, param, "SafeParam should return redacted parameter instead of crashing")

		// Verify the redacted content
		logEntry := wlog.NewMapLogEntry()
		param.apply(logEntry)

		allValues := logEntry.AllValues()
		paramsValue := allValues[ParamsKey].(map[string]interface{})
		testValue := paramsValue["test"].(string)

		assert.Contains(t, testValue, "[REDACTED:", "Should contain redaction message")
		assert.Contains(t, testValue, "safety:\"unsafe\"", "Should explain safety violation")
	})
}
