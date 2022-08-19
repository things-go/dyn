package middleware

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/things-go/dyn/core/signature"
)

// 签名加/解密过程
// 随机生成一个 randomKey
// 如果body不为空并且需要加密 cipherBody = Base64(AesCBCEncrypt(randomKey,body)), 消息体传密文
// 如果body不为空并且不需要加密 cipherBody = Base64(body), 消息体传原文
// 拼接签名串 str = timestamp + method + url + cipherBody
// sign = Base64(HMAC(randomKey,str))
// secret = Base64(RsaEncrypt(randomKey, pubkey))
// 服务端验签解密过程则是上述的逆过程
//
// 请求头:
//  Timestamp: 时间戳, 单位ms
//  Secret: randomKey的密文, 即上述 secret
//  AesCbcEncrypt: 是否加密消息体, 1: 表示加密, 0: 表示不加密
//  Sign: 签名, 即上述 sign

// SignOption 签名选项
type SignOption func(*SignConfig)

// SignConfig 签名配置
type SignConfig struct {
	privKey       *rsa.PrivateKey
	availWindow   time.Duration
	skip          func(*gin.Context) bool
	errorFallback func(*gin.Context, error)
}

// PrivKey 设置私钥
func PrivKey(privKey *rsa.PrivateKey) SignOption {
	return func(o *SignConfig) {
		o.privKey = privKey
	}
}

// MustPrivKeyFromFile 设置私钥
func MustPrivKeyFromFile(privKeyFile string) SignOption {
	keyData, err := os.ReadFile(privKeyFile)
	if err != nil {
		panic(err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		panic(err)
	}
	return PrivKey(key)
}

// WithAvailWindow 有效窗口时间, 小于等于0表示不验证
func WithAvailWindow(t time.Duration) SignOption {
	return func(o *SignConfig) {
		o.availWindow = t
	}
}

// WithSkip 忽略验证的接口
func WithSkip(skip func(c *gin.Context) bool) SignOption {
	return func(o *SignConfig) {
		if skip != nil {
			o.skip = skip
		}
	}
}

// WithUnauthorizedFallback sets the fallback handler when requests are unauthorized.
func WithErrorFallback(f func(c *gin.Context, err error)) SignOption {
	return func(o *SignConfig) {
		if f != nil {
			o.errorFallback = f
		}
	}
}

// VerifySign 签名验证器
func VerifySign(opts ...SignOption) gin.HandlerFunc {
	cfg := SignConfig{
		skip: func(c *gin.Context) bool { return false },
		errorFallback: func(c *gin.Context, err error) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "无效请求",
				"detail":  err.Error(),
			})
		},
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return func(c *gin.Context) {
		if cfg.skip(c) {
			c.Next()
			return
		}

		timestamp := c.GetHeader("Timestamp") // 时间戳
		secret := c.GetHeader("Secret")       // 密钥
		if timestamp == "" || secret == "" {
			cfg.errorFallback(c, errors.New("无效的timestamp, secret"))
			return
		}
		milliTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			cfg.errorFallback(c, errors.New("无效的timestamp"))
			return
		}
		if cfg.availWindow > 0 && time.Now().Sub(time.UnixMilli(milliTimestamp)) > cfg.availWindow {
			cfg.errorFallback(c, errors.New("该请求已过期失效"))
			return
		}

		// 根据secret解出randomKey
		randomKey, err := signature.RsaDecrypt(cfg.privKey, secret)
		if err != nil {
			cfg.errorFallback(c, err)
			return
		}

		// body处理
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			cfg.errorFallback(c, err)
			return
		}
		var origBody = body
		var cipherBody string

		if len(body) > 0 {
			encrypt := c.GetHeader("AesCbcEncrypt") // 是否加密body
			if encrypt == "1" {
				cipherBody = *(*string)(unsafe.Pointer(&body))
				origBody, err = signature.AesCbcDecrypt(randomKey, cipherBody)
			} else {
				cipherBody = base64.StdEncoding.EncodeToString(body)
				origBody = body
			}
		}
		if err != nil {
			cfg.errorFallback(c, err)
			return
		}
		c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(origBody))

		str := timestamp + strings.ToUpper(c.Request.Method) + c.Request.RequestURI + cipherBody
		calcSign := signature.HmacSha256(randomKey, str)
		sign := c.GetHeader("Sign")

		// log.Printf("加密串: %s\r\n", str)
		// log.Printf("客户端上传的sign: %s\r\n", sign)
		// log.Printf("计算得到的sign: %s\r\n", calcSign)
		if calcSign != sign {
			cfg.errorFallback(c, errors.New("无效的签名"))
			return
		}
		c.Next()
	}
}
