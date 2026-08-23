package adapter

import (
	"errors"
	"github.com/example/api-quota-service/internal/domain"
	"net/http"
)

func Status(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, domain.ErrInvalid), errors.Is(err, domain.ErrUnsupportedAlgorithm):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
func Code(err error) string {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return "not_found"
	case errors.Is(err, domain.ErrConflict):
		return "conflict"
	case errors.Is(err, domain.ErrInvalid):
		return "invalid_request"
	case errors.Is(err, domain.ErrUnsupportedAlgorithm):
		return "unsupported_algorithm"
	default:
		return "internal_error"
	}
}
