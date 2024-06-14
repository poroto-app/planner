package apperrors

import "errors"

var (
	ErrUnauthorized = errors.New("permission denied")
)
