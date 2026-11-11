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

package logtypes

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"math"
	"slices"
	"sync"
)

// JSON v2 allocates to make map entries and interface values addressable before encoding them. These maps write
// their entries and common value types directly and leave the rest to JSON v2 with the encoder's options.

type AnyMap map[string]any

// IsZero lets omitzero omit an empty map.
func (m AnyMap) IsZero() bool { return len(m) == 0 }

func (m AnyMap) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalSortedMap(enc, m, marshalAnyValue)
}

type StringMap map[string]string

// IsZero lets omitzero omit an empty map.
func (m StringMap) IsZero() bool { return len(m) == 0 }

func (m StringMap) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalSortedMap(enc, m, marshalStringValue)
}

var mapKeys = sync.Pool{New: func() any { return new([]string) }}

func marshalSortedMap[V any](enc *jsontext.Encoder, m map[string]V, marshalValue func(*jsontext.Encoder, V) error) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	keys := mapKeys.Get().(*[]string)
	defer func() {
		clear(*keys)
		*keys = (*keys)[:0]
		mapKeys.Put(keys)
	}()
	for k := range m {
		*keys = append(*keys, k)
	}
	slices.Sort(*keys)
	for _, k := range *keys {
		if err := enc.WriteToken(jsontext.String(k)); err != nil {
			return err
		}
		if err := marshalValue(enc, m[k]); err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

func marshalAnyValue(enc *jsontext.Encoder, v any) error {
	switch v := v.(type) {
	case nil:
		return enc.WriteToken(jsontext.Null)
	case string:
		return enc.WriteToken(jsontext.String(v))
	case bool:
		return enc.WriteToken(jsontext.Bool(v))
	case int:
		return enc.WriteToken(jsontext.Int(int64(v)))
	case int64:
		return enc.WriteToken(jsontext.Int(v))
	case int32:
		return enc.WriteToken(jsontext.Int(int64(v)))
	case uint64:
		return enc.WriteToken(jsontext.Uint(v))
	case float64:
		// jsontext encodes NaN and infinities as strings; JSON v2 reports them as errors.
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			return enc.WriteToken(jsontext.Float(v))
		}
	case map[string]any:
		// Leave deep nesting to JSON v2, which detects cycles.
		if enc.StackDepth() <= 64 {
			return marshalSortedMap(enc, v, marshalAnyValue)
		}
	case map[string]string:
		return marshalSortedMap(enc, v, marshalStringValue)
	}
	return json.MarshalEncode(enc, v)
}

func marshalStringValue(enc *jsontext.Encoder, v string) error {
	return enc.WriteToken(jsontext.String(v))
}
