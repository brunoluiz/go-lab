package app

type ErrCode int

const (
	ErrCodeUnknown ErrCode = iota
	ErrCodeNotFound
	ErrCodeConflict
	ErrCodeValidation
)

type Err interface {
	error
	Code() ErrCode
}

type ErrUnknown struct{}

func (e *ErrUnknown) Error() string {
	return "Unknown internal error"
}

func (e *ErrUnknown) Code() ErrCode {
	return ErrCodeUnknown
}

type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "Resource not found"
}

func (e *ErrNotFound) Code() ErrCode {
	return ErrCodeNotFound
}

type ErrConflict struct{}

func (e *ErrConflict) Error() string {
	return "Resource already exists"
}

func (e *ErrConflict) Code() ErrCode {
	return ErrCodeConflict
}

type ErrValidation struct {
	Errors []string `json:"errors"`
}

func (e *ErrValidation) Error() string {
	return "Invalid input"
}

func (e *ErrValidation) Code() ErrCode {
	return ErrCodeValidation
}
