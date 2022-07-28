package limit

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ErrInLimitFailureTimes = errors.New("limit: in limit failure times")
	ErrOverMaxFailureTimes = errors.New("limit: over the max failure times")
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

// NewPeriodLimit returns a PeriodLimit with given parameters.
func NewPeriodFailureLimit(periodSecond, quota int, keyPrefix string,
	store *redis.Client, opts ...PeriodLimitOption) *PeriodFailureLimit {
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

func (p *PeriodFailureLimit) align() {
	p.isAlign = true
}

// CheckErr requests a permit.
// return result:
// nil: 表示成功
// ErrUnknownCode: lua脚本错误
// ErrInLimitFailureTimes: 表示还在最大失败次数范围内
// ErrOverMaxFailureTimes: 表示超过了最大失败验证次数
// NOTE: success 为 false, 只会出现 ErrInLimitFailureTimes 或 ErrOverMaxFailureTimes
func (p *PeriodFailureLimit) CheckErr(ctx context.Context, key string, err error) error {
	return p.Check(ctx, key, err == nil)
}

// Check requests a permit.
// return result:
// nil: 表示成功
// ErrUnknownCode: lua脚本错误
// ErrInLimitFailureTimes: 表示还在最大失败次数范围内
// ErrOverMaxFailureTimes: 表示超过了最大失败验证次数
// NOTE: success 为 false, 只会出现 ErrInLimitFailureTimes 或 ErrOverMaxFailureTimes
func (p *PeriodFailureLimit) Check(ctx context.Context, key string, success bool) error {
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
		return err
	}
	code, ok := result.(int64)
	if !ok {
		return ErrUnknownCode
	}
	switch code {
	case 0:
		return nil
	case 1:
		return ErrInLimitFailureTimes
	case 2:
		return ErrOverMaxFailureTimes
	default:
		return ErrUnknownCode
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

// GetInt get count
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
