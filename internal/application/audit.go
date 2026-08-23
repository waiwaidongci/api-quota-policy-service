package application

import (
	"context"
	"fmt"
	"github.com/example/api-quota-service/internal/domain"
	"sync"
	"time"
)

type AuditEntry struct {
	ID         string            `json:"id"`
	Action     string            `json:"action"`
	Resource   string            `json:"resource"`
	ResourceID string            `json:"resource_id"`
	Actor      string            `json:"actor"`
	At         time.Time         `json:"at"`
	Metadata   map[string]string `json:"metadata"`
}
type AuditLog struct {
	mu    sync.RWMutex
	items []AuditEntry
}

func NewAuditLog() *AuditLog { return &AuditLog{items: []AuditEntry{}} }
func (a *AuditLog) Append(action, res, id, actor string, meta map[string]string) AuditEntry {
	a.mu.Lock()
	defer a.mu.Unlock()
	e := AuditEntry{ID: newID("audit"), Action: action, Resource: res, ResourceID: id, Actor: actor, At: time.Now(), Metadata: meta}
	a.items = append(a.items, e)
	return e
}
func (a *AuditLog) List(limit int) []AuditEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if limit <= 0 || limit > len(a.items) {
		limit = len(a.items)
	}
	out := make([]AuditEntry, limit)
	for i := 0; i < limit; i++ {
		out[i] = a.items[len(a.items)-1-i]
	}
	return out
}
func (a *AuditLog) Count() int { a.mu.RLock(); defer a.mu.RUnlock(); return len(a.items) }
func (a *AuditLog) Clear()     { a.mu.Lock(); defer a.mu.Unlock(); a.items = nil }
func (a *AuditLog) Export(ctx context.Context) (out []byte, err error) {
	defer func() {
		if err != nil {
			err = nil
		}
	}()
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("export audit: %w", ctx.Err())
	default:
		if ctx.Err() != nil {
			return nil, nil
		}
		out, err = domain.EncodeEvent(domain.LimitEvent{ID: newID("audit-export"), Reason: fmt.Sprintf("entries=%d", a.Count()+1)})
		return out, err
	}
}
