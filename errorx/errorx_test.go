package errorx_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/things-go/dyn/errorx"
)

func Test_NilError(t *testing.T) {
	var err *errorx.Error

	err = err.TakeOption()
	require.Equal(t, err.Code(), int32(0))
	require.Equal(t, err.Message(), "")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "<nil>")
}

func Test_Error(t *testing.T) {
	err := errorx.New(400, "请求参数错误")
	require.Error(t, err)
	require.Equal(t, err.Error(), "请求参数错误")

	innerErr := errors.New("内部错误1")
	err = err.TakeOption(errorx.WithCause(innerErr), errorx.WithMetadata("k1", "v1"))

	require.Equal(t, err.Code(), int32(400))
	require.Equal(t, err.Message(), "请求参数错误")
	require.Equal(t, err.Metadata(), map[string]string{"k1": "v1"})
	require.Equal(t, err.Error(), "请求参数错误: 内部错误1")

	err = err.TakeOption(errorx.WithError("内部错误2"))
	require.Equal(t, err.Error(), "请求参数错误: 内部错误2")

	err = err.TakeOption(errorx.WithErrorf("内部错误3"))
	require.Equal(t, err.Error(), "请求参数错误: 内部错误3")

	err = err.TakeOption(errorx.WithMessage("另一个错误"))
	require.Equal(t, err.Error(), "另一个错误: 内部错误3")

	err = err.TakeOption(errorx.WithMessagef("另一个错误(%v)", "设备号111"))
	require.Equal(t, err.Error(), "另一个错误(设备号111): 内部错误3")
}

type testError struct {
	s string
}

func (e *testError) Error() string { return e.s }

func newTestError(s string) *testError {
	return &testError{s: s}
}

func Test_Error_Unwrap(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		var err *errorx.Error

		gotErr := new(testError)
		ok := errors.As(err, &gotErr)
		require.False(t, ok)
	})
	t.Run("should not found", func(t *testing.T) {
		err1 := errorx.Newf(400, "请求参数错误")

		gotErr := new(testError)
		ok := errors.As(err1, &gotErr)
		require.False(t, ok)
	})
	t.Run("Unwrap", func(t *testing.T) {
		err := newTestError("内部错误")
		err1 := errorx.Newf(400, "请求参数错误").TakeOption(errorx.WithCause(err))

		gotErr := new(testError)
		ok := errors.As(err1, &gotErr)
		require.True(t, ok)
		require.Equal(t, gotErr.Error(), "内部错误")
	})
}

func Test_FromError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		var err error
		gotErr := errorx.Parse(err)
		require.Nil(t, gotErr)
	})
	t.Run("not Error", func(t *testing.T) {
		err := newTestError("内部错误")
		gotErr := errorx.Parse(err)

		require.Equal(t, gotErr.Code(), int32(500))
		require.Equal(t, gotErr.Message(), "服务器错误")
		require.Equal(t, gotErr.Unwrap(), err)
		require.Equal(t, gotErr.Error(), "服务器错误: 内部错误")
	})
	t.Run("Error", func(t *testing.T) {
		err := errorx.New(400, "请求参数错误")
		gotErr := errorx.Parse(err)

		require.Equal(t, gotErr.Code(), int32(400))
		require.Equal(t, gotErr.Message(), "请求参数错误")
		require.Equal(t, gotErr.Metadata(), map[string]string(nil))
		require.Equal(t, gotErr.Error(), "请求参数错误")
	})
}

func Test_Error_EqualCode(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		var err1 error
		require.False(t, errorx.EqualCode(err1, 400))
	})
	t.Run("nil Error", func(t *testing.T) {
		var err2 *errorx.Error
		require.False(t, errorx.EqualCode(err2, 400))
	})
	t.Run("not equal", func(t *testing.T) {
		err1 := new(testError)
		require.False(t, errorx.EqualCode(err1, 400))

		err2 := errorx.Newf(400, "请求参数错误1")
		require.False(t, errorx.EqualCode(err2, 500))
	})
	t.Run("equal", func(t *testing.T) {
		err1 := errorx.Newf(400, "请求参数错误1")
		require.True(t, errorx.EqualCode(err1, 400))
	})
}
