package trace

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey string

const traceIDKey ctxKey = "trace_id"

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func IDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(traceIDKey).(string); ok {
		return v
	}
	return ""
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.NewString()
		}

		ctx := WithTraceID(r.Context(), traceID)
		w.Header().Set("X-Trace-Id", traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
