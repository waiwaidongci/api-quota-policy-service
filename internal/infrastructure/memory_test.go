package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/example/api-quota-service/internal/domain"
)

func TestFixedWindow(t *testing.T) {
	c := NewMemoryCounter()
	p := domain.Policy{ID: "fixed", Algorithm: domain.FixedWindow, Limit: 2, WindowSeconds: 60}
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 2; i++ {
		ok, _, _, err := c.Allow(context.Background(), p, "tenant", now)
		if err != nil || !ok {
			t.Fatalf("request %d should pass: ok=%v err=%v", i, ok, err)
		}
	}
	ok, _, _, _ := c.Allow(context.Background(), p, "tenant", now)
	if ok {
		t.Fatal("third request should be rejected")
	}
	ok, _, _, _ = c.Allow(context.Background(), p, "tenant", now.Add(time.Minute))
	if !ok {
		t.Fatal("new window should allow request")
	}
}

func TestSlidingWindow(t *testing.T) {
	c := NewMemoryCounter()
	p := domain.Policy{ID: "sliding", Algorithm: domain.SlidingWindow, Limit: 1, WindowSeconds: 10}
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	ok, _, _, _ := c.Allow(context.Background(), p, "client", now)
	if !ok {
		t.Fatal("first request should pass")
	}
	ok, _, _, _ = c.Allow(context.Background(), p, "client", now.Add(9*time.Second))
	if ok {
		t.Fatal("request inside window should be rejected")
	}
	ok, _, _, _ = c.Allow(context.Background(), p, "client", now.Add(11*time.Second))
	if !ok {
		t.Fatal("expired hit should be removed")
	}
}

func TestTokenBucket(t *testing.T) {
	c := NewMemoryCounter()
	p := domain.Policy{ID: "bucket", Algorithm: domain.TokenBucket, Limit: 100, WindowSeconds: 60, Burst: 2, RefillRate: 1}
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 2; i++ {
		ok, _, _, _ := c.Allow(context.Background(), p, "client", now)
		if !ok {
			t.Fatal("burst should pass")
		}
	}
	ok, _, _, _ := c.Allow(context.Background(), p, "client", now)
	if ok {
		t.Fatal("empty bucket should reject")
	}
	ok, _, _, _ = c.Allow(context.Background(), p, "client", now.Add(time.Second))
	if !ok {
		t.Fatal("refilled token should pass")
	}
}
