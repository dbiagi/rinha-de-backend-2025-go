package errors

import "errors"

var (
	ErrTimeout = errors.New("timeout on payment processor")
)
