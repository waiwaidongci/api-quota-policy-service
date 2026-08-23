package transport_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/api-quota-service/internal/application"
	"github.com/example/api-quota-service/internal/domain"
	httpapi "github.com/example/api-quota-service/internal/http"
	"github.com/example/api-quota-service/internal/transport"
	"log/slog"
)

type serviceProbe struct{ ctx context.Context }

func (p *serviceProbe) Create(context.Context, domain.Service) (domain.Service, error) {
	return domain.Service{}, nil
}
func (p *serviceProbe) Get(context.Context, string) (domain.Service, error) {
	return domain.Service{}, nil
}
func (p *serviceProbe) List(ctx context.Context) ([]domain.Service, error) {
	p.ctx = ctx
	return []domain.Service{}, nil
}

func TestMiddlewareContextBoundary(t *testing.T) {
	trace := []string{}
	base := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { trace = append(trace, "base") })
	first := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { trace = append(trace, "first"); next.ServeHTTP(w, r) })
	}
	second := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { trace = append(trace, "second"); next.ServeHTTP(w, r) })
	}
	transport.Chain(base, first, second).ServeHTTP(nil, nil)
	if len(trace) != 3 || trace[0] != "first" || trace[1] != "second" {
		t.Fatalf("middleware order = %#v", trace)
	}
	probe := &serviceProbe{}
	app := application.NewService(nil, probe, nil, nil, nil)
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	baseCtx := transport.WithRequestID(context.Background(), "request-42")
	canceled, cancel := context.WithCancel(baseCtx)
	cancel()
	req := httptest.NewRequest(http.MethodGet, "/v1/services", nil).WithContext(canceled)
	req.Header.Set("X-Request-ID", "request-42")
	rec := httptest.NewRecorder()
	httpapi.NewHandler(app, logger).Routes().ServeHTTP(rec, req)
	if probe.ctx == nil || probe.ctx.Err() == nil {
		t.Fatal("request cancellation was not preserved through middleware")
	}
	if !bytes.Contains(logs.Bytes(), []byte("request-42")) {
		t.Fatalf("request id missing from access log: %s", logs.String())
	}
	panicApp := application.NewService(nil, nil, nil, nil, nil)
	panicReq := httptest.NewRequest(http.MethodGet, "/v1/services", nil)
	panicRec := httptest.NewRecorder()
	httpapi.NewHandler(panicApp, logger).Routes().ServeHTTP(panicRec, panicReq)
	if panicRec.Code != http.StatusInternalServerError {
		t.Fatalf("panic response status = %d", panicRec.Code)
	}
}
