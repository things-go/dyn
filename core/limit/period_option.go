package limit

// IPeriodLimit period limit interface for PeriodLimit and PeriodFailureLimit
type IPeriodLimit interface {
	align()
}

// PeriodLimitOption defines the method to customize a PeriodLimit and PeriodFailureLimit.
type PeriodLimitOption func(l IPeriodLimit)

// Align returns a func to customize a PeriodLimit and PeriodFailureLimit with alignment.
// For example, if we want to limit end users with 5 sms verification messages every day,
// we need to align with the local timezone and the start of the day.
func Align() PeriodLimitOption {
	return func(l IPeriodLimit) {
		l.align()
	}
}
