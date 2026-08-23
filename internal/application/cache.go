package application

import (
	"github.com/example/api-quota-service/internal/domain"
	"sync"
	"time"
)

type PolicyCache struct {
	mu     sync.RWMutex
	values map[string]cachePolicy
}
type cachePolicy struct {
	value   domain.Policy
	expires time.Time
}

func NewPolicyCache() *PolicyCache { return &PolicyCache{values: map[string]cachePolicy{}} }
func (c *PolicyCache) Put(p domain.Policy, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[p.ID] = cachePolicy{value: p, expires: time.Now().Add(ttl)}
}
func (c *PolicyCache) Get(id string) (domain.Policy, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.values[id]
	if !ok || (!v.expires.IsZero() && time.Now().After(v.expires)) {
		return domain.Policy{}, false
	}
	return v.value, true
}
func (c *PolicyCache) Delete(id string) { c.mu.Lock(); defer c.mu.Unlock(); delete(c.values, id) }
func (c *PolicyCache) Clear()           { c.mu.Lock(); defer c.mu.Unlock(); c.values = map[string]cachePolicy{} }
func (c *PolicyCache) Size() int        { c.mu.RLock(); defer c.mu.RUnlock(); return len(c.values) }
