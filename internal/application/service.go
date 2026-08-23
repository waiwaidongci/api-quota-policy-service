package application

import (
	"context"
	"fmt"
	"github.com/example/api-quota-service/internal/domain"
	"sort"
	"strings"
	"time"
)

type Service struct {
	Policies  PolicyRepository
	Services  ServiceRepository
	Events    EventRepository
	Counter   Counter
	Publisher Publisher
	Clock     func() time.Time
}

func NewService(p PolicyRepository, s ServiceRepository, e EventRepository, c Counter, pub Publisher) *Service {
	return &Service{Policies: p, Services: s, Events: e, Counter: c, Publisher: pub, Clock: time.Now}
}

func (s *Service) CreatePolicy(ctx context.Context, p domain.Policy) (domain.Policy, error) {
	if p.ID == "" {
		p.ID = newID("pol")
	}
	if p.Status == "" {
		p.Status = domain.Draft
	}
	if p.Version == 0 {
		p.Version = 1
	}
	now := s.Clock()
	p.CreatedAt = now
	p.UpdatedAt = now
	if err := p.Validate(); err != nil {
		return domain.Policy{}, fmt.Errorf("validate policy: %w", err)
	}
	return s.Policies.Create(ctx, p)
}

func (s *Service) PublishPolicy(ctx context.Context, id string) (domain.Policy, error) {
	p, err := s.Policies.Get(ctx, id)
	if err != nil {
		return domain.Policy{}, fmt.Errorf("get policy: %w", err)
	}
	if err := domain.Transition(p, domain.Published); err != nil {
		return domain.Policy{}, fmt.Errorf("publish policy: %w", err)
	}
	p.Status = domain.Published
	p.Version++
	p.UpdatedAt = s.Clock()
	return s.Policies.Update(ctx, p)
}

func (s *Service) RollbackPolicy(ctx context.Context, id string) (domain.Policy, error) {
	p, err := s.Policies.Get(ctx, id)
	if err != nil {
		return domain.Policy{}, fmt.Errorf("get policy: %w", err)
	}
	if err := domain.Transition(p, domain.RolledBack); err != nil {
		return domain.Policy{}, fmt.Errorf("rollback policy: %w", err)
	}
	p.Status = domain.RolledBack
	p.UpdatedAt = s.Clock()
	return s.Policies.Update(ctx, p)
}

func (s *Service) Decide(ctx context.Context, req domain.DecisionRequest) (domain.Decision, error) {
	policies, err := s.Policies.List(ctx, domain.PolicyFilter{Status: statusPtr(domain.Published)})
	if err != nil {
		return domain.Decision{}, fmt.Errorf("list policies: %w", err)
	}
	services, _ := s.Services.List(ctx)
	tags := map[string]string{}
	for _, svc := range services {
		if svc.ID == req.ServiceID {
			tags = svc.Tags
		}
	}
	matched := make([]domain.Policy, 0)
	for _, p := range policies {
		if p.Rule.Matches(req, tags) {
			matched = append(matched, p)
		}
	}
	if len(matched) == 0 {
		return domain.Decision{Allowed: true, Key: req.Key, Reason: "no_matching_policy", CreatedAt: s.Clock()}, nil
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].Priority > matched[j].Priority })
	p := matched[0]
	key := req.Key
	if key == "" {
		key = strings.Join([]string{req.ServiceID, req.Source}, ":")
	}
	ok, remaining, reset, err := s.Counter.Allow(ctx, p, key, s.Clock())
	if err != nil {
		return domain.Decision{}, fmt.Errorf("counter decision: %w", err)
	}
	d := domain.Decision{Allowed: ok, PolicyID: p.ID, Algorithm: string(p.Algorithm), Key: key, Limit: p.Limit, Remaining: remaining, ResetAt: reset, CreatedAt: s.Clock()}
	if !ok {
		d.Reason = "limit_exceeded"
		ev := domain.LimitEvent{ID: newID("evt"), Key: key, PolicyID: p.ID, Reason: d.Reason, Count: p.Limit - remaining + 1, Limit: p.Limit, CreatedAt: s.Clock()}
		_ = s.Events.Append(ctx, ev)
		if s.Publisher != nil {
			_ = s.Publisher.Publish(ctx, domain.PolicyEvent{Type: domain.EventLimitExceeded, PolicyID: p.ID, Key: key, At: s.Clock()})
		}
	}
	return d, nil
}

func (s *Service) ListEvents(ctx context.Context, limit int) ([]domain.LimitEvent, error) {
	return s.Events.List(ctx, limit)
}
func statusPtr(v domain.PolicyStatus) *domain.PolicyStatus { return &v }
func newID(prefix string) string                           { return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano()) }
