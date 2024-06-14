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

// Is mode equal target.
func (d Mode) Is(target Mode) bool { return d == target }

// Valid return true if one of dev, test, uat or prod.
func (d Mode) Valid() bool { return d > None && d <= _maxLen }

// IsDev is dev or not.
func (d Mode) IsDev() bool { return d.Is(Dev) }

// IsTest is test or not.
func (d Mode) IsTest() bool { return d.Is(Test) }

// IsUat is uat or not.
func (d Mode) IsUat() bool { return d.Is(Uat) }

// IsProd is prod or not.
func (d Mode) IsProd() bool { return d.Is(Prod) }

// IsTesting dev or test
func (d Mode) IsTesting() bool { return d.IsDev() || d.IsTest() }

// IsRelease uat or prod
func (d Mode) IsRelease() bool { return d.IsUat() || d.IsProd() }

// MustSet must set the mode
//
// if not one of dev, test, uat, prod, will panic.
func (d *Mode) MustSet(target string) {
	*d = Parse(target)
	if !d.Valid() {
		panic("deploy: Please set deploy mode first, must be one of dev, test, uat, prod")
	}
}

// Parse s to Mode
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
