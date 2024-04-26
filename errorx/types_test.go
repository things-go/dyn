package errorx_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/things-go/dyn/errorx"
)

func Test_Types(t *testing.T) {
	err := errorx.NewBadRequest()
	require.Equal(t, err.Code(), int32(400))
	require.Equal(t, err.Message(), "请求参数错误")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "请求参数错误")

	err = errorx.NewUnauthorized()
	require.Equal(t, err.Code(), int32(401))
	require.Equal(t, err.Message(), "未授权")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "未授权")

	err = errorx.NewForbidden()
	require.Equal(t, err.Code(), int32(403))
	require.Equal(t, err.Message(), "禁止访问")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "禁止访问")

	err = errorx.NewNotFound()
	require.Equal(t, err.Code(), int32(404))
	require.Equal(t, err.Message(), "没有找到,资源不存在")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "没有找到,资源不存在")

	err = errorx.NewMethodNotAllowed()
	require.Equal(t, err.Code(), int32(405))
	require.Equal(t, err.Message(), "方法不允许")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "方法不允许")

	err = errorx.NewRequestTimeout()
	require.Equal(t, err.Code(), int32(408))
	require.Equal(t, err.Message(), "请求超时")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "请求超时")

	err = errorx.NewConflict()
	require.Equal(t, err.Code(), int32(409))
	require.Equal(t, err.Message(), "资源冲突")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "资源冲突")

	err = errorx.NewInternalServer()
	require.Equal(t, err.Code(), int32(500))
	require.Equal(t, err.Message(), "服务器错误")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "服务器错误")

	err = errorx.NewNotImplemented()
	require.Equal(t, err.Code(), int32(501))
	require.Equal(t, err.Message(), "未实现")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "未实现")

	err = errorx.NewBadGateway()
	require.Equal(t, err.Code(), int32(502))
	require.Equal(t, err.Message(), "网关错误")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "网关错误")

	err = errorx.NewServiceUnavailable()
	require.Equal(t, err.Code(), int32(503))
	require.Equal(t, err.Message(), "服务器不可用")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "服务器不可用")

	err = errorx.NewGatewayTimeout()
	require.Equal(t, err.Code(), int32(504))
	require.Equal(t, err.Message(), "网关超时")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "网关超时")

	err = errorx.NewClientClosed()
	require.Equal(t, err.Code(), int32(499))
	require.Equal(t, err.Message(), "客户端关闭")
	require.Equal(t, err.Metadata(), map[string]string(nil))
	require.Equal(t, err.Error(), "客户端关闭")
}
