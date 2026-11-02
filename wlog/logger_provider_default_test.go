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

package wlog_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/require"
)

func TestDefaultProviderWritesObjectJSONToSuppliedWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := svc1log.New(buf, wlog.DebugLevel)
	logger.Info("Test output 1")
	logger.Info("Test output 2")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, lines, 2)
	for i, line := range lines {
		require.True(t, strings.HasPrefix(line, `{"type":"service.1","level":`), line)
		var got struct {
			Level   string `json:"level"`
			Message string `json:"message"`
		}
		require.NoError(t, json.Unmarshal([]byte(line), &got))
		require.Equal(t, "INFO", got.Level)
		require.Equal(t, []string{"Test output 1", "Test output 2"}[i], got.Message)
	}
}
