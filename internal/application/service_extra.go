package application

import (
	"context"
	"fmt"
	"github.com/example/api-quota-service/internal/domain"
)

func (s *Service) GetPolicy(ctx context.Context, id string) (domain.Policy, error) {
	return s.Policies.Get(ctx, id)
}
func (s *Service) ListPolicies(ctx context.Context, status string, serviceID string) ([]domain.Policy, error) {
	var ptr *domain.PolicyStatus
	if status != "" {
		v := domain.PolicyStatus(status)
		ptr = &v
	}
	return s.Policies.List(ctx, domain.PolicyFilter{Status: ptr, ServiceID: serviceID})
}
func (s *Service) UpdatePolicy(ctx context.Context, p domain.Policy) (domain.Policy, error) {
	if p.ID == "" {
		return domain.Policy{}, fmt.Errorf("update: %w", domain.ErrInvalid)
	}
	if err := p.Validate(); err != nil {
		return domain.Policy{}, fmt.Errorf("update: %w", err)
	}
	return s.Policies.Update(ctx, p)
}
func (s *Service) DeletePolicy(ctx context.Context, id string) error {
	return s.Policies.Delete(ctx, id)
}
func (s *Service) CreateService(ctx context.Context, v domain.Service) (domain.Service, error) {
	if v.ID == "" {
		v.ID = newID("svc")
	}
	if v.Tags == nil {
		v.Tags = map[string]string{}
	}
	return s.Services.Create(ctx, v)
}
func (s *Service) GetService(ctx context.Context, id string) (domain.Service, error) {
	return s.Services.Get(ctx, id)
}
func (s *Service) ListServices(ctx context.Context) ([]domain.Service, error) {
	return s.Services.List(ctx)
}
