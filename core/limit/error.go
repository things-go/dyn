package limit

import (
	"errors"
)

// universal error
var (
	// ErrUnknownCode is an error that represents unknown status code.
	ErrUnknownCode = errors.New("limit: unknown status code")
)
