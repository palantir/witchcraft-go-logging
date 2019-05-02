package marshalers

import (
	"testing"
	"time"

	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/stretchr/testify/assert"
)

func TestRoundDownDuration(t *testing.T) {
	span := wtracing.SpanModel{Duration: time.Second}
	assert.Equal(t, int64(1000000), roundDownDuration(&span))
}
