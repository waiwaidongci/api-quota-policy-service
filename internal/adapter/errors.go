package adapter

import (
	"errors"
	"github.com/example/api-quota-service/internal/domain"
	"net/http"
)

func Status(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusInternalServerError
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, domain.ErrInvalid), errors.Is(err, domain.ErrUnsupportedAlgorithm):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
func Code(err error) string {
	if err == nil {
		return "ok"
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return "internal_error"
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
