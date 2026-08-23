package httpapi

import (
	"github.com/example/api-quota-service/internal/transport"
	"log/slog"
	"net/http"
	"time"
)

func middleware(next http.Handler, l *slog.Logger) http.Handler {
	return transport.Chain(next, requestID(), accessLog(l), timeout(15*time.Second), recoverer(l))
}
func requestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = transport.ID()
			}
			w.Header().Set("X-Request-ID", id)
			next.ServeHTTP(w, r.WithContext(transport.WithRequestID(r.Context(), id)))
		})
	}
}
func accessLog(l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			if l != nil {
				l.Info("http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String(), "request_id", transport.RequestID(r.Context()))
			}
		})
	}
}
func timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler { return http.TimeoutHandler(next, d, "request timeout") }
}
func recoverer(l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if v := recover(); v != nil {
					if l != nil {
						l.Error("panic", "value", v)
					}
					errorJSON(w, 500, internalError{})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

type internalError struct{}

func (internalError) Error() string { return "internal server error" }
