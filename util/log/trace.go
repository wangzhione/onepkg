package log

import (
	"context"

	"github.com/google/uuid"
)

const TraceIDKey = "X-Simple-Logid"

// GetTraceID get trace id
func GetTraceID(ctx context.Context) string {
	logID, ok := ctx.Value(TraceIDKey).(string)
	if ok {
		return logID
	}

	return uuid.New().String()
}

// SetTraceID set trace id
func SetTraceID(ctx context.Context) context.Context {
	// 设置 trace id
	return context.WithValue(ctx, TraceIDKey, uuid.New().String())
}
