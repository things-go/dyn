package limit

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestPeriodLimit_Take(t *testing.T) {
	testPeriodLimit(t)
}

func TestPeriodLimit_TakeWithAlign(t *testing.T) {
	testPeriodLimit(t, Align())
}

func TestPeriodLimit_RedisUnavailable(t *testing.T) {
	const (
		seconds = 1
		quota   = 5
	)

	s, err := miniredis.Run()
	assert.Nil(t, err)

	l := NewPeriodLimit(seconds, quota, "periodlimit", redis.NewClient(&redis.Options{Addr: s.Addr()}))
	s.Close()
	val, err := l.Take("first")
	assert.NotNil(t, err)
	assert.Equal(t, Unknown, val)
}

func testPeriodLimit(t *testing.T, opts ...PeriodLimitOption) {
	const (
		seconds = 1
		total   = 100
		quota   = 5
	)

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
		default:
			t.Error("unknown status")
		}
	}

	assert.Equal(t, quota-1, allowed)
	assert.Equal(t, 1, hitQuota)
	assert.Equal(t, total-quota, overQuota)
}

func TestQuotaFull(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)

	l := NewPeriodLimit(1, 1, "periodlimit", redis.NewClient(&redis.Options{Addr: s.Addr()}))
	val, err := l.Take("first")
	assert.Nil(t, err)
	assert.Equal(t, HitQuota, val)
}
