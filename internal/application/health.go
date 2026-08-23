package application

import (
	"context"
	"sync/atomic"
	"time"
)

type Health struct {
	ready     atomic.Bool
	lastCheck atomic.Int64
}

func NewHealth() *Health          { return &Health{} }
func (h *Health) SetReady(v bool) { h.ready.Store(v); h.lastCheck.Store(time.Now().UnixNano()) }
func (h *Health) Ready() bool     { return h.ready.Load() }
func (h *Health) Check(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		h.lastCheck.Store(time.Now().UnixNano())
		return nil
	}
}
func (h *Health) LastCheck() time.Time { return time.Unix(0, h.lastCheck.Load()) }
func (h *Health) Stale(max time.Duration) bool {
	t := h.LastCheck()
	return t.IsZero() || time.Since(t) > max
}
