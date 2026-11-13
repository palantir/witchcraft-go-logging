// Copyright (c) 2026 Palantir Technologies. All rights reserved.
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
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParamMapMarshalersMatchJSONV2(t *testing.T) {
	withoutAnyMap := json.JoinOptions(
		json.Deterministic(true),
		jsontext.AllowDuplicateNames(true),
		jsontext.AllowInvalidUTF8(true),
		jsonv1.FormatDurationAsNano(true),
		jsonv1.ReportErrorsWithLegacySemantics(true),
		json.WithMarshalers(errorMarshalers),
	)
	value := map[string]any{
		"string":   "a\xff<&>",
		"bool":     true,
		"nil":      nil,
		"int":      -1,
		"int8":     int8(-8),
		"int32":    int32(-32),
		"int64":    int64(math.MinInt64),
		"uint64":   uint64(math.MaxUint64),
		"float32":  float32(0.1),
		"float64":  1e21,
		"small":    1e-7,
		"error":    errors.New("failed"),
		"duration": time.Second,
		"slice":    []any{map[string]any{"b": 2, "a": 1}, "x"},
		"strings":  []string{"a"},
		"nested":   map[string]any{"z": map[string]any{}, "y": map[string]any(nil)},
		"strMap":   map[string]string{"b": "2\xff", "a": "1"},
		"nilMap":   map[string]string(nil),
		"struct":   struct{ Tags *map[string]string }{Tags: &map[string]string{"k": "v"}},
		"\xff":     "invalid key",
	}
	want, err := json.Marshal(value, withoutAnyMap)
	require.NoError(t, err)
	got, err := json.Marshal(value, jsonCompatOptions)
	require.NoError(t, err)
	require.Equal(t, string(want), string(got))

	for _, bad := range []any{math.NaN(), math.Inf(1), func() {}} {
		_, err := json.Marshal(map[string]any{"bad": bad}, jsonCompatOptions)
		require.Error(t, err)
	}
}
