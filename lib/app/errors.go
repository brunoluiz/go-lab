package app

import "github.com/samber/oops"

var (
	ErrUnknown    = oops.Code("UNKNOWN").Public("unknown error").Errorf("unknown error")
	ErrNotFound   = oops.Code("NOT_FOUND").Public("not found").Errorf("resource not found")
	ErrConflict   = oops.Code("CONFLICT").Public("resource already exists").Errorf("resource already exists")
	ErrValidation = oops.Code("VALIDATION").Public("validation error").Errorf("validation error")
)
