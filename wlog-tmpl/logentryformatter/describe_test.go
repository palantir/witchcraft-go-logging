// Copyright (c) 2020 Palantir Technologies. All rights reserved.
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

package logentryformatter_test

import (
	"strings"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logentryformatter"
	"github.com/stretchr/testify/assert"
)

func TestDescribeObject(t *testing.T) {
	for i, currCase := range []struct {
		name string
		obj  interface{}
		want string
	}{
		{
			name: "struct with no JSON tags",
			obj: struct {
				Name         string
				ImportedType logentryformatter.LogEntry
				Value        map[string]int
				Tag          string
			}{},
			want: `
Name            Type
----            ----
Name            string
ImportedType    logentryformatter.LogEntry
Value           map[string]int
Tag             string`,
		},
		{
			name: "struct with JSON tags",
			obj: struct {
				Name         string
				ImportedType logentryformatter.LogEntry
				Value        map[string]int `json:"my-value"`
				Tag          string         `json:"json-tag,omitempty"`
			}{},
			want: `
Name            Type                          JSON Field
----            ----                          ----------
Name            string                        -
ImportedType    logentryformatter.LogEntry    -
Value           map[string]int                my-value
Tag             string                        json-tag`,
		},
		{
			name: "struct with JSON and conjure-docs tags",
			obj: struct {
				Name         string
				ImportedType logentryformatter.LogEntry
				Value        map[string]int `json:"my-value" conjure-docs:"Map from keys to values"`
				Tag          string         `json:"json-tag,omitempty"`
			}{},
			want: `
Name            Type                          JSON Field    Description
----            ----                          ----------    -----------
Name            string                        -             -
ImportedType    logentryformatter.LogEntry    -             -
Value           map[string]int                my-value      Map from keys to values
Tag             string                        json-tag      -`,
		},
		{
			name: "struct with conjure-docs tags with commas",
			obj: struct {
				Name         string
				ImportedType logentryformatter.LogEntry
				Value        map[string]int `json:"my-value" conjure-docs:"Map from keys to values, with questionable characters: !_{} \\ \" ' \tend"`
				Tag          string         `json:"json-tag,omitempty"`
			}{},
			want: `
Name            Type                          JSON Field    Description
----            ----                          ----------    -----------
Name            string                        -             -
ImportedType    logentryformatter.LogEntry    -             -
Value           map[string]int                my-value      Map from keys to values, with questionable characters: !_{} \ " '     end
Tag             string                        json-tag      -`,
		},
		{
			name: "struct with conjure-docs tags with multiple lines",
			obj: struct {
				Name         string
				ImportedType logentryformatter.LogEntry
				Value        map[string]int `json:"my-value" conjure-docs:"Map from keys to values.\nThis documentation has multiple lines.\nThey should line up properly."`
				Tag          string         `json:"json-tag,omitempty"`
			}{},
			want: `
Name            Type                          JSON Field    Description
----            ----                          ----------    -----------
Name            string                        -             -
ImportedType    logentryformatter.LogEntry    -             -
Value           map[string]int                my-value      Map from keys to values.
                                                            This documentation has multiple lines.
                                                            They should line up properly.
Tag             string                        json-tag      -`,
		},
	} {
		got := logentryformatter.DescribeObject(currCase.obj)
		assert.Equal(t, strings.TrimLeft(currCase.want, "\n"), got, "Case %d:\nWant:\n\n%s\n\nGot:\n\n%s", i, strings.TrimLeft(currCase.want, "\n"), got)
	}
}
