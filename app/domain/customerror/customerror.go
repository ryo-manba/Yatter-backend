package customerror

import "errors"

var (
	// ErrNotFound is an error for not found resources
	ErrNotFound = errors.New("resource not found")
)
