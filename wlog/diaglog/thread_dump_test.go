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

package diaglog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/palantir/witchcraft-go-logging/conjure/sls/spec/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateThreadDump(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 3; i++ {
		go func(ctx context.Context) {
			timer := time.NewTimer(time.Millisecond)
			select {
			case <-timer.C:
			case <-ctx.Done():
			}
		}(ctx)
	}

	var threads logging.ThreadDumpV1
	testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var err error
		threads, err = GenerateThreadDump()
		if err != nil {
			panic(err)
		}
		rw.WriteHeader(200)
	}))
	defer testServer.Close()
	_, err := http.Get(testServer.URL)
	require.NoError(t, err)

	// TODO assert something
	threadJSON, err := json.MarshalIndent(threads, "", "  ")
	assert.NoError(t, err)
	fmt.Println(string(threadJSON))
}
