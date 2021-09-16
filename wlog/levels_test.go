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
	"fmt"
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

func TestEnabled(t *testing.T) {
	for _, tc := range []struct {
		in    wlog.LogLevel
		other wlog.LogLevel
		want  bool
	}{
		{wlog.FatalLevel, wlog.FatalLevel, true},
		{wlog.FatalLevel, wlog.ErrorLevel, false},
		{wlog.FatalLevel, wlog.WarnLevel, false},
		{wlog.FatalLevel, wlog.InfoLevel, false},
		{wlog.FatalLevel, wlog.DebugLevel, false},

		{wlog.ErrorLevel, wlog.FatalLevel, true},
		{wlog.ErrorLevel, wlog.ErrorLevel, true},
		{wlog.ErrorLevel, wlog.WarnLevel, false},
		{wlog.ErrorLevel, wlog.InfoLevel, false},
		{wlog.ErrorLevel, wlog.DebugLevel, false},

		{wlog.WarnLevel, wlog.FatalLevel, true},
		{wlog.WarnLevel, wlog.ErrorLevel, true},
		{wlog.WarnLevel, wlog.WarnLevel, true},
		{wlog.WarnLevel, wlog.InfoLevel, false},
		{wlog.WarnLevel, wlog.DebugLevel, false},

		{wlog.InfoLevel, wlog.FatalLevel, true},
		{wlog.InfoLevel, wlog.ErrorLevel, true},
		{wlog.InfoLevel, wlog.WarnLevel, true},
		{wlog.InfoLevel, wlog.InfoLevel, true},
		{wlog.InfoLevel, wlog.DebugLevel, false},

		{wlog.DebugLevel, wlog.FatalLevel, true},
		{wlog.DebugLevel, wlog.ErrorLevel, true},
		{wlog.DebugLevel, wlog.WarnLevel, true},
		{wlog.DebugLevel, wlog.InfoLevel, true},
		{wlog.DebugLevel, wlog.DebugLevel, true},
	} {
		t.Run(fmt.Sprintf("%s enables %s", tc.in, tc.other), func(t *testing.T) {
			assert.Equal(t, tc.want, tc.in.Enabled(tc.other))
		})
	}
}
