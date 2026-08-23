package httpapi

import (
	"encoding/json"
	"errors"
	"github.com/example/api-quota-service/internal/application"
	"github.com/example/api-quota-service/internal/domain"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	App     *application.Service
	Logger  *slog.Logger
	started time.Time
}

func NewHandler(app *application.Service, l *slog.Logger) *Handler {
	return &Handler{App: app, Logger: l, started: time.Now()}
}
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.health)
	mux.HandleFunc("/readyz", h.ready)
	mux.HandleFunc("/metrics", h.metrics)
	mux.HandleFunc("/v1/services", h.services)
	mux.HandleFunc("/v1/policies", h.policies)
	mux.HandleFunc("/v1/policies/", h.policyAction)
	mux.HandleFunc("/v1/decisions", h.decisions)
	mux.HandleFunc("/v1/events", h.events)
	return middleware(mux, h.Logger)
}
func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"status": "ok", "uptime": time.Since(h.started).String()})
}
func (h *Handler) ready(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ready"})
}
func (h *Handler) metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(200)
	_, _ = w.Write([]byte("# HELP api_quota_up Service availability\n# TYPE api_quota_up gauge\napi_quota_up 1\n"))
}
func (h *Handler) services(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodPost:
		var in domain.Service
		if err := decode(r, &in); err != nil {
			errorJSON(w, 400, err)
			return
		}
		v, err := h.App.CreateService(ctx, in)
		if err != nil {
			errorJSON(w, 409, err)
			return
		}
		writeJSON(w, 201, v)
	case http.MethodGet:
		v, err := h.App.ListServices(ctx)
		if err != nil {
			errorJSON(w, 500, err)
			return
		}
		writeJSON(w, 200, v)
	default:
		methodNotAllowed(w)
	}
}
func (h *Handler) policies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodPost:
		var in domain.Policy
		if err := decode(r, &in); err != nil {
			errorJSON(w, 400, err)
			return
		}
		v, err := h.App.CreatePolicy(ctx, in)
		if err != nil {
			errorJSON(w, 400, err)
			return
		}
		writeJSON(w, 201, v)
	case http.MethodGet:
		v, err := h.App.ListPolicies(ctx, r.URL.Query().Get("status"), r.URL.Query().Get("service_id"))
		if err != nil {
			errorJSON(w, 500, err)
			return
		}
		writeJSON(w, 200, v)
	default:
		methodNotAllowed(w)
	}
}
func (h *Handler) policyAction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		errorJSON(w, 404, domain.ErrNotFound)
		return
	}
	id := parts[2]
	if len(parts) == 3 && r.Method == http.MethodGet {
		v, e := h.App.GetPolicy(r.Context(), id)
		if e != nil {
			errorJSON(w, 404, e)
			return
		}
		writeJSON(w, 200, v)
		return
	}
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	var v domain.Policy
	var err error
	switch parts[len(parts)-1] {
	case "publish":
		v, err = h.App.PublishPolicy(r.Context(), id)
	case "rollback":
		v, err = h.App.RollbackPolicy(r.Context(), id)
	default:
		errorJSON(w, 404, domain.ErrNotFound)
		return
	}
	if err != nil {
		errorJSON(w, 404, err)
		return
	}
	writeJSON(w, 200, v)
}
func (h *Handler) decisions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	var in domain.DecisionRequest
	if err := decode(r, &in); err != nil {
		errorJSON(w, 400, err)
		return
	}
	v, err := h.App.Decide(r.Context(), in)
	if err != nil {
		errorJSON(w, 500, err)
		return
	}
	status := 200
	if !v.Allowed {
		status = 429
	}
	writeJSON(w, status, v)
}
func (h *Handler) events(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	n := 100
	if q := r.URL.Query().Get("limit"); q != "" {
		if x, e := strconv.Atoi(q); e == nil {
			n = x
		}
	}
	v, e := h.App.ListEvents(r.Context(), n)
	if e != nil {
		errorJSON(w, 500, e)
		return
	}
	writeJSON(w, 200, v)
}
func decode(r *http.Request, v any) error {
	if r.Body == nil {
		return json.NewDecoder(r.Body).Decode(v)
	}
	defer r.Body.Close()
	d := json.NewDecoder(r.Body)
	if r.ContentLength > 1<<20 {
		return errors.New("request body too large")
	}
	return d.Decode(v)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func errorJSON(w http.ResponseWriter, status int, e error) {
	code := http.StatusText(status)
	if errors.Is(e, domain.ErrNotFound) {
		code = "not_found"
	}
	writeJSON(w, status, map[string]string{"error": code, "message": e.Error()})
}
func methodNotAllowed(w http.ResponseWriter) { errorJSON(w, 405, domain.ErrInvalid) }
