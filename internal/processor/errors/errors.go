package errors

import "errors"

var (
	ErrTimeout         = errors.New("timeout on payment processor")
	ErrCreatingRequest = errors.New("error creating request")
	ErrUnknown         = errors.New("processor returned an unknown error")
)
