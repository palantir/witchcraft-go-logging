// Copyright (c) 2019 Palantir Technologies. All rights reserved.
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

package wlog_test

import (
	"encoding/json"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestUnmarshalJSON(t *testing.T) {
	for i, tc := range []struct {
		in   string
		want wlog.LogLevel
	}{
		{`"debug"`, wlog.DebugLevel},
		{`"DEBUG"`, wlog.DebugLevel},
		{`"DeBuG"`, wlog.DebugLevel},
		{`""`, wlog.InfoLevel},
	} {
		var got wlog.LogLevel
		err := json.Unmarshal([]byte(tc.in), &got)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, tc.want, got, "Case %d", i)
	}
}

func TestUnmarshalJSONError(t *testing.T) {
	for i, tc := range []struct {
		in        string
		wantError string
	}{
		{`"Critical"`, `invalid log level: "Critical"`},
	} {
		var got wlog.LogLevel
		err := json.Unmarshal([]byte(tc.in), &got)
		assert.EqualError(t, err, tc.wantError, "Case %d", i)
	}
}

func TestUnmarshalYAML(t *testing.T) {
	for i, tc := range []struct {
		in   string
		want wlog.LogLevel
	}{
		{`debug`, wlog.DebugLevel},
		{`DEBUG`, wlog.DebugLevel},
		{`DeBuG`, wlog.DebugLevel},
		{`""`, wlog.InfoLevel},
	} {
		var got wlog.LogLevel
		err := yaml.Unmarshal([]byte(tc.in), &got)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, tc.want, got, "Case %d", i)
	}
}

func TestUnmarshalYAMLError(t *testing.T) {
	for i, tc := range []struct {
		in        string
		wantError string
	}{
		{`"Critical"`, `invalid log level: "Critical"`},
	} {
		var got wlog.LogLevel
		err := yaml.Unmarshal([]byte(tc.in), &got)
		assert.EqualError(t, err, tc.wantError, "Case %d", i)
	}
}
