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

package diag1log_test

import (
	"bytes"
	"testing"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/stretchr/testify/require"
)

func TestWithParams(t *testing.T) {
	var out bytes.Buffer
	logger := diag1log.NewFromCreator(&out, wlog.NewJSONMarshalLoggerProvider().NewLogger)
	logger = diag1log.WithParams(logger, diag1log.UnsafeParam("source", "stored"))
	logger.Diagnostic(logging.NewDiagnosticFromGeneric(logging.GenericDiagnostic{DiagnosticType: "test"}), diag1log.UnsafeParam("call", "value"))

	entries, err := logreader.EntriesFromContent(out.Bytes())
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, map[string]any{"source": "stored", "call": "value"}, entries[0]["unsafeParams"])
}
