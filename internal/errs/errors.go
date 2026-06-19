package errs

import "errors"

var (
	ErrNotFound = errors.New("resource not found")
)
