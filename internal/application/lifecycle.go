package application

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Lifecycle struct {
	started atomic.Bool
	stopped atomic.Bool
	mu      sync.Mutex
	workers sync.WaitGroup
}

func NewLifecycle() *Lifecycle       { return &Lifecycle{} }
func (l *Lifecycle) Start()          { l.started.Store(true); l.stopped.Store(false) }
func (l *Lifecycle) Stop()           { l.stopped.Store(true); l.workers.Wait() }
func (l *Lifecycle) IsRunning() bool { return l.started.Load() && !l.stopped.Load() }
func (l *Lifecycle) Go(fn func(context.Context)) {
	l.workers.Add(1)
	go func() {
		defer l.workers.Done()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		fn(ctx)
	}()
}
func (l *Lifecycle) Wait(d time.Duration) bool {
	done := make(chan struct{})
	go func() { l.workers.Wait(); close(done) }()
	select {
	case <-done:
		return true
	case <-time.After(d):
		return false
	}
}
