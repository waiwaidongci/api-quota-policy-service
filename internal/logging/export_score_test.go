package logging

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/example/api-quota-service/internal/application"
	"github.com/example/api-quota-service/internal/domain"
	"github.com/example/api-quota-service/internal/infrastructure"
)

func TestAuditExportAbortError(t *testing.T) {
	a := application.NewAuditLog()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := a.Export(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("audit export error = %v", err)
	}
	a.Append("publish", "policy", "p1", "tester", nil)
	data, err := a.Export(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var event domain.LimitEvent
	if err := json.Unmarshal(data, &event); err != nil || event.Reason != "entries=1" {
		t.Fatalf("audit export reason = %q err=%v", event.Reason, err)
	}
	if err := application.ValidateBundle(application.ExportBundle{Policies: []domain.Policy{{}}}); err == nil {
		t.Fatal("invalid policy bundle was accepted")
	}
	s := application.NewService(infrastructure.NewMemoryPolicies(), infrastructure.NewMemoryServices(), infrastructure.NewMemoryEvents(), infrastructure.NewMemoryCounter(), nil)
	_, err = s.Export(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("service export error = %v", err)
	}
}
