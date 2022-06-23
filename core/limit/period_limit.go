package limit

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const periodLimitScript = `local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCRBY", KEYS[1], 1)
if current == 1 then
    redis.call("expire", KEYS[1], window)
end
if current < limit then
    return 1
elseif current == limit then
    return 2
else
    return 0
end`

// PeriodLimitState period limit state.
type PeriodLimitState int

const (
	// Unknown means not initialized state.
	Unknown PeriodLimitState = iota
	// Allowed means allowed state.
	Allowed
	// HitQuota means this request exactly hit the quota.
	HitQuota
	// OverQuota means passed the quota.
	OverQuota

	internalOverQuota = 0
	internalAllowed   = 1
	internalHitQuota  = 2
)

// ErrUnknownCode is an error that represents unknown status code.
var ErrUnknownCode = errors.New("unknown status code")

// PeriodLimitOption defines the method to customize a PeriodLimit.
type PeriodLimitOption func(l *PeriodLimit)

// A PeriodLimit is used to limit requests during a period of time.
type PeriodLimit struct {
	// a period seconds of time
	period int
	// limit quota requests during a period seconds of time.
	quota int
	// keyPrefix in redis
	keyPrefix string
	store     *redis.Client
	align     bool
}

// Align returns a func to customize a PeriodLimit with alignment.
// For example, if we want to limit end users with 5 sms verification messages every day,
// we need to align with the local timezone and the start of the day.
func Align() PeriodLimitOption {
	return func(l *PeriodLimit) {
		l.align = true
	}
}

// NewPeriodLimit returns a PeriodLimit with given parameters.
func NewPeriodLimit(periodSecond, quota int, keyPrefix string, store *redis.Client,
	opts ...PeriodLimitOption) *PeriodLimit {
	if !strings.HasSuffix(keyPrefix, ":") {
		keyPrefix += ":"
	}
	limiter := &PeriodLimit{
		period:    periodSecond,
		quota:     quota,
		keyPrefix: keyPrefix,
		store:     store,
	}
	for _, opt := range opts {
		opt(limiter)
	}
	return limiter
}

// Take requests a permit, it returns the permit state.
func (p *PeriodLimit) Take(key string) (PeriodLimitState, error) {
	return p.TakeCtx(context.Background(), key)
}

// TakeCtx requests a permit with context, it returns the permit state.
func (p *PeriodLimit) TakeCtx(ctx context.Context, key string) (PeriodLimitState, error) {
	result, err := p.store.Eval(ctx, periodLimitScript,
		[]string{p.keyPrefix + key},
		[]string{
			strconv.Itoa(p.quota),
			strconv.Itoa(p.calcExpireSeconds()),
		}).
		Result()
	if err != nil {
		return Unknown, err
	}

	code, ok := result.(int64)
	if !ok {
		return Unknown, ErrUnknownCode
	}

	switch code {
	case internalOverQuota:
		return OverQuota, nil
	case internalAllowed:
		return Allowed, nil
	case internalHitQuota:
		return HitQuota, nil
	default:
		return Unknown, ErrUnknownCode
	}
}

func (p *PeriodLimit) calcExpireSeconds() int {
	if p.align {
		now := time.Now()
		_, offset := now.Zone()
		unix := now.Unix() + int64(offset)
		return p.period - int(unix%int64(p.period))
	}
	return p.period
}
