package verification

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	sendCodeScript = `
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
	verifyCodeScript = `
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

// error defined
var (
	ErrUnknownCode         = errors.New("verification: unknown status code")
	ErrMaxSendPerDay       = errors.New("verification: reach the maximum send times")
	ErrResendTooFrequently = errors.New("verification: resend too frequently")
	ErrCodeRequired        = errors.New("verification: code is required")
	ErrCodeExpired         = errors.New("verification: code is expired")
	ErrCodeMaxError        = errors.New("verification: reach the maximum error times")
	ErrCodeVerification    = errors.New("verification: code verification failed")
)

// Provider the provider
type Provider interface {
	Name() string
	SendCode(target, code string) error
}

// Verification verification code
type Verification struct {
	p                 Provider      // Provider send code
	store             *redis.Client // store client
	keyPrefix         string        // store 存验证码key的前缀, 默认 verification:
	keyExpires        time.Duration // store 存验证码key的过期时间, 默认: 24 小时
	maxErrorCount     int           // 最大验证失败次数, 默认: 3
	maxSendPerDay     int           // 验证码一天最大发送次数, 默认: 10
	availWindowSec    int           // 验证码有效窗口时间, 默认180, 单位: 秒
	resendIntervalSec int           // 重发验证码间隔时间, 默认60, 单位: 秒
}

// Option sms选项
type Option func(*Verification)

// WithKeyPrefix redis存验证码key的前缀, 默认 SMS:
func WithKeyPrefix(k string) Option {
	return func(v *Verification) {
		if k != "" {
			if !strings.HasSuffix(k, ":") {
				k += ":"
			}
			v.keyPrefix = k
		}
	}
}

// WithKeyExpires redis存验证码key的过期时间, 默认 24小时
func WithKeyExpires(expires time.Duration) Option {
	return func(v *Verification) {
		v.keyExpires = expires
	}
}

// WithMaxErrorCount 验证码最大验证失败次数, 默认: 3
func WithMaxErrorCount(cnt int) Option {
	return func(v *Verification) {
		v.maxErrorCount = cnt
	}
}

// WithMaxSendPerDay 验证码一天最大发送次数, 默认: 10
func WithMaxSendPerDay(cnt int) Option {
	return func(v *Verification) {
		v.maxSendPerDay = cnt
	}
}

// WithAvailWindowSecond 验证码有效窗口时间, 默认180, 单位: 秒
func WithAvailWindowSecond(sec int) Option {
	return func(v *Verification) {
		v.availWindowSec = sec
	}
}

// WithResendIntervalSecond 重发验证码间隔时间, 默认60, 单位: 秒
func WithResendIntervalSecond(sec int) Option {
	return func(v *Verification) {
		v.resendIntervalSec = sec
	}
}

// New new a Verification
func New(p Provider, store *redis.Client, opts ...Option) *Verification {
	v := &Verification{
		p,
		store,
		"verification:",
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
func (sf *Verification) Name() string { return sf.p.Name() }

// SendCode send code and store in redis cache.
func (sf *Verification) SendCode(target, code string) error {
	result, err := sf.store.Eval(context.Background(), sendCodeScript,
		[]string{sf.keyPrefix + target},
		[]string{
			code,
			strconv.Itoa(sf.maxSendPerDay),
			strconv.Itoa(sf.resendIntervalSec),
			strconv.FormatInt(time.Now().Unix(), 10),
			strconv.FormatInt(int64(sf.keyExpires/time.Second), 10),
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
		return sf.p.SendCode(target, code)
	case 1:
		return ErrMaxSendPerDay
	case 2:
		return ErrResendTooFrequently
	default:
		return ErrUnknownCode
	}
}

// VerifyCode verify code from redis cache.
func (sf *Verification) VerifyCode(target, code string) error {
	result, err := sf.store.Eval(context.Background(), verifyCodeScript,
		[]string{sf.keyPrefix + target},
		[]string{
			code,
			strconv.Itoa(sf.maxErrorCount),
			strconv.Itoa(sf.availWindowSec),
			strconv.FormatInt(time.Now().Unix(), 10),
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
		return nil
	case 1:
		return ErrCodeRequired
	case 2:
		return ErrCodeExpired
	case 3:
		return ErrCodeMaxError
	case 4:
		return ErrCodeVerification
	default:
		return ErrUnknownCode
	}
}
