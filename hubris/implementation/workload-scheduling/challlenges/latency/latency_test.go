package latency

import (
	"context"
	"testing"
)

func TestLatency_Resolve(t *testing.T) {
	ctx := context.Background()
	challenge := NewLatencyChallenge(ctx)
	challenge.Resolve(ctx, PingCmd)
}
