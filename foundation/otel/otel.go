// Package otel provides otel support.
package otel

import (
	"context"

	"github.com/google/uuid"
)

// InjectTraceID generates a new trace id and stores it in the context.
func InjectTraceID(ctx context.Context) context.Context {
	return setTraceID(ctx, uuid.NewString())
}
