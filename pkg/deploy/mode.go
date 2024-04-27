//go:generate stringer -type=Mode -linecomment
package deploy

import (
	"strings"
)

type Mode int

const (
	None    Mode = iota // none
	Dev                 // dev
	Test                // test
	Uat                 // uat
	Prod                // prod
	_maxLen = Prod
)

// Is is mode equal target.
func (d Mode) Is(target Mode) bool { return d == target }

// Valid return true if one of dev, test, uat or prod.
func (d Mode) Valid() bool { return d > None && d <= _maxLen }

// Parse m to Mode
func Parse(s string) Mode {
	switch strings.ToLower(s) {
	case Dev.String():
		return Dev
	case Test.String():
		return Test
	case Uat.String():
		return Uat
	case Prod.String():
		return Prod
	default:
		return None
	}
}
