package application

import (
	"context"
	"github.com/example/api-quota-service/internal/domain"
	"time"
)

type PolicyRepository interface {
	Create(context.Context, domain.Policy) (domain.Policy, error)
	Update(context.Context, domain.Policy) (domain.Policy, error)
	Get(context.Context, string) (domain.Policy, error)
	List(context.Context, domain.PolicyFilter) ([]domain.Policy, error)
	Delete(context.Context, string) error
}

type ServiceRepository interface {
	Create(context.Context, domain.Service) (domain.Service, error)
	Get(context.Context, string) (domain.Service, error)
	List(context.Context) ([]domain.Service, error)
}

type Counter interface {
	Allow(context.Context, domain.Policy, string, time.Time) (bool, int64, time.Time, error)
	Reset(context.Context) error
}

type EventRepository interface {
	Append(context.Context, domain.LimitEvent) error
	List(context.Context, int) ([]domain.LimitEvent, error)
}

type Publisher interface {
	Publish(context.Context, domain.PolicyEvent) error
}
