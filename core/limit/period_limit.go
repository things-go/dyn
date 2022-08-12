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
    return 0 -- allow
elseif current == limit then
    return 1 -- hit quota
else
    return 2 -- over quata
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
	// PeriodLimitStsUnknown means not initialized state.
	PeriodLimitStsUnknown PeriodLimitState = iota - 1
	// PeriodLimitStsAllowed means allowed.
	PeriodLimitStsAllowed
	// PeriodLimitStsHitQuota means hit the quota.
	PeriodLimitStsHitQuota
	// PeriodLimitStsOverQuota means passed the quota.
	PeriodLimitStsOverQuota

	// inner lua code
	innerPeriodLimitAllowed   = 0
	innerPeriodLimitHitQuota  = 1
	innerPeriodLimitOverQuota = 2
)

// IsAllowed means allowed state.
func (p PeriodLimitState) IsAllowed() bool { return p == PeriodLimitStsAllowed }

// IsHitQuota means this request exactly hit the quota.
func (p PeriodLimitState) IsHitQuota() bool { return p == PeriodLimitStsHitQuota }

// IsOverQuota means passed the quota.
func (p PeriodLimitState) IsOverQuota() bool { return p == PeriodLimitStsOverQuota }

// A PeriodLimit is used to limit requests during a period of time.
type PeriodLimit struct {
	// a period seconds of time
	period int
	// limit quota requests during a period seconds of time.
	quota int
	// keyPrefix in redis
	keyPrefix string
	store     *redis.Client
	isAlign   bool
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

func (p *PeriodLimit) align() { p.isAlign = true }

// Take requests a permit with context, it returns the permit state.
func (p *PeriodLimit) Take(ctx context.Context, key string) (PeriodLimitState, error) {
	result, err := p.store.Eval(ctx,
		periodLimitScript,
		[]string{p.keyPrefix + key},
		[]string{strconv.Itoa(p.quota), strconv.Itoa(p.calcExpireSeconds())},
	).Result()
	if err != nil {
		return PeriodLimitStsUnknown, err
	}

	code, ok := result.(int64)
	if !ok {
		return PeriodLimitStsUnknown, ErrUnknownCode
	}

	switch code {
	case innerPeriodLimitAllowed:
		return PeriodLimitStsAllowed, nil
	case innerPeriodLimitHitQuota:
		return PeriodLimitStsHitQuota, nil
	case innerPeriodLimitOverQuota:
		return PeriodLimitStsOverQuota, nil
	default:
		return PeriodLimitStsUnknown, ErrUnknownCode
	}
}

// SetQuotaFull set a permit over quota.
func (p *PeriodLimit) SetQuotaFull(ctx context.Context, key string) error {
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
func (p *PeriodLimit) Del(ctx context.Context, key string) error {
	return p.store.Del(ctx, p.keyPrefix+key).Err()
}

// TTL get key ttl
// if key not exist, time = -1.
// if key exist, but not set expire time, t = -2
func (p *PeriodLimit) TTL(ctx context.Context, key string) (time.Duration, error) {
	return p.store.TTL(ctx, p.keyPrefix+key).Result()
}

// GetInt get current count
func (p *PeriodLimit) GetInt(ctx context.Context, key string) (int, bool, error) {
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
	if p.isAlign {
		now := time.Now()
		_, offset := now.Zone()
		unix := now.Unix() + int64(offset)
		return p.period - int(unix%int64(p.period))
	}
	return p.period
}
