package limit

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

const (
	seconds = 1
	quota   = 5
	total   = 100
)

func TestPeriodLimit_Take(t *testing.T) {
	testPeriodLimit(t)
}

func TestPeriodLimit_TakeWithAlign(t *testing.T) {
	testPeriodLimit(t, Align())
}

func TestPeriodLimit_RedisUnavailable(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)

	l := NewPeriodLimit(seconds, quota, "periodlimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))
	mr.Close()
	val, err := l.Take("first")
	assert.NotNil(t, err)
	assert.Equal(t, Unknown, val)
}

func testPeriodLimit(t *testing.T, opts ...PeriodLimitOption) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)

	defer mr.Close()

	l := NewPeriodLimit(seconds, quota, "periodlimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}), opts...)
	var allowed, hitQuota, overQuota int
	for i := 0; i < total; i++ {
		val, err := l.Take("first")
		if err != nil {
			t.Error(err)
		}
		switch val {
		case Allowed:
			allowed++
		case HitQuota:
			hitQuota++
		case OverQuota:
			overQuota++
		case Unknown:
			fallthrough
		default:
			t.Error("unknown status")
		}
	}

	assert.Equal(t, quota-1, allowed)
	assert.Equal(t, 1, hitQuota)
	assert.Equal(t, total-quota, overQuota)
}

func TestQuotaFull(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)
	defer mr.Close()

	l := NewPeriodLimit(1, 1, "periodlimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))
	val, err := l.Take("first")
	assert.Nil(t, err)
	assert.True(t, val.IsHitQuota())
}

func TestSetQuotaFull(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)
	defer mr.Close()

	l := NewPeriodLimit(seconds, quota, "periodlimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))

	err = l.SetQuotaFull("first")
	assert.Nil(t, err)

	val, err := l.Take("first")
	assert.Nil(t, err)
	assert.Equal(t, OverQuota, val)
}

func TestDel(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)
	defer mr.Close()

	l := NewPeriodLimit(seconds, quota, "periodlimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))

	v, b, err := l.GetInt("first")
	assert.Nil(t, err)
	assert.False(t, b)
	assert.Equal(t, 0, v)

	// 第一次ttl, 不存在
	tt, err := l.TTL("first")
	assert.Nil(t, err)
	assert.Equal(t, int(tt), -2)

	err = l.SetQuotaFull("first")
	assert.Nil(t, err)

	// 第二次ttl, key 存在
	tt, err = l.TTL("first")
	assert.Nil(t, err)
	assert.LessOrEqual(t, int(tt.Seconds()), seconds)

	v, b, err = l.GetInt("first")
	assert.Nil(t, err)
	assert.True(t, b)
	assert.Equal(t, quota, v)

	val, err := l.Take("first")
	assert.Nil(t, err)
	assert.True(t, val.IsOverQuota())

	err = l.Del("first")
	assert.Nil(t, err)

	val, err = l.Take("first")
	assert.Nil(t, err)
	assert.True(t, val.IsAllowed())
}
