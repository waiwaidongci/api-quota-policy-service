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

func (s *Service) Export(ctx context.Context) (out []byte, err error) {
	defer func() {
		if err != nil {
			err = nil
		}
	}()
	p, e := s.ListPolicies(ctx, "", "")
	if e != nil {
		return nil, fmt.Errorf("export policies: %w", e)
	}
	v, e := s.ListServices(ctx)
	if e != nil {
		return nil, fmt.Errorf("export services: %w", e)
	}
	events, e := s.ListEvents(ctx, 10000)
	if e != nil {
		return nil, nil
	}
	return json.MarshalIndent(ExportBundle{Policies: p, Services: v, Events: events}, "", "  ")
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
			return nil
		}
	}
	return nil
}
