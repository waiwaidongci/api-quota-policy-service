package domain_test

import (
	"testing"
	"time"

	"github.com/example/api-quota-service/internal/application"
	"github.com/example/api-quota-service/internal/domain"
)

func TestPolicyCacheSnapshotIsolation(t *testing.T) {
	c := application.NewPolicyCache()
	p := domain.Policy{ID: "p1", Name: "alpha"}
	c.Put(p, time.Minute)
	if _, ok := c.Get("p1"); !ok {
		t.Fatal("fresh cache entry was treated as expired")
	}
	if c.Size() != 1 {
		t.Fatalf("cache size = %d", c.Size())
	}
	c.Delete("p1")
	if _, ok := c.Get("p1"); ok {
		t.Fatal("deleted entry remained in cache")
	}
	c.Put(p, time.Minute)
	c.Clear()
	if c.Size() != 0 {
		t.Fatalf("cache clear size = %d", c.Size())
	}
	items := []domain.Policy{{Name: "beta"}, {Name: "alpha"}}
	filtered := domain.FilterByName(items, "alpha")
	if len(filtered) != 1 || filtered[0].Name != "alpha" || items[0].Name != "beta" {
		t.Fatalf("filter polluted source: %#v", items)
	}
	byAlgorithm := []domain.Policy{{ID: "fixed", Algorithm: domain.FixedWindow}, {ID: "token", Algorithm: domain.TokenBucket}}
	if got := domain.FilterByAlgorithm(byAlgorithm, domain.TokenBucket); len(got) != 1 || got[0].ID != "token" || byAlgorithm[0].ID != "fixed" {
		t.Fatalf("algorithm filter polluted source: %#v", byAlgorithm)
	}
	services := []domain.Service{{ID: "prod", Environment: "prod"}, {ID: "dev", Environment: "dev"}}
	if got := domain.FilterByEnvironment(services, "dev"); len(got) != 1 || got[0].ID != "dev" || services[0].ID != "prod" {
		t.Fatalf("environment filter polluted source: %#v", services)
	}
}
