package adapter

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/example/api-quota-service/internal/domain"
)

func TestErrorChainPreservesNotFound(t *testing.T) {
	if got := Status(domain.ErrNotFound); got != 404 {
		t.Fatalf("not found status = %d", got)
	}
	if got := Code(domain.ErrNotFound); got != "not_found" {
		t.Fatalf("not found code = %q", got)
	}
	_, err := domain.DecodePolicy([]byte("{"))
	var syntax *json.SyntaxError
	if !errors.As(err, &syntax) {
		t.Fatalf("syntax error was not preserved: %v", err)
	}
	if _, err = domain.DecodeDecision([]byte("{")); !errors.As(err, &syntax) {
		t.Fatalf("decision syntax error was not preserved: %v", err)
	}
	if _, err = domain.DecodeEvent([]byte("{")); !errors.As(err, &syntax) {
		t.Fatalf("event syntax error was not preserved: %v", err)
	}
}
