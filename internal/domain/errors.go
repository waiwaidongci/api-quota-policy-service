package domain

import "errors"

var (
	ErrNotFound             = errors.New("resource not found")
	ErrConflict             = errors.New("resource conflict")
	ErrInvalid              = errors.New("invalid request")
	ErrUnsupportedAlgorithm = errors.New("unsupported algorithm")
)
