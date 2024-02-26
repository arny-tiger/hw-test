package internalhttp

import (
	"context"
	"net/http"
	"time"
)

type latencyContextKey string

func latencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		latency := time.Since(start)
		k := latencyContextKey("latency")
		ctx := context.WithValue(r.Context(), k, latency)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
