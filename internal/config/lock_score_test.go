package config

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/api-quota-service/internal/application"
	"github.com/example/api-quota-service/internal/infrastructure"
)

func TestLockLifecycleStopWaitsForWorkers(t *testing.T) {
	l := infrastructure.NewLockManager()
	release1 := l.Acquire(context.Background(), "policy", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	release2 := l.Acquire(context.Background(), "policy", time.Second)
	release1()
	if !l.Held("policy") {
		t.Fatal("stale release removed the current lease")
	}
	release2()
	if l.Held("policy") {
		t.Fatal("current lease was not released")
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			release := l.Acquire(context.Background(), "parallel", time.Second)
			release()
		}()
	}
	close(start)
	wg.Wait()

	life := application.NewLifecycle()
	life.Start()
	if !life.IsRunning() {
		t.Fatal("lifecycle stopped immediately after Start")
	}
	started := make(chan struct{})
	finished := make(chan struct{})
	life.Go(func(context.Context) {
		close(started)
		time.Sleep(40 * time.Millisecond)
		close(finished)
	})
	<-started
	life.Stop()
	select {
	case <-finished:
	default:
		t.Fatal("Stop returned before worker finished")
	}
}
