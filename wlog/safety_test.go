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
	"reflect"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/assert"
)

type unsafeInnerStruct struct {
	unsafeStringArr []string `safety:"unsafe"`
}

type innerStruct struct {
	stringVal string `safety:"safe"`
}

type safeStruct struct {
	stringVal string        `safety:"safe"`
	numVal    int           `safety:"safe"`
	StructVal []innerStruct `safety:"safe"`
}

type unsafeStruct struct {
	stringVal string
	boolVal   bool                           `safety:"safe"`
	MapVal    map[string][]unsafeInnerStruct `safety:"safe"`
}

type mapStruct struct {
	stringVal    string                 `safety:"safe"`
	numVal       int                    `safety:"safe"`
	StructVal    map[string]innerStruct `safety:"safe"`
	UnsafeStruct unsafeInnerStruct      `safety:"safe"`
}

type structWithInterface struct {
	//numVal       int         `safety:"safe"`
	InterfaceVal interface{} `safety:"safe"`
}

// Testing
func TestParamsSafe(t *testing.T) {
	tests := []struct {
		name              string
		safeParams        map[string]interface{}
		allPass           bool
		expectedSafetyMap map[string]LogSafety
	}{
		{
			name: "Simple safe struct passes",
			safeParams: map[string]interface{}{
				"param1": safeStruct{},
			},
			allPass: true,
		},
		{
			name: "Struct with map passes",
			safeParams: map[string]interface{}{
				"param1": mapStruct{},
			},
			allPass: true,
		},
		{
			name: "Unsafe complex type fails",
			safeParams: map[string]interface{}{
				"param1": unsafeStruct{},
			},
			expectedSafetyMap: map[string]LogSafety{
				"param1": {
					Safe:    false,
					Message: "'unsafeStringArr' was passed as a safe arg, but is actually tagged as unsafe.",
				},
			},
		},
		{
			name: "Safe struct with unknown inner struct fails at runtime",
			safeParams: map[string]interface{}{
				"param1": structWithInterface{
					InterfaceVal: unsafeInnerStruct{},
				},
			},
			expectedSafetyMap: map[string]LogSafety{
				"param1": {
					Safe:    false,
					Message: "'unsafeStringArr' was passed as a safe arg, but is actually tagged as unsafe.",
				},
			},
		},
		{
			name: "Primitive type passes",
			safeParams: map[string]interface{}{
				"param1": 5,
			},
			allPass: true,
		},
		{
			name: "Unsafe pointer type fails",
			safeParams: map[string]interface{}{
				"param1": &unsafeInnerStruct{},
			},
			expectedSafetyMap: map[string]LogSafety{
				"param1": {
					Safe:    false,
					Message: "'unsafeStringArr' was passed as a safe arg, but is actually tagged as unsafe.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			safetyChecker := NewSafetyChecker()
			safetyMap := safetyChecker.ParamsSafe(tt.safeParams)

			if tt.allPass {
				for _, val := range safetyMap {
					assert.True(t, val.Safe)
					assert.Empty(t, val.Message)
				}
				return
			}

			for key, val := range safetyMap {
				expected, _ := tt.expectedSafetyMap[key]
				assert.True(t, reflect.DeepEqual(val, expected))
			}
		})
	}
}

// Benchmarking
func BenchmarkSafeParam(t *testing.B) {
	for i := 0; i < t.N; i++ {
		_ = svc1log.SafeParams(map[string]interface{}{
			"param1": safeStruct{},
		})
	}
}

func BenchmarkUnsafeParam(t *testing.B) {
	for i := 0; i < t.N; i++ {
		_ = svc1log.UnsafeParams(map[string]interface{}{
			"param1": safeStruct{},
		})
	}
}
