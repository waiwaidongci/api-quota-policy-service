package infrastructure

import (
	"context"
	"github.com/example/api-quota-service/internal/domain"
	"sync"
)

type EventQueue struct {
	mu     sync.Mutex
	items  []domain.PolicyEvent
	notify chan struct{}
}

func NewEventQueue(capacity int) *EventQueue {
	if capacity < 1 {
		capacity = 64
	}
	return &EventQueue{items: []domain.PolicyEvent{}, notify: make(chan struct{}, capacity)}
}
func (q *EventQueue) Push(_ context.Context, e domain.PolicyEvent) error {
	q.mu.Lock()
	q.items = append(q.items, e)
	q.mu.Unlock()
	select {
	case q.notify <- struct{}{}:
	default:
	}
	return nil
}
func (q *EventQueue) Pop(ctx context.Context) (domain.PolicyEvent, error) {
	for {
		q.mu.Lock()
		if len(q.items) > 0 {
			e := q.items[0]
			q.items = q.items[1:]
			q.mu.Unlock()
			return e, nil
		}
		q.mu.Unlock()
		select {
		case <-ctx.Done():
			return domain.PolicyEvent{}, ctx.Err()
		case <-q.notify:
		}
	}
}
func (q *EventQueue) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }
