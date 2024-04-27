package deploy

var _mode = None

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
