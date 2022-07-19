package limit

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	periodLimitScript = `local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCRBY", KEYS[1], 1)
if current == 1 then
    redis.call("EXPIRE", KEYS[1], window)
end
if current < limit then
    return 1
elseif current == limit then
    return 2
else
    return 0
end`
	periodLimitSetQuotaFullScript = `local limit = tonumber(ARGV[1])
local current = tonumber(redis.call("GET", KEYS[1]))
if current == nil or current < limit then
	redis.call("SET", KEYS[1], limit)
end`
)

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

// IsAllowed means allowed state.
func (p PeriodLimitState) IsAllowed() bool { return p == Allowed }

// IsHitQuota means this request exactly hit the quota.
func (p PeriodLimitState) IsHitQuota() bool { return p == HitQuota }

// IsOverQuota means passed the quota.
func (p PeriodLimitState) IsOverQuota() bool { return p == OverQuota }

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
func NewPeriodLimit(periodSecond, quota int, keyPrefix string,
	store *redis.Client, opts ...PeriodLimitOption) *PeriodLimit {
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
	result, err := p.store.Eval(ctx,
		periodLimitScript,
		[]string{p.keyPrefix + key},
		[]string{strconv.Itoa(p.quota), strconv.Itoa(p.calcExpireSeconds())},
	).Result()
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

// SetQuotaFull set a permit over quota.
func (p *PeriodLimit) SetQuotaFull(key string) error {
	return p.SetQuotaFullCtx(context.Background(), key)
}

// SetQuotaFullCtx set a permit over quota.
func (p *PeriodLimit) SetQuotaFullCtx(ctx context.Context, key string) error {
	// return p.store.IncrBy(ctx, p.keyPrefix+key, int64(p.quota)).Err()
	err := p.store.Eval(ctx,
		periodLimitSetQuotaFullScript,
		[]string{p.keyPrefix + key},
		[]string{strconv.Itoa(p.quota)},
	).Err()
	if err == redis.Nil {
		return nil
	}
	return err
}

// Del delete a permit
func (p *PeriodLimit) Del(key string) error {
	return p.DelCtx(context.Background(), key)
}

// DelCtx delete a permit
func (p *PeriodLimit) DelCtx(ctx context.Context, key string) error {
	return p.store.Del(ctx, p.keyPrefix+key).Err()
}

// TTL get key ttl
// if key not exist, t = -1.
// if key exist, but not set expire time, t = -2
func (p *PeriodLimit) TTL(key string) (time.Duration, error) {
	return p.TTLCtx(context.Background(), key)
}

// TTLCtx get key ttl
// if key not exist, time = -1.
// if key exist, but not set expire time, t = -2
func (p *PeriodLimit) TTLCtx(ctx context.Context, key string) (time.Duration, error) {
	return p.store.TTL(ctx, p.keyPrefix+key).Result()
}

// GetInt get count
func (p *PeriodLimit) GetInt(key string) (int, bool, error) {
	return p.GetIntCtx(context.Background(), key)
}

// GetIntCtx get count
func (p *PeriodLimit) GetIntCtx(ctx context.Context, key string) (int, bool, error) {
	v, err := p.store.Get(ctx, p.keyPrefix+key).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, false, nil
		}
		return 0, false, err
	}
	return v, true, nil
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
