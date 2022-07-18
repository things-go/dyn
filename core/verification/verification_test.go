package verification

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	target  = "112233"
	code    = "123456"
	badCode = "654321"
)

var _ Provider = (*TestProvider)(nil)

type TestProvider struct{}

func (t TestProvider) Name() string { return "test_provider" }

func (t TestProvider) SendCode(target, code string) error { return nil }

func TestName(t *testing.T) {
	l := New(new(TestProvider), redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"}))
	require.Equal(t, "test_provider", l.Name())
}

func TestSendCode_RedisUnavailable(t *testing.T) {
	mr, err := miniredis.Run()
	require.Nil(t, err)

	l := New(new(TestProvider), redis.NewClient(&redis.Options{Addr: mr.Addr()}))
	mr.Close()

	err = l.SendCode(target, code)
	assert.NotNil(t, err)
}
func TestSendCode_Success(t *testing.T) {
	mr, err := miniredis.Run()
	require.Nil(t, err)
	defer mr.Close()

	l := New(new(TestProvider),
		redis.NewClient(&redis.Options{Addr: mr.Addr()}),
		WithKeyPrefix("verification"),
		WithKeyExpires(time.Hour),
	)
	err = l.SendCode(target, code)
	require.NoError(t, err)
}

func TestSendCode_MaxSendPerDay(t *testing.T) {
	mr, err := miniredis.Run()
	require.Nil(t, err)
	defer mr.Close()

	l := New(new(TestProvider),
		redis.NewClient(&redis.Options{Addr: mr.Addr()}),
		WithMaxSendPerDay(1),
		WithResendIntervalSecond(1),
	)

	err = l.SendCode(target, code)
	require.NoError(t, err)

	time.Sleep(time.Second + time.Millisecond*10)
	err = l.SendCode(target, code)
	require.ErrorIs(t, err, ErrMaxSendPerDay)
}

func TestSendCode_Concurrency_MaxSendPerDay(t *testing.T) {
	var success uint32
	var failed uint32

	mr, err := miniredis.Run()
	require.Nil(t, err)
	defer mr.Close()

	l := New(new(TestProvider),
		redis.NewClient(&redis.Options{Addr: mr.Addr()}),
		WithMaxSendPerDay(1),
	)

	wg := &sync.WaitGroup{}
	wg.Add(15)
	for i := 0; i < 15; i++ {
		go func() {
			defer wg.Done()

			err := l.SendCode(target, code)
			if err != nil {
				require.ErrorIs(t, err, ErrMaxSendPerDay)
				atomic.AddUint32(&failed, 1)
			} else {
				atomic.AddUint32(&success, 1)
			}
		}()
	}

	wg.Wait()
	require.Equal(t, uint32(1), success)
	require.Equal(t, uint32(14), failed)
}

func TestSendCode_ResendTooFrequently(t *testing.T) {
	mr, err := miniredis.Run()
	require.Nil(t, err)
	defer mr.Close()

	l := New(new(TestProvider),
		redis.NewClient(&redis.Options{Addr: mr.Addr()}),
		WithResendIntervalSecond(1),
	)

	err = l.SendCode(target, code)
	require.NoError(t, err)
	err = l.SendCode(target, code)
	require.ErrorIs(t, err, ErrResendTooFrequently)
}

func TestSendCode_Concurrency_ResendTooFrequently(t *testing.T) {
	var success uint32
	var failed uint32

	mr, err := miniredis.Run()
	require.Nil(t, err)
	defer mr.Close()

	l := New(new(TestProvider),
		redis.NewClient(&redis.Options{Addr: mr.Addr()}),
		WithResendIntervalSecond(1),
	)

	wg := &sync.WaitGroup{}
	wg.Add(15)
	for i := 0; i < 15; i++ {
		go func() {
			defer wg.Done()

			err := l.SendCode(target, code)
			if err != nil {
				require.ErrorIs(t, err, ErrResendTooFrequently)
				atomic.AddUint32(&failed, 1)
			} else {
				atomic.AddUint32(&success, 1)
			}
		}()
	}

	wg.Wait()
	require.Equal(t, uint32(1), success)
	require.Equal(t, uint32(14), failed)
}

func TestVerifyCode_Success(t *testing.T) {
	mr, err := miniredis.Run()
	require.Nil(t, err)
	defer mr.Close()

	l := New(new(TestProvider),
		redis.NewClient(&redis.Options{Addr: mr.Addr()}),
		WithResendIntervalSecond(1),
	)

	err = l.SendCode(target, code)
	require.Nil(t, err)

	err = l.VerifyCode(target, code)
	assert.NoError(t, err)
}

func TestVerifyCode_CodeRequired(t *testing.T) {
	mr, err := miniredis.Run()
	require.Nil(t, err)
	defer mr.Close()

	l := New(new(TestProvider), redis.NewClient(&redis.Options{Addr: mr.Addr()}))

	err = l.VerifyCode(target, code)
	assert.Error(t, err, ErrCodeRequired)
}

func TestVerifyCode_CodeExpired(t *testing.T) {
	mr, err := miniredis.Run()
	require.Nil(t, err)
	defer mr.Close()

	l := New(new(TestProvider),
		redis.NewClient(&redis.Options{Addr: mr.Addr()}),
		WithAvailWindowSecond(1),
	)
	err = l.SendCode(target, code)
	require.Nil(t, err)

	time.Sleep(time.Second * 2)
	err = l.VerifyCode(target, code)
	assert.Error(t, err, ErrCodeExpired)
}
func TestVerifyCode_CodeMaxError(t *testing.T) {
	mr, err := miniredis.Run()
	require.Nil(t, err)
	defer mr.Close()

	l := New(new(TestProvider),
		redis.NewClient(&redis.Options{Addr: mr.Addr()}),
		WithMaxErrorCount(3),
	)
	err = l.SendCode(target, code)
	require.Nil(t, err)

	for i := 0; i < 3; i++ {
		err = l.VerifyCode(target, badCode)
		assert.Error(t, err, ErrCodeVerification)
	}
	err = l.VerifyCode(target, badCode)
	assert.Error(t, err, ErrCodeMaxError)
}

func TestVerifyCode_Concurrency_CodeMaxError(t *testing.T) {
	var failedMaxError uint32
	var failedVerify uint32

	mr, err := miniredis.Run()
	require.Nil(t, err)
	defer mr.Close()

	l := New(new(TestProvider),
		redis.NewClient(&redis.Options{Addr: mr.Addr()}),
		WithMaxErrorCount(3),
	)

	err = l.SendCode(target, code)
	require.Nil(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(15)
	for i := 0; i < 15; i++ {
		go func() {
			defer wg.Done()

			err := l.VerifyCode(target, badCode)
			if err != nil {
				if errors.Is(err, ErrCodeMaxError) {
					atomic.AddUint32(&failedMaxError, 1)
				} else {
					atomic.AddUint32(&failedVerify, 1)
				}
			}
		}()
	}

	wg.Wait()
	require.Equal(t, uint32(3), failedVerify)
	require.Equal(t, uint32(12), failedMaxError)
}

func TestErrorMap(t *testing.T) {
	var ErrCodeVerificationExample = errors.New("example verification error")
	var ErrCodeMaxErrorExample = errors.New("example max error")

	mr, err := miniredis.Run()
	require.Nil(t, err)
	defer mr.Close()

	l := New(new(TestProvider),
		redis.NewClient(&redis.Options{Addr: mr.Addr()}),
		WithErrorMap(map[error]error{
			ErrCodeVerification: ErrCodeVerificationExample,
			ErrCodeMaxError:     ErrCodeMaxErrorExample,
		}),
	)
	err = l.SendCode(target, code)
	require.Nil(t, err)

	for i := 0; i < 3; i++ {
		err = l.VerifyCode(target, badCode)
		assert.Error(t, err, ErrCodeVerificationExample)
	}
	err = l.VerifyCode(target, badCode)
	assert.Error(t, err, ErrCodeMaxErrorExample)
}
