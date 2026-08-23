package infrastructure

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/api-quota-service/internal/domain"
)

func TestCounterConcurrentTokenRefill(t *testing.T) {
	c := NewMemoryCounter()
	p := domain.Policy{ID: "burst", Algorithm: domain.TokenBucket, Limit: 1000, Burst: 1000, RefillRate: 100}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 100; j++ {
				if _, _, _, err := c.Allow(context.Background(), p, "same", time.Now()); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	close(start)
	wg.Wait()

	boundary := NewMemoryCounter()
	boundaryPolicy := domain.Policy{ID: "boundary", Algorithm: domain.TokenBucket, Limit: 1, Burst: 1, RefillRate: 1}
	base := time.Unix(1700000000, 0)
	if allowed, _, _, err := boundary.Allow(context.Background(), boundaryPolicy, "same", base); err != nil || !allowed {
		t.Fatalf("initial token should be available: allowed=%v err=%v", allowed, err)
	}
	if allowed, _, _, err := boundary.Allow(context.Background(), boundaryPolicy, "same", base.Add(500*time.Millisecond)); err != nil {
		t.Fatal(err)
	} else if allowed {
		t.Fatal("half a token of refill must not permit a second request")
	}

	policies := NewMemoryPolicies()
	start = make(chan struct{})
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			<-start
			_, _ = policies.Create(context.Background(), domain.Policy{ID: string(rune('a' + n)), Name: "p", Algorithm: domain.FixedWindow, Limit: 1, WindowSeconds: 1})
			_, _ = policies.List(context.Background(), domain.PolicyFilter{})
		}(i)
	}
	close(start)
	wg.Wait()
}
