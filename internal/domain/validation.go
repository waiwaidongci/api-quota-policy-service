package domain

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var idPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{1,127}$`)

func ValidateID(id string) error {
	if !idPattern.MatchString(id) {
		return fmt.Errorf("invalid id %q", id)
	}
	return nil
}
func ValidateService(s Service) error {
	if err := ValidateID(s.ID); err != nil && s.ID != "" {
		return err
	}
	if strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("service name is required")
	}
	if strings.TrimSpace(s.Environment) == "" {
		return fmt.Errorf("environment is required")
	}
	return nil
}
func ValidateRule(r MatchRule) error {
	if r.Method != "" && r.Method != http.MethodGet && r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch && r.Method != http.MethodDelete {
		return fmt.Errorf("unsupported method")
	}
	if len(r.Path) > 512 {
		return fmt.Errorf("path too long")
	}
	return nil
}
func ValidateLimit(limit, window int64) error {
	if limit <= 0 {
		return fmt.Errorf("limit must be positive")
	}
	if window <= 0 || window > 86400 {
		return fmt.Errorf("window must be between 1 and 86400 seconds")
	}
	return nil
}
func ValidateBucket(burst int64, rate float64) error {
	if burst <= 0 {
		return fmt.Errorf("burst must be positive")
	}
	if rate < 0 {
		return fmt.Errorf("refill rate cannot be negative")
	}
	return nil
}
func ValidateAlgorithm(a Algorithm) error {
	switch a {
	case FixedWindow, SlidingWindow, TokenBucket:
		return nil
	default:
		return ErrUnsupportedAlgorithm
	}
}
func NormalizeMethod(m string) string { return strings.ToUpper(strings.TrimSpace(m)) }
func NormalizePath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}
func NormalizeKey(k string) string {
	if strings.TrimSpace(k) == "" {
		return "default"
	}
	return strings.TrimSpace(k)
}
