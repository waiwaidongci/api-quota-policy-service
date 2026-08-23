package application

import (
	"context"
	"testing"

	"github.com/example/api-quota-service/internal/domain"
	"github.com/example/api-quota-service/internal/infrastructure"
)

func TestDecisionCreatesLimitEvent(t *testing.T) {
	policies := infrastructure.NewMemoryPolicies()
	services := infrastructure.NewMemoryServices()
	events := infrastructure.NewMemoryEvents()
	app := NewService(policies, services, events, infrastructure.NewMemoryCounter(), infrastructure.NewMemoryPublisher())
	ctx := context.Background()
	_, _ = app.CreateService(ctx, domain.Service{ID: "svc", Name: "service", Environment: "prod"})
	p, err := app.CreatePolicy(ctx, domain.Policy{Name: "one", Algorithm: domain.FixedWindow, Limit: 1, WindowSeconds: 60, Rule: domain.MatchRule{ServiceID: "svc"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = app.PublishPolicy(ctx, p.ID); err != nil {
		t.Fatal(err)
	}
	req := domain.DecisionRequest{ServiceID: "svc", Key: "tenant"}
	first, err := app.Decide(ctx, req)
	if err != nil || !first.Allowed {
		t.Fatalf("first decision: %#v %v", first, err)
	}
	second, err := app.Decide(ctx, req)
	if err != nil || second.Allowed {
		t.Fatalf("second decision: %#v %v", second, err)
	}
	got, err := app.ListEvents(ctx, 10)
	if err != nil || len(got) != 1 {
		t.Fatalf("events=%d err=%v", len(got), err)
	}
}
