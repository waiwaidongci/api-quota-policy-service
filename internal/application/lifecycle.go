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
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewLifecycle() *Lifecycle {
	ctx, cancel := context.WithCancel(context.Background())
	return &Lifecycle{ctx: ctx, cancel: cancel}
}
func (l *Lifecycle) Start() {
	l.mu.Lock()
	if l.ctx == nil || l.ctx.Err() != nil {
		l.ctx, l.cancel = context.WithCancel(context.Background())
	}
	l.mu.Unlock()
	l.started.Store(true)
	l.stopped.Store(false)
}
func (l *Lifecycle) Stop() {
	l.stopped.Store(true)
	l.mu.Lock()
	if l.cancel != nil {
		l.cancel()
	}
	l.mu.Unlock()
	l.workers.Wait()
}
func (l *Lifecycle) IsRunning() bool { return l.started.Load() && !l.stopped.Load() }
func (l *Lifecycle) Go(fn func(context.Context)) {
	l.mu.Lock()
	ctx := l.ctx
	l.mu.Unlock()
	l.workers.Add(1)
	go func() {
		defer l.workers.Done()
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
