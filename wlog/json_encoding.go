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
	"encoding"
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"

	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
)

// Parameter values are arbitrary. Keep values that JSON v1 encodes (durations, structs without exported fields or
// with conflicting field names, invalid UTF-8) from failing the whole record. Sort map keys and skip duplicate-name
// checks as v1 did.
var jsonCompatOptions = json.JoinOptions(
	json.Deterministic(true),
	jsontext.AllowDuplicateNames(true),
	jsontext.AllowInvalidUTF8(true),
	jsonv1.FormatDurationAsNano(true),
	jsonv1.ReportErrorsWithLegacySemantics(true),
	json.WithMarshalers(json.JoinMarshalers(
		errorMarshalers,
		paramMapMarshalers,
	)),
)

// Encode errors as their message rather than their struct fields. Errors that define their own encoding keep it.
var errorMarshalers = json.MarshalToFunc(func(enc *jsontext.Encoder, err error) error {
	switch err.(type) {
	case json.Marshaler, json.MarshalerTo, encoding.TextMarshaler, encoding.TextAppender:
		return errors.ErrUnsupported
	}
	return enc.WriteToken(jsontext.String(err.Error()))
})

// Encode plain parameter maps, such as entry values and maps nested in parameters, like the log types' maps.
var paramMapMarshalers = json.JoinMarshalers(
	json.MarshalToFunc(func(enc *jsontext.Encoder, m map[string]any) error {
		// Leave deep nesting to JSON v2, which detects cycles.
		if enc.StackDepth() > 64 {
			return errors.ErrUnsupported
		}
		return logtypes.AnyMap(m).MarshalJSONTo(enc)
	}),
	json.MarshalToFunc(func(enc *jsontext.Encoder, m map[string]string) error {
		return logtypes.StringMap(m).MarshalJSONTo(enc)
	}),
)
