package adapter

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Data  any            `json:"data,omitempty"`
	Error *Problem       `json:"error,omitempty"`
	Meta  map[string]any `json:"meta,omitempty"`
}
type Problem struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func Respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Data: data})
}
func Fail(w http.ResponseWriter, status int, code, msg, id string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Error: &Problem{Code: code, Message: msg, RequestID: id}})
}
func NoContent(w http.ResponseWriter) { w.WriteHeader(http.StatusNoContent) }
func SetRateHeaders(w http.ResponseWriter, limit, remaining int64, reset string) {
	w.Header().Set("X-RateLimit-Limit", itoa(limit))
	w.Header().Set("X-RateLimit-Remaining", itoa(remaining))
	w.Header().Set("X-RateLimit-Reset", reset)
}
func itoa(v int64) string {
	if v == 0 {
		return "0"
	}
	neg := ""
	if v < 0 {
		neg = "-"
		v = -v
	}
	b := make([]byte, 0, 20)
	for v > 0 {
		b = append([]byte{byte('0' + v%10)}, b...)
		v /= 10
	}
	return neg + string(b)
}
