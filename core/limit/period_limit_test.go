package limit

import (
	"context"
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
	val, err := l.Take(context.Background(), "first")
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
		val, err := l.Take(context.Background(), "first")
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

func TestPeriodLimit_QuotaFull(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)
	defer mr.Close()

	l := NewPeriodLimit(1, 1, "periodlimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))
	val, err := l.Take(context.Background(), "first")
	assert.Nil(t, err)
	assert.True(t, val.IsHitQuota())
}

func TestPeriodLimit_SetQuotaFull(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)
	defer mr.Close()

	l := NewPeriodLimit(seconds, quota, "periodlimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))

	err = l.SetQuotaFull(context.Background(), "first")
	assert.Nil(t, err)

	val, err := l.Take(context.Background(), "first")
	assert.Nil(t, err)
	assert.Equal(t, OverQuota, val)
}

func TestPeriodLimit_Del(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)
	defer mr.Close()

	l := NewPeriodLimit(seconds, quota, "periodlimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))

	v, b, err := l.GetInt(context.Background(), "first")
	assert.Nil(t, err)
	assert.False(t, b)
	assert.Equal(t, 0, v)

	// ?????????ttl, ?????????
	tt, err := l.TTL(context.Background(), "first")
	assert.Nil(t, err)
	assert.Equal(t, int(tt), -2)

	err = l.SetQuotaFull(context.Background(), "first")
	assert.Nil(t, err)

	// ?????????ttl, key ??????
	tt, err = l.TTL(context.Background(), "first")
	assert.Nil(t, err)
	assert.LessOrEqual(t, int(tt.Seconds()), seconds)

	v, b, err = l.GetInt(context.Background(), "first")
	assert.Nil(t, err)
	assert.True(t, b)
	assert.Equal(t, quota, v)

	val, err := l.Take(context.Background(), "first")
	assert.Nil(t, err)
	assert.True(t, val.IsOverQuota())

	err = l.Del(context.Background(), "first")
	assert.Nil(t, err)

	val, err = l.Take(context.Background(), "first")
	assert.Nil(t, err)
	assert.True(t, val.IsAllowed())
}
