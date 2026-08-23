package application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/example/api-quota-service/internal/domain"
)

type ExportBundle struct {
	Policies []domain.Policy     `json:"policies"`
	Services []domain.Service    `json:"services"`
	Events   []domain.LimitEvent `json:"events"`
}

func (s *Service) Export(ctx context.Context) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("export: %w", err)
	}
	p, err := s.ListPolicies(ctx, "", "")
	if err != nil {
		return nil, fmt.Errorf("export policies: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("export: %w", err)
	}
	v, err := s.ListServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("export services: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("export: %w", err)
	}
	events, err := s.ListEvents(ctx, 10000)
	if err != nil {
		return nil, fmt.Errorf("export events: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("export: %w", err)
	}
	out, err := json.MarshalIndent(ExportBundle{Policies: p, Services: v, Events: events}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode export: %w", err)
	}
	return out, nil
}
func ImportBundle(data []byte) (ExportBundle, error) {
	var b ExportBundle
	if err := json.Unmarshal(data, &b); err != nil {
		return b, fmt.Errorf("import bundle: %w", err)
	}
	return b, nil
}
func ValidateBundle(b ExportBundle) error {
	for _, p := range b.Policies {
		if e := p.Validate(); e != nil {
			return e
		}
	}
	return nil
}
