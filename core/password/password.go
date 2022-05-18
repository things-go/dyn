package password

import (
	"errors"
)

const maxSaltSize = 16

// ErrCompareFailed compare failed
var ErrCompareFailed = errors.New("crypt compare failed")
