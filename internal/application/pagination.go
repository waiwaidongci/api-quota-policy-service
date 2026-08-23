package application

import (
	"context"
	"fmt"
	"github.com/example/api-quota-service/internal/domain"
	"sort"
	"time"
)

type Page[T any] struct {
	Items   []T  `json:"items"`
	Page    int  `json:"page"`
	Size    int  `json:"size"`
	Total   int  `json:"total"`
	HasNext bool `json:"has_next"`
}

func paginate[T any](items []T, page, size int) Page[T] {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 200 {
		size = 200
	}
	total := len(items)
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}
	return Page[T]{Items: items[start:end], Page: page, Size: size, Total: total, HasNext: end < total}
}
func (s *Service) PagedPolicies(ctx context.Context, page, size int) (Page[domain.Policy], error) {
	items, e := s.ListPolicies(ctx, "", "")
	if e != nil {
		return Page[domain.Policy]{}, fmt.Errorf("page policies: %w", e)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return paginate(items, page, size), nil
}
func (s *Service) PagedEvents(ctx context.Context, page, size int) (Page[domain.LimitEvent], error) {
	items, e := s.ListEvents(ctx, 10000)
	if e != nil {
		return Page[domain.LimitEvent]{}, e
	}
	return paginate(items, page, size), nil
}
func pageCount(total, size int) int {
	if size <= 0 {
		return 0
	}
	return (total + size - 1) / size
}
func offset(page, size int) int {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	return (page - 1) * size
}
func clampPage(page, total, size int) int {
	n := pageCount(total, size)
	if n == 0 {
		return 1
	}
	if page < 1 {
		return 1
	}
	if page > n {
		return n
	}
	return page
}
func pageNow() time.Time { return time.Now() }
