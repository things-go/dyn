//go:generate stringer -type=Deploy -linecomment
package deploy

import (
	"log"
)

type Deploy int

const (
	None       Deploy = iota // none
	Local                    // local
	Develop                  // dev
	Debug                    // debug
	Uat                      // uat
	Production               // prod
)

var deploy = None

func Convert(m string) Deploy {
	switch m {
	case Local.String():
		return Local
	case Develop.String():
		return Develop
	case Debug.String():
		return Debug
	case Uat.String():
		return Uat
	case Production.String():
		return Production
	default:
		return None
	}
}

// Set 设置布署模式
func Set(m Deploy) {
	deploy = m
}

// Get 获取当前的布署模式
func Get() Deploy { return deploy }

// IsLocal 是否本地模式
func IsLocal() bool { return deploy == Local }

// IsDevelop 是否开发模式
func IsDevelop() bool { return deploy == Develop }

// IsDebug 是否调试模式
func IsDebug() bool { return deploy == Debug }

// IsUat 是否预发布模式
func IsUat() bool { return deploy == Uat }

// IsProduction 是否生产模式
func IsProduction() bool { return deploy == Production }

// IsTest 测试: 本地,开发或者调试
func IsTest() bool { return IsLocal() || IsDevelop() || IsDebug() }

// IsRelease 预发或者生产环境
func IsRelease() bool { return IsUat() || IsProduction() }

// MustSetDeploy 设置布署模式, 不得为 unknown 模式, 否则panic
func MustSetDeploy(m string) {
	Set(Convert(m))
	CheckMustDeploy()
}

// GetDeploy 获取当前的布署模式
func GetDeploy() string {
	return Get().String()
}

// CheckMustDeploy 校验当前的布署环境必须设置非 unknown 模式, 否则panic
func CheckMustDeploy() {
	if deploy == None {
		log.Fatalf("Please set deploy mode first, must be one of local, dev, debug, uat, prod")
	}
}
