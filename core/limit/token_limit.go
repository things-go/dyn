package limit

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	xrate "golang.org/x/time/rate"

	"github.com/things-go/dyn/log"
)

// KEYS[1] as tokens_key
// KEYS[2] as timestamp_key
const (
	tokenLimitScript = `local rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])
local fill_time = capacity/rate
local ttl = math.floor(fill_time*2)
local last_tokens = tonumber(redis.call("GET", KEYS[1]))
if last_tokens == nil then
    last_tokens = capacity
end

local last_refreshed = tonumber(redis.call("GET", KEYS[2]))
if last_refreshed == nil then
    last_refreshed = 0
end

local delta = math.max(0, now-last_refreshed)
local filled_tokens = math.min(capacity, last_tokens+(delta*rate))
local allowed = filled_tokens >= requested
local new_tokens = filled_tokens
if allowed then
    new_tokens = filled_tokens - requested
end

redis.call("SETEX", KEYS[1], ttl, new_tokens)
redis.call("SETEX", KEYS[2], ttl, now)

return allowed`
	tokenFormat     = "{%s}.tokens"
	timestampFormat = "{%s}.ts"
	pingInterval    = time.Millisecond * 100
)

// TokenLimit controls how frequently events are allowed to happen with in one second.
type TokenLimit struct {
	rate           int
	burst          int
	store          *redis.Client
	tokenKey       string
	timestampKey   string
	rescueLock     sync.Mutex
	isRedisAlive   uint32
	rescueLimiter  *xrate.Limiter
	monitorStarted bool
}

// NewTokenLimit returns a new TokenLimit that allows events up to rate and permits
// bursts of at most burst tokens.
func NewTokenLimit(rate, burst int, key string, store *redis.Client) *TokenLimit {
	return &TokenLimit{
		rate:          rate,
		burst:         burst,
		store:         store,
		tokenKey:      fmt.Sprintf(tokenFormat, key),
		timestampKey:  fmt.Sprintf(timestampFormat, key),
		isRedisAlive:  1,
		rescueLimiter: xrate.NewLimiter(xrate.Every(time.Second/time.Duration(rate)), burst),
	}
}

// Allow is shorthand for AllowN(time.Now(), 1).
func (t *TokenLimit) Allow() bool {
	return t.AllowN(time.Now(), 1)
}

// AllowN reports whether n events may happen at time now.
// Use this method if you intend to drop / skip events that exceed the rate.
// Otherwise, use Reserve or Wait.
func (t *TokenLimit) AllowN(now time.Time, n int) bool {
	return t.reserveN(now, n)
}

func (t *TokenLimit) reserveN(now time.Time, n int) bool {
	if atomic.LoadUint32(&t.isRedisAlive) == 0 {
		return t.rescueLimiter.AllowN(now, n)
	}

	resp, err := t.store.Eval(context.Background(), tokenLimitScript,
		[]string{t.tokenKey, t.timestampKey},
		[]string{
			strconv.Itoa(t.rate),
			strconv.Itoa(t.burst),
			strconv.FormatInt(now.Unix(), 10),
			strconv.Itoa(n),
		}).Result()
	// redis allowed == false
	// Lua boolean false -> r Nil bulk reply
	if err == redis.Nil {
		return false
	}
	if err != nil {
		log.Errorf("fail to use rate limiter: %s, use in-process limiter for rescue", err)
		t.startMonitor()
		return t.rescueLimiter.AllowN(now, n)
	}

	code, ok := resp.(int64)
	if !ok {
		log.Errorf("fail to eval redis script: %v, use in-process limiter for rescue", resp)
		t.startMonitor()
		return t.rescueLimiter.AllowN(now, n)
	}

	// redis allowed == true
	// Lua boolean true -> r integer reply with value of 1
	return code == 1
}

func (t *TokenLimit) startMonitor() {
	t.rescueLock.Lock()
	defer t.rescueLock.Unlock()

	if t.monitorStarted {
		return
	}

	t.monitorStarted = true
	atomic.StoreUint32(&t.isRedisAlive, 0)

	go t.waitForRedis()
}

func (t *TokenLimit) waitForRedis() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		t.rescueLock.Lock()
		t.monitorStarted = false
		t.rescueLock.Unlock()
	}()

	for range ticker.C {
		v, err := t.store.Ping(context.Background()).Result()
		if err != nil {
			continue
		}
		if v == "PONG" {
			atomic.StoreUint32(&t.isRedisAlive, 1)
			return
		}
	}
}
