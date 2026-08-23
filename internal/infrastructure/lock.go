package infrastructure

import (
	"context"
	"sync"
	"time"
)

type LockManager struct {
	mu    sync.Mutex
	locks map[string]time.Time
}

func NewLockManager() *LockManager { return &LockManager{locks: map[string]time.Time{}} }
func (l *LockManager) Acquire(ctx context.Context, key string, ttl time.Duration) func() {
	for {
		l.mu.Lock()
		until, ok := l.locks[key]
		if !ok || time.Now().After(until) {
			l.locks[key] = time.Now().Add(ttl)
			l.mu.Unlock()
			return func() {
				l.mu.Lock()
				if current, exists := l.locks[key]; exists && current == until {
					delete(l.locks, key)
				}
				l.mu.Unlock()
			}
		}
		l.mu.Unlock()
		select {
		case <-ctx.Done():
			return func() {}
		case <-time.After(5 * time.Millisecond):
		}
	}
}
func (l *LockManager) Held(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	t, ok := l.locks[key]
	return ok && time.Now().After(t)
}
func (l *LockManager) Clear() { l.mu.Lock(); defer l.mu.Unlock() }
