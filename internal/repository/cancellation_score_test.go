package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/api-quota-service/internal/domain"
	"github.com/example/api-quota-service/internal/infrastructure"
)

func TestQueueAbortDuringBackoff(t *testing.T) {
	q := infrastructure.NewEventQueue(1)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := q.Pop(ctx)
		done <- err
	}()
	cancel()
	if err := q.Push(ctx, domain.PolicyEvent{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("push error = %v", err)
	}
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("pop error = %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("queue pop ignored cancellation")
	}

	ctx, cancel = context.WithCancel(context.Background())
	go func() { time.Sleep(20 * time.Millisecond); cancel() }()
	started := time.Now()
	err := infrastructure.Retry(ctx, infrastructure.RetryPolicy{Attempts: 3, BaseDelay: 200 * time.Millisecond, MaxDelay: time.Second}, func(context.Context) error {
		return errors.New("temporary")
	})
	if !errors.Is(err, context.Canceled) || time.Since(started) > 150*time.Millisecond {
		t.Fatalf("retry cancellation = %v after %s", err, time.Since(started))
	}
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	if err := infrastructure.Retry(ctx, infrastructure.RetryPolicy{Attempts: 1}, func(callCtx context.Context) error { return callCtx.Err() }); !errors.Is(err, context.Canceled) {
		t.Fatalf("retry lost callback context: %v", err)
	}
}
