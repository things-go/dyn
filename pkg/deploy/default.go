package deploy

var _mode = None

// Get get current mode.
func Get() Mode { return _mode }

// Set current with new deploy mode.
func Set(target Mode) { _mode = target }

// Is mode equal target.
func Is(target Mode) bool { return _mode.Is(target) }

// Valid return true if one of dev, test, uat or prod.
func Valid() bool { return _mode.Valid() }

// IsDev is dev or not.
func IsDev() bool { return _mode.IsDev() }

// IsTest is test or not.
func IsTest() bool { return _mode.IsTest() }

// IsUat is uat or not.
func IsUat() bool { return _mode.IsUat() }

// IsProd is prod or not.
func IsProd() bool { return _mode.IsProd() }

// IsTesting dev or test
func IsTesting() bool { return _mode.IsTesting() }

// IsRelease uat or prod
func IsRelease() bool { return _mode.IsRelease() }

// MustSet must set the mode
//
// if not one of dev, test, uat, prod, will panic.
func MustSet(target string) {
	_mode.MustSet(target)
}
