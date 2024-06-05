//go:generate errno-gen -t BizError
package errno

type BizError int // 业务错误

const (
	Timeout      BizError = 1000 + iota // 操作超时
	UserNotExist                        // 用户不存在
)
