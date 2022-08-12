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
	sts, err := l.CheckErr(context.Background(), "first", nil)
	assert.Error(t, err)
	assert.Equal(t, PeriodFailureLimitStsUnknown, sts)
}

func testPeriodFailureLimit(t *testing.T, opts ...PeriodLimitOption) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)

	defer mr.Close()

	l := NewPeriodFailureLimit(seconds, quota, "PeriodFailureLimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}), opts...)
	var inLimitCnt, overFailureTimeCnt int
	for i := 0; i < total; i++ {
		sts, err := l.CheckErr(context.Background(), "first", internalErr)
		assert.NoError(t, err)
		switch sts {
		case PeriodFailureLimitStsInQuota:
			inLimitCnt++
		case PeriodFailureLimitStsOverQuota:
			overFailureTimeCnt++
		default:
			t.Errorf("unknown status, must be on of [%d, %d]", PeriodFailureLimitStsInQuota, PeriodFailureLimitStsOverQuota)
		}
	}
	assert.Equal(t, quota, inLimitCnt)
	assert.Equal(t, total-quota, overFailureTimeCnt)

	sts, err := l.CheckErr(context.Background(), "first", nil)
	assert.NoError(t, err)
	assert.Equal(t, PeriodFailureLimitStsOverQuota, sts)
}

func TestPeriodFailureLimit_Check_In_Limit_Failure_Time_Then_Success(t *testing.T) {
	mr, err := miniredis.Run()
	assert.Nil(t, err)

	defer mr.Close()

	l := NewPeriodFailureLimit(seconds, quota, "PeriodFailureLimit", redis.NewClient(&redis.Options{Addr: mr.Addr()}))
	var inLimitCnt, overFailureTimeCnt int
	for i := 0; i < quota-1; i++ {
		sts, err := l.CheckErr(context.Background(), "first", internalErr)
		assert.NoError(t, err)
		switch sts {
		case PeriodFailureLimitStsInQuota:
			inLimitCnt++
		case PeriodFailureLimitStsOverQuota:
			overFailureTimeCnt++
		default:
			t.Errorf("unknown status, must be on of [%d, %d]", PeriodFailureLimitStsInQuota, PeriodFailureLimitStsOverQuota)
		}
	}
	assert.Equal(t, quota-1, inLimitCnt)
	assert.Equal(t, 0, overFailureTimeCnt)

	sts, err := l.CheckErr(context.Background(), "first", nil)
	assert.NoError(t, err)
	assert.Equal(t, PeriodFailureLimitStsSuccess, sts)

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
		sts, err := l.CheckErr(context.Background(), "first", internalErr)
		assert.NoError(t, err)
		switch sts {
		case PeriodFailureLimitStsInQuota:
			inLimitCnt++
		case PeriodFailureLimitStsOverQuota:
			overFailureTimeCnt++
		default:
			t.Errorf("unknown status, must be on of [%d, %d]", PeriodFailureLimitStsInQuota, PeriodFailureLimitStsOverQuota)
		}
	}
	assert.Equal(t, quota, inLimitCnt)
	assert.Equal(t, 1, overFailureTimeCnt)

	sts, err := l.CheckErr(context.Background(), "first", nil)
	assert.NoError(t, err)
	assert.Equal(t, PeriodFailureLimitStsOverQuota, sts)

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

	sts, err := l.CheckErr(context.Background(), "first", nil)
	assert.NoError(t, err)
	assert.Equal(t, PeriodFailureLimitStsOverQuota, sts)
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

	// 第一次ttl, 不存在
	tt, err := l.TTL(context.Background(), "first")
	assert.Nil(t, err)
	assert.Equal(t, int(tt), -2)

	err = l.SetQuotaFull(context.Background(), "first")
	assert.Nil(t, err)

	// 第二次ttl, key 存在
	tt, err = l.TTL(context.Background(), "first")
	assert.Nil(t, err)
	assert.LessOrEqual(t, int(tt.Seconds()), seconds)

	v, b, err = l.GetInt(context.Background(), "first")
	assert.Nil(t, err)
	assert.True(t, b)
	assert.Equal(t, quota, v)

	sts, err := l.CheckErr(context.Background(), "first", internalErr)
	assert.NoError(t, err)
	assert.Equal(t, PeriodFailureLimitStsOverQuota, sts)

	err = l.Del(context.Background(), "first")
	assert.Nil(t, err)

	sts, err = l.CheckErr(context.Background(), "first", internalErr)
	assert.NoError(t, err)
	assert.Equal(t, PeriodFailureLimitStsInQuota, sts)

	sts, err = l.CheckErr(context.Background(), "first", nil)
	assert.NoError(t, err)
	assert.Equal(t, PeriodFailureLimitStsSuccess, sts)
}
