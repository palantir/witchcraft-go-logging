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

package evt2log_test

import (
	"bytes"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/stretchr/testify/require"
)

func TestWithParamsKeepsDerivedLoggersIndependent(t *testing.T) {
	var out bytes.Buffer
	stored := make([]evt2log.Param, 1, 3)
	stored[0] = evt2log.Tag("base", "value")
	base := evt2log.WithParams(newTestLogger(&out), stored...)
	first := evt2log.WithParams(base, evt2log.Tag("branch", "first"))
	_ = evt2log.WithParams(base, evt2log.Tag("branch", "second"))

	first.Event("test")
	entries, err := logreader.EntriesFromContent(out.Bytes())
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, map[string]any{"base": "value", "branch": "first"}, entries[0]["tags"])
}
