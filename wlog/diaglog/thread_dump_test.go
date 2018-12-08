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
