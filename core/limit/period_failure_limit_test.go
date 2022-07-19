package limit

import (
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestPeriodFailureLimit_Check(t *testing.T) {
	testPeriodFailureLimit(t)
}

func TestPeriodFailureLimit_CheckWithAlign(t *testing.T) {
	testPeriodFailureLimit(t, Align())
}

func TestPeriodFailureLimit_RedisUnavailable(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)

	l := NewPeriodFailureLimit(seconds, quota, "PeriodFailureLimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))
	mr.Close()
	err = l.Check("first", true)
	assert.NotNil(t, err)
}

func testPeriodFailureLimit(t *testing.T, opts ...PeriodLimitOption) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)

	defer mr.Close()

	l := NewPeriodFailureLimit(seconds, quota, "PeriodFailureLimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}), opts...)
	var inLimitCnt, overFailureTimeCnt int
	for i := 0; i < total; i++ {
		err = l.Check("first", false)
		assert.Error(t, err)
		if errors.Is(err, ErrInLimitFailureTimes) {
			inLimitCnt++
		} else if errors.Is(err, ErrOverMaxFailureTimes) {
			overFailureTimeCnt++
		} else {
			t.Error("unknown status")
		}
	}
	assert.Equal(t, quota, inLimitCnt)
	assert.Equal(t, total-quota, overFailureTimeCnt)

	err = l.Check("first", true)
	assert.ErrorIs(t, err, ErrOverMaxFailureTimes)
}

func TestPeriodFailureLimit_Check_In_Limit_Failure_Time_Then_Success(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)

	defer mr.Close()

	l := NewPeriodFailureLimit(seconds, quota, "PeriodFailureLimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))
	var inLimitCnt, overFailureTimeCnt int
	for i := 0; i < quota-1; i++ {
		err = l.Check("first", false)
		assert.Error(t, err)
		if errors.Is(err, ErrInLimitFailureTimes) {
			inLimitCnt++
		} else if errors.Is(err, ErrOverMaxFailureTimes) {
			overFailureTimeCnt++
		} else {
			t.Error("unknown status")
		}
	}
	assert.Equal(t, quota-1, inLimitCnt)
	assert.Equal(t, 0, overFailureTimeCnt)

	err = l.Check("first", true)
	assert.NoError(t, err)

	v, existed, err := l.GetInt("first")
	assert.NoError(t, err)
	assert.False(t, existed)
	assert.Zero(t, v)
}

func TestPeriodFailureLimit_Check_Over_Limit_Failure_Time_Then_Success_Always_OverFailureTimeError(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)

	defer mr.Close()

	l := NewPeriodFailureLimit(seconds, quota, "PeriodFailureLimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))
	var inLimitCnt, overFailureTimeCnt int
	for i := 0; i < quota+1; i++ {
		err = l.Check("first", false)
		assert.Error(t, err)
		if errors.Is(err, ErrInLimitFailureTimes) {
			inLimitCnt++
		} else if errors.Is(err, ErrOverMaxFailureTimes) {
			overFailureTimeCnt++
		} else {
			t.Error("unknown status")
		}
	}
	assert.Equal(t, quota, inLimitCnt)
	assert.Equal(t, 1, overFailureTimeCnt)

	err = l.Check("first", true)
	assert.ErrorIs(t, err, ErrOverMaxFailureTimes)

	v, existed, err := l.GetInt("first")
	assert.NoError(t, err)
	assert.True(t, existed)
	assert.Equal(t, quota+1, v)
}

func TestPeriodFailureLimit_SetQuotaFull(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)
	defer mr.Close()

	l := NewPeriodFailureLimit(seconds, quota, "PeriodFailureLimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))

	err = l.SetQuotaFull("first")
	assert.Nil(t, err)

	err = l.Check("first", true)
	assert.ErrorIs(t, err, ErrOverMaxFailureTimes)
}

func TestPeriodFailureLimit_Del(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)
	defer mr.Close()

	l := NewPeriodFailureLimit(seconds, quota, "PeriodFailureLimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))

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

	err = l.Check("first", false)
	assert.ErrorIs(t, err, ErrOverMaxFailureTimes)

	err = l.Del("first")
	assert.Nil(t, err)

	err = l.Check("first", false)
	assert.ErrorIs(t, err, ErrInLimitFailureTimes)

	err = l.Check("first", true)
	assert.Nil(t, err)
}
