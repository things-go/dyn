package limit

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	verifiedLimitSendCodeScript = `
local key = KEYS[1] -- key
local code = ARGV[1] -- 验证码
local maxSendPerDay = tonumber(ARGV[2]) -- 验证码一天最大发送次数
local resendWindowTime = tonumber(ARGV[3]) -- 验证码重发限制窗口时间
local now = tonumber(ARGV[4]) -- 当前时间, 单位秒
local expires = tonumber(ARGV[5]) -- key 过期时间, 单位秒

if (redis.call('exists', key) == 1) then
    local sendCnt = tonumber(redis.call('HGET', key, "send"))
    if sendCnt >= maxSendPerDay then
        return 1 -- 超过每天发送限制次数
    end

    local lastedAt = tonumber(redis.call('HGET', key, "lasted"))
    if lastedAt + resendWindowTime > now then
        return 2 -- 发送过于频繁, 即还在重发限制窗口
    end
    sendCnt = sendCnt + 1

	redis.call('HMSET', key, 'code', code, 'err', 0, 'send', sendCnt, 'lasted', now)
else
    redis.call('HMSET', key, 'code', code, 'err', 0, 'send', 1, 'lasted', now)
	redis.call("expire", key, expires)
end

return 0 
`
	verifiedLimitVerifyCodeScript = `
local key = KEYS[1] -- key
local code = ARGV[1] -- 验证码
local maxErrorCount = tonumber(ARGV[2]) -- 验证码最大验证失败次数
local availWindowTime = tonumber(ARGV[3]) -- 验证码有效窗口时间, 单位秒
local now = tonumber(ARGV[4]) -- 当前时间, 单位秒

if redis.call('exists', key) == 0 then
    return 1  -- 未发送短信验证码
end

local lastedAt = tonumber(redis.call('HGET', key, "lasted"))
local errCnt = tonumber(redis.call('HGET', key, "err"))
local currentCode = redis.call('HGET', key, "code")
if lastedAt + availWindowTime < now then
    return 2  -- 验证码已过期
end
if errCnt >= maxErrorCount then
    return 3  -- 验证码错误次数超过限制
end
if currentCode == code then
    redis.call('HSET', key, "lasted", lastedAt - availWindowTime) -- 设置成过期
    return 0
else
    redis.call('HINCRBY', key, "err", 1)
    return 4 -- 验证码错误
end
`
)

// error defined for verified
var (
	ErrMaxSendPerDay       = errors.New("limit: reach the maximum send times")
	ErrResendTooFrequently = errors.New("limit: resend too frequently")
	ErrCodeRequired        = errors.New("limit: code is required")
	ErrCodeExpired         = errors.New("limit: code is expired")
	ErrCodeMaxError        = errors.New("limit: reach the maximum error times")
	ErrCodeVerification    = errors.New("limit: code verified failed")
)

// VerifiedProvider the provider
type VerifiedProvider interface {
	Name() string
	SendCode(target, code string) error
}

// VerifiedLimit verified code limit
type VerifiedLimit struct {
	p                 VerifiedProvider // VerifiedProvider send code
	store             *redis.Client    // store client
	keyPrefix         string           // store 存验证码key的前缀, 默认 verified:
	keyExpires        time.Duration    // store 存验证码key的过期时间, 默认: 24 小时
	maxErrorCount     int              // 最大验证失败次数, 默认: 3
	maxSendPerDay     int              // 验证码一天最大发送次数, 默认: 10
	availWindowSec    int              // 验证码有效窗口时间, 默认180, 单位: 秒
	resendIntervalSec int              // 重发验证码间隔时间, 默认60, 单位: 秒
}

// Option sms选项
type Option func(*VerifiedLimit)

// WithVerifiedKeyPrefix redis存验证码key的前缀, 默认 SMS:
func WithVerifiedKeyPrefix(k string) Option {
	return func(v *VerifiedLimit) {
		if k != "" {
			if !strings.HasSuffix(k, ":") {
				k += ":"
			}
			v.keyPrefix = k
		}
	}
}

// WithVerifiedKeyExpires redis存验证码key的过期时间, 默认 24小时
func WithVerifiedKeyExpires(expires time.Duration) Option {
	return func(v *VerifiedLimit) {
		v.keyExpires = expires
	}
}

// WithVerifiedMaxErrorCount 验证码最大验证失败次数, 默认: 3
func WithVerifiedMaxErrorCount(cnt int) Option {
	return func(v *VerifiedLimit) {
		v.maxErrorCount = cnt
	}
}

// WithVerifiedMaxSendPerDay 验证码一天最大发送次数, 默认: 10
func WithVerifiedMaxSendPerDay(cnt int) Option {
	return func(v *VerifiedLimit) {
		v.maxSendPerDay = cnt
	}
}

// WithVerifiedAvailWindowSecond 验证码有效窗口时间, 默认180, 单位: 秒
func WithVerifiedAvailWindowSecond(sec int) Option {
	return func(v *VerifiedLimit) {
		v.availWindowSec = sec
	}
}

// WithVerifiedResendIntervalSecond 重发验证码间隔时间, 默认60, 单位: 秒
func WithVerifiedResendIntervalSecond(sec int) Option {
	return func(v *VerifiedLimit) {
		v.resendIntervalSec = sec
	}
}

// NewVerified  new a verified limit
func NewVerified(p VerifiedProvider, store *redis.Client, opts ...Option) *VerifiedLimit {
	v := &VerifiedLimit{
		p,
		store,
		"verified:",
		time.Hour * 24,
		3,
		10,
		180,
		60,
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// Name the provider name
func (v *VerifiedLimit) Name() string { return v.p.Name() }

// SendCode send code and store in redis cache.
func (v *VerifiedLimit) SendCode(target, code string) error {
	result, err := v.store.Eval(context.Background(), verifiedLimitSendCodeScript,
		[]string{v.keyPrefix + target},
		[]string{
			code,
			strconv.Itoa(v.maxSendPerDay),
			strconv.Itoa(v.resendIntervalSec),
			strconv.FormatInt(time.Now().Unix(), 10),
			strconv.FormatInt(int64(v.keyExpires/time.Second), 10),
		},
	).Result()
	if err != nil {
		return err
	}
	sts, ok := result.(int64)
	if !ok {
		return ErrUnknownCode
	}
	switch sts {
	case 0:
		err = v.p.SendCode(target, code)
	case 1:
		err = ErrMaxSendPerDay
	case 2:
		err = ErrResendTooFrequently
	default:
		err = ErrUnknownCode
	}
	return err
}

// VerifyCode verify code from redis cache.
func (v *VerifiedLimit) VerifyCode(target, code string) error {
	result, err := v.store.Eval(context.Background(), verifiedLimitVerifyCodeScript,
		[]string{v.keyPrefix + target},
		[]string{
			code,
			strconv.Itoa(v.maxErrorCount),
			strconv.Itoa(v.availWindowSec),
			strconv.FormatInt(time.Now().Unix(), 10),
		},
	).Result()
	if err != nil {
		return err
	}

	sts, ok := result.(int64)
	if !ok {
		err = ErrUnknownCode
	}
	switch sts {
	case 0:
		return nil
	case 1:
		err = ErrCodeRequired
	case 2:
		err = ErrCodeExpired
	case 3:
		err = ErrCodeMaxError
	case 4:
		err = ErrCodeVerification
	default:
		err = ErrUnknownCode
	}
	return err
}
