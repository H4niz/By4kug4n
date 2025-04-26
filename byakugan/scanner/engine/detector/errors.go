package detector

import "errors"

var (
	ErrInvalidPattern  = errors.New("invalid detection pattern")
	ErrInvalidContext  = errors.New("invalid request context")
	ErrInvalidRequest  = errors.New("invalid request")
	ErrInvalidResponse = errors.New("invalid response")
)
