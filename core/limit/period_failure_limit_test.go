package limit

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var internalErr = errors.New("internal error")

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
	err = l.CheckErr(context.Background(), "first", nil)
	assert.NotNil(t, err)
}

func testPeriodFailureLimit(t *testing.T, opts ...PeriodLimitOption) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)

	defer mr.Close()

	l := NewPeriodFailureLimit(seconds, quota, "PeriodFailureLimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}), opts...)
	var inLimitCnt, overFailureTimeCnt int
	for i := 0; i < total; i++ {
		err = l.CheckErr(context.Background(), "first", internalErr)
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

	err = l.CheckErr(context.Background(), "first", nil)
	assert.ErrorIs(t, err, ErrOverMaxFailureTimes)
}

func TestPeriodFailureLimit_Check_In_Limit_Failure_Time_Then_Success(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)

	defer mr.Close()

	l := NewPeriodFailureLimit(seconds, quota, "PeriodFailureLimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))
	var inLimitCnt, overFailureTimeCnt int
	for i := 0; i < quota-1; i++ {
		err = l.CheckErr(context.Background(), "first", internalErr)
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

	err = l.CheckErr(context.Background(), "first", nil)
	assert.NoError(t, err)

	v, existed, err := l.GetInt(context.Background(), "first")
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
		err = l.CheckErr(context.Background(), "first", internalErr)
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

	err = l.CheckErr(context.Background(), "first", nil)
	assert.ErrorIs(t, err, ErrOverMaxFailureTimes)

	v, existed, err := l.GetInt(context.Background(), "first")
	assert.NoError(t, err)
	assert.True(t, existed)
	assert.Equal(t, quota+1, v)
}

func TestPeriodFailureLimit_SetQuotaFull(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)
	defer mr.Close()

	l := NewPeriodFailureLimit(seconds, quota, "PeriodFailureLimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))

	err = l.SetQuotaFull(context.Background(), "first")
	assert.Nil(t, err)

	err = l.CheckErr(context.Background(), "first", nil)
	assert.ErrorIs(t, err, ErrOverMaxFailureTimes)
}

func TestPeriodFailureLimit_Del(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)
	defer mr.Close()

	l := NewPeriodFailureLimit(seconds, quota, "PeriodFailureLimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))

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

	err = l.CheckErr(context.Background(), "first", internalErr)
	assert.ErrorIs(t, err, ErrOverMaxFailureTimes)

	err = l.Del(context.Background(), "first")
	assert.Nil(t, err)

	err = l.CheckErr(context.Background(), "first", internalErr)
	assert.ErrorIs(t, err, ErrInLimitFailureTimes)

	err = l.CheckErr(context.Background(), "first", nil)
	assert.Nil(t, err)
}
