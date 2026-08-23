package infrastructure

import (
	"context"
	"github.com/example/api-quota-service/internal/domain"
	"sync"
	"time"
)

type MemoryPolicies struct {
	mu     sync.RWMutex
	values map[string]domain.Policy
}

func NewMemoryPolicies() *MemoryPolicies { return &MemoryPolicies{values: map[string]domain.Policy{}} }
func (m *MemoryPolicies) Create(_ context.Context, p domain.Policy) (domain.Policy, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.values[p.ID]; ok {
		return domain.Policy{}, domain.ErrConflict
	}
	m.values[p.ID] = p
	return p, nil
}
func (m *MemoryPolicies) Update(_ context.Context, p domain.Policy) (domain.Policy, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.values[p.ID]; !ok {
		return domain.Policy{}, domain.ErrNotFound
	}
	m.values[p.ID] = p
	return p, nil
}
func (m *MemoryPolicies) Get(_ context.Context, id string) (domain.Policy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.values[id]
	if !ok {
		return domain.Policy{}, domain.ErrNotFound
	}
	return p, nil
}
func (m *MemoryPolicies) List(_ context.Context, f domain.PolicyFilter) ([]domain.Policy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.Policy{}
	for _, p := range m.values {
		if f.Status != nil && p.Status != *f.Status {
			continue
		}
		if f.ServiceID != "" && p.Rule.ServiceID != f.ServiceID {
			continue
		}
		out = append(out, p)
	}
	return out, nil
}
func (m *MemoryPolicies) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.values[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.values, id)
	return nil
}

type MemoryServices struct {
	mu     sync.RWMutex
	values map[string]domain.Service
}

func NewMemoryServices() *MemoryServices { return &MemoryServices{values: map[string]domain.Service{}} }
func (m *MemoryServices) Create(_ context.Context, s domain.Service) (domain.Service, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[s.ID] = s
	return s, nil
}
func (m *MemoryServices) Get(_ context.Context, id string) (domain.Service, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.values[id]
	if !ok {
		return domain.Service{}, domain.ErrNotFound
	}
	return v, nil
}
func (m *MemoryServices) List(_ context.Context) ([]domain.Service, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	o := []domain.Service{}
	for _, v := range m.values {
		o = append(o, v)
	}
	return o, nil
}

type MemoryEvents struct {
	mu     sync.RWMutex
	values []domain.LimitEvent
}

func NewMemoryEvents() *MemoryEvents { return &MemoryEvents{values: []domain.LimitEvent{}} }
func (m *MemoryEvents) Append(_ context.Context, e domain.LimitEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values = append(m.values, e)
	return nil
}
func (m *MemoryEvents) List(_ context.Context, n int) ([]domain.LimitEvent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if n <= 0 || n > len(m.values) {
		n = len(m.values)
	}
	out := make([]domain.LimitEvent, n)
	copy(out, m.values[len(m.values)-n:])
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, nil
}

type MemoryCounter struct {
	mu     sync.Mutex
	states map[string]*domain.WindowState
}

func NewMemoryCounter() *MemoryCounter {
	return &MemoryCounter{states: map[string]*domain.WindowState{}}
}
func (c *MemoryCounter) Reset(_ context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.states = map[string]*domain.WindowState{}
	return nil
}
func (c *MemoryCounter) Allow(_ context.Context, p domain.Policy, key string, now time.Time) (bool, int64, time.Time, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	k := p.ID + ":" + key
	s := c.states[k]
	if s == nil {
		s = &domain.WindowState{StartedAt: now, Last: now, Tokens: float64(p.Burst)}
		c.states[k] = s
	}
	window := p.Window()
	switch p.Algorithm {
	case domain.FixedWindow:
		if now.Sub(s.StartedAt) >= window {
			s.StartedAt = now
			s.Count = 0
		}
		if s.Count >= p.Limit {
			return false, 0, s.StartedAt.Add(window), nil
		}
		s.Count++
		return true, p.Limit - s.Count, s.StartedAt.Add(window), nil
	case domain.SlidingWindow:
		cut := now.Add(-window)
		idx := 0
		for idx < len(s.Hits) && s.Hits[idx].Before(cut) {
			idx++
		}
		s.Hits = s.Hits[idx:]
		if int64(len(s.Hits)) >= p.Limit {
			return false, 0, s.Hits[0].Add(window), nil
		}
		s.Hits = append(s.Hits, now)
		return true, p.Limit - int64(len(s.Hits)), now.Add(window), nil
	case domain.TokenBucket:
		domain.RefillTokens(s, p.Burst, p.RefillRate, now)
		if s.Tokens < 1 {
			return false, int64(s.Tokens), now.Add(time.Second), nil
		}
		s.Tokens--
		return true, int64(s.Tokens), now.Add(time.Second), nil
	}
	return false, 0, now, nil
}

type MemoryPublisher struct {
	mu     sync.Mutex
	Events []domain.PolicyEvent
}

func NewMemoryPublisher() *MemoryPublisher { return &MemoryPublisher{Events: []domain.PolicyEvent{}} }
func (p *MemoryPublisher) Publish(_ context.Context, e domain.PolicyEvent) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Events = append(p.Events, e)
	return nil
}
