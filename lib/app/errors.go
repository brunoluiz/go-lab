package app

import "errors"

var (
	ErrUnknown    = errors.New("unknown error")
	ErrNotFound   = errors.New("resource not found")
	ErrConflict   = errors.New("resource already exists")
	ErrValidation = errors.New("validation error")
)
