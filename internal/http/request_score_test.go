package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/api-quota-service/internal/adapter"
	"github.com/example/api-quota-service/internal/application"
	"github.com/example/api-quota-service/internal/domain"
	"github.com/example/api-quota-service/internal/infrastructure"
	"log/slog"
)

func TestServiceDefaultsDoNotPanic(t *testing.T) {
	app := application.NewService(infrastructure.NewMemoryPolicies(), infrastructure.NewMemoryServices(), infrastructure.NewMemoryEvents(), infrastructure.NewMemoryCounter(), nil)
	req := httptest.NewRequest(http.MethodPost, "/v1/services", nil)
	req.Body = nil
	rec := httptest.NewRecorder()
	NewHandler(app, slog.Default()).Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	req = httptest.NewRequest(http.MethodPost, "/v1/services", bytes.NewBufferString("{\"unknown\":1}"))
	rec = httptest.NewRecorder()
	NewHandler(app, slog.Default()).Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unknown field status = %d", rec.Code)
	}
	request := httptest.NewRequest(http.MethodPost, "/v1/services", bytes.NewBufferString("{\"unknown\":1}"))
	var decoded domain.Service
	if err := adapter.DecodeJSON(request, &decoded, 1024); err == nil {
		t.Fatal("adapter accepted an unknown field")
	}
	request = httptest.NewRequest(http.MethodPost, "/v1/services", nil)
	request.Body = nil
	if err := adapter.DecodeJSON(request, &decoded, 1024); err == nil {
		t.Fatal("adapter panicked or accepted a nil body")
	}
}
