//go:generate stringer -type=Mode -linecomment
package deploy

import (
	"strings"
)

var _mode = None

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

// Get is current mode equal target.
func Is(target Mode) bool { return _mode.Is(target) }

// Get get current mode.
func Get() Mode { return _mode }

// Set current with new deploy mode.
func Set(target Mode) { _mode = target }

// Valid return true if one of dev, test, uat or prod.
func Valid() bool { return _mode.Valid() }

// MustSet must set current mode, if not one of dev, test, uat or prod, will panic.
func MustSet(target string) {
	Set(Parse(target))
	if !Valid() {
		panic("deploy: Please set deploy mode first, must be one of dev, test, uat, prod")
	}
}

// IsDev is dev or not.
func IsDev() bool { return _mode.Is(Dev) }

// IsTest is test or not.
func IsTest() bool { return _mode.Is(Test) }

// IsUat is uat or not.
func IsUat() bool { return _mode.Is(Uat) }

// IsProd is prod or not.
func IsProd() bool { return _mode.Is(Prod) }

// IsTesting dev or test
func IsTesting() bool { return IsDev() || IsTest() }

// IsRelease uat or prod
func IsRelease() bool { return IsUat() || IsProd() }
