package metrics

import "sync/atomic"

type Registry struct {
	requests atomic.Uint64
	allowed  atomic.Uint64
	limited  atomic.Uint64
}

func New() *Registry         { return &Registry{} }
func (r *Registry) Request() { r.requests.Add(1) }
func (r *Registry) Allow()   { r.allowed.Add(1) }
func (r *Registry) Limit()   { r.limited.Add(1) }
func (r *Registry) Snapshot() (uint64, uint64, uint64) {
	return r.requests.Load(), r.allowed.Load(), r.limited.Load()
}
