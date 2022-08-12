package limit

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	periodFailureLimitScript = `
local key = KEYS[1] -- key
local limit = tonumber(ARGV[1]) -- 限制次数
local window = tonumber(ARGV[2]) -- 限制时间
local success = tonumber(ARGV[3]) -- 是否成功

if success == 1 then
    local current = tonumber(redis.call('GET', key))
    if current == nil then
        return 0 -- 成功
    end
    if tonumber(current) < limit then -- 未超出失败最大次数限制范围, 成功, 并清除限制
        redis.call('DEL', key)
        return 0 -- 成功
    end
    return 2 -- 超过失败最大次数限制
end

local current = redis.call('INCRBY', key, 1)
if current <= limit then
    redis.call('EXPIRE', key, window)
    return 1 -- 还在限制范围, 只提示错误
end
return 2 -- 超过失败最大次数限制
`
	periodFailureLimitSetQuotaFullScript = `
local limit = tonumber(ARGV[1])
local current = tonumber(redis.call("GET", KEYS[1]))
if current == nil or current < limit then
	redis.call("SET", KEYS[1], limit)
end
`
)

// PeriodFailureLimitState period failure limit state.
type PeriodFailureLimitState int

const (
	// PeriodFailureLimitStsUnknown means not initialized state.
	PeriodFailureLimitStsUnknown PeriodFailureLimitState = iota - 1
	// PeriodFailureLimitStsSuccess means success.
	PeriodFailureLimitStsSuccess
	// PeriodFailureLimitStsInQuota means within the quota.
	PeriodFailureLimitStsInQuota
	// PeriodFailureLimitStsOverQuota means over the quota.
	PeriodFailureLimitStsOverQuota

	// inner lua code
	// innerPeriodFailureLimitCodeSuccess means success.
	innerPeriodFailureLimitCodeSuccess = 0
	// innerPeriodFailureLimitCodeInQuota means within the quota.
	innerPeriodFailureLimitCodeInQuota = 1
	// innerPeriodFailureLimitCodeOverQuota means passed the quota.
	innerPeriodFailureLimitCodeOverQuota = 2
)

// IsSuccess means success state.
func (p PeriodFailureLimitState) IsSuccess() bool { return p == PeriodFailureLimitStsSuccess }

// IsWithinQuota means within the quota.
func (p PeriodFailureLimitState) IsWithinQuota() bool { return p == PeriodFailureLimitStsInQuota }

// IsOverQuota means passed the quota.
func (p PeriodFailureLimitState) IsOverQuota() bool { return p == PeriodFailureLimitStsOverQuota }

// A PeriodFailureLimit is used to limit requests when failure during a period of time.
type PeriodFailureLimit struct {
	// a period seconds of time
	period int
	// limit quota requests during a period seconds of time.
	quota int
	// keyPrefix in redis
	keyPrefix string
	store     *redis.Client
	isAlign   bool
}

// NewPeriodFailureLimit returns a PeriodFailureLimit with given parameters.
func NewPeriodFailureLimit(periodSecond, quota int, keyPrefix string, store *redis.Client,
	opts ...PeriodLimitOption) *PeriodFailureLimit {
	if !strings.HasSuffix(keyPrefix, ":") {
		keyPrefix += ":"
	}
	limiter := &PeriodFailureLimit{
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

func (p *PeriodFailureLimit) align() { p.isAlign = true }

// CheckErr requests a permit state.
// same as Check
func (p *PeriodFailureLimit) CheckErr(ctx context.Context, key string, err error) (PeriodFailureLimitState, error) {
	return p.Check(ctx, key, err == nil)
}

// Check requests a permit.
func (p *PeriodFailureLimit) Check(ctx context.Context, key string, success bool) (PeriodFailureLimitState, error) {
	s := "0"
	if success {
		s = "1"
	}
	result, err := p.store.Eval(ctx,
		periodFailureLimitScript,
		[]string{p.keyPrefix + key},
		[]string{
			strconv.Itoa(p.quota),
			strconv.Itoa(p.calcExpireSeconds()),
			s,
		},
	).Result()
	if err != nil {
		return PeriodFailureLimitStsUnknown, err
	}
	code, ok := result.(int64)
	if !ok {
		return PeriodFailureLimitStsUnknown, ErrUnknownCode
	}
	switch code {
	case innerPeriodFailureLimitCodeSuccess:
		return PeriodFailureLimitStsSuccess, nil
	case innerPeriodFailureLimitCodeInQuota:
		return PeriodFailureLimitStsInQuota, nil
	case innerPeriodFailureLimitCodeOverQuota:
		return PeriodFailureLimitStsOverQuota, nil
	default:
		return PeriodFailureLimitStsUnknown, ErrUnknownCode
	}
}

// SetQuotaFull set a permit over quota.
func (p *PeriodFailureLimit) SetQuotaFull(ctx context.Context, key string) error {
	err := p.store.Eval(ctx,
		periodFailureLimitSetQuotaFullScript,
		[]string{p.keyPrefix + key},
		[]string{strconv.Itoa(p.quota)},
	).Err()
	if err == redis.Nil {
		return nil
	}
	return err
}

// Del delete a permit
func (p *PeriodFailureLimit) Del(ctx context.Context, key string) error {
	return p.store.Del(ctx, p.keyPrefix+key).Err()
}

// TTL get key ttl
// if key not exist, time = -1.
// if key exist, but not set expire time, t = -2
func (p *PeriodFailureLimit) TTL(ctx context.Context, key string) (time.Duration, error) {
	return p.store.TTL(ctx, p.keyPrefix+key).Result()
}

// GetInt get current failure count
func (p *PeriodFailureLimit) GetInt(ctx context.Context, key string) (int, bool, error) {
	v, err := p.store.Get(ctx, p.keyPrefix+key).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, false, nil
		}
		return 0, false, err
	}
	return v, true, nil
}

func (p *PeriodFailureLimit) calcExpireSeconds() int {
	if p.isAlign {
		now := time.Now()
		_, offset := now.Zone()
		unix := now.Unix() + int64(offset)
		return p.period - int(unix%int64(p.period))
	}
	return p.period
}
