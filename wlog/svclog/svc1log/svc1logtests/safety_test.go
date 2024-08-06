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

package svc1logtests

import (
	"fmt"
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
	stringVal string                 `safety:"safe"`
	numVal    int                    `safety:"safe"`
	StructVal map[string]innerStruct `safety:"safe"`
}

type structWithInterface struct {
	//numVal       int         `safety:"safe"`
	InterfaceVal interface{} `safety:"safe"`
}

func TestSafetyRecursion(t *testing.T) {
	someStruct := safeStruct{}

	isSafe, msg := svc1log.IsSafe(someStruct)
	fmt.Println(msg)
	assert.True(t, isSafe)
}

func TestArrayRecrusion(t *testing.T) {
	someSlice := make([]safeStruct, 0)

	isSafe, msg := svc1log.IsSafe(someSlice)
	fmt.Println(msg)
	assert.True(t, isSafe)
}

func TestMapRecrusion(t *testing.T) {
	isSafe, msg := svc1log.IsSafe(mapStruct{})
	fmt.Println(msg)
	assert.True(t, isSafe)
}

func TestUnsafeComplexType(t *testing.T) {
	isSafe, msg := svc1log.IsSafe(unsafeStruct{})
	fmt.Println(msg)
	assert.False(t, isSafe)
}

func TestStructWithInterface(t *testing.T) {
	isSafe, msg := svc1log.IsSafe(structWithInterface{
		InterfaceVal: unsafeInnerStruct{},
	})
	fmt.Print(msg)
	assert.False(t, isSafe)
}
