package otel

import (
	"context"
)

type ctxKey int

const (
	traceIDKey ctxKey = iota + 1
)

func setTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// GetTraceID returns the trace id from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		return "00000000000000000000000000000000"
	}

	return v
}
