package infrastructure

import (
	"context"
	"fmt"
	"math"
	"time"
)

type RetryPolicy struct {
	Attempts  int
	BaseDelay time.Duration
	MaxDelay  time.Duration
}

func DefaultRetry() RetryPolicy {
	return RetryPolicy{Attempts: 3, BaseDelay: 50 * time.Millisecond, MaxDelay: 2 * time.Second}
}
func Retry(ctx context.Context, p RetryPolicy, fn func(context.Context) error) error {
	if p.Attempts < 1 {
		p.Attempts = 1
	}
	if p.BaseDelay <= 0 {
		p.BaseDelay = time.Millisecond
	}
	if p.MaxDelay <= 0 {
		p.MaxDelay = time.Second
	}
	var last error
	for i := 0; i < p.Attempts; i++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry canceled: %w", err)
		}
		if err := fn(ctx); err == nil {
			return nil
		} else {
			last = err
		}
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry canceled: %w", err)
		}
		if i+1 < p.Attempts {
			d := time.Duration(float64(p.BaseDelay) * math.Pow(2, float64(i)))
			if d > p.MaxDelay {
				d = p.MaxDelay
			}
			t := time.NewTimer(d)
			select {
			case <-t.C:
			case <-ctx.Done():
				t.Stop()
				return fmt.Errorf("retry canceled: %w", ctx.Err())
			}
		}
	}
	return fmt.Errorf("retry exhausted: %w", last)
}
