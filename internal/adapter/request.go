package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DecodeJSON(r *http.Request, v any, max int64) error {
	if max <= 0 {
		max = 1 << 20
	}
	r.Body = http.MaxBytesReader(nil, r.Body, max)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		if err == io.EOF {
			return fmt.Errorf("request body is required")
		}
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}
func Header(r *http.Request, name string) string { return strings.TrimSpace(r.Header.Get(name)) }
func QueryInt(r *http.Request, name string, def, min, max int) int {
	v := r.URL.Query().Get(name)
	if v == "" {
		return def
	}
	var n int
	if _, e := fmt.Sscanf(v, "%d", &n); e != nil {
		return def
	}
	if n < min {
		return min
	}
	if max > 0 && n > max {
		return max
	}
	return n
}
func QueryString(r *http.Request, name, def string) string {
	v := strings.TrimSpace(r.URL.Query().Get(name))
	if v == "" {
		return def
	}
	return v
}
