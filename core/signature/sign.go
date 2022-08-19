package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"time"
)

type H map[string]any

func SignHexSha256(mp map[string]any, secret string) string {
	return Sign(mp, secret, func(src string) string {
		bs := sha256.Sum256([]byte(src))
		return hex.EncodeToString(bs[:])
	})
}

func SignHmacSha256(mp map[string]any, secret string) string {
	return Sign(mp, secret, func(src string) string {
		bs := hmac.New(sha256.New, []byte(src)).Sum(nil)
		return base64.StdEncoding.EncodeToString(bs)
	})
}

func Sign(mp map[string]any, secret string, sign func(string) string) string {
	src := ConcatMap(mp, false) + secret
	return sign(src)
}

// Iat 签发时间字符串
func Iat() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// VerifyIat 验证签发时间是否在有效期内
func VerifyIat(iat string, timeout time.Duration) bool {
	ns, err := strconv.ParseInt(iat, 10, 64)
	if err != nil {
		return false
	}
	t := time.Unix(ns/int64(time.Second), ns%int64(time.Second))
	return t.Add(timeout).After(time.Now())
}

func IatSign(mp map[string]any) (iat, sign string) {
	iat = Iat()
	return iat, SignHmacSha256(mp, iat)
}

// VerifySign 验证签发时间是否在有效期内, 并验证签名是否正确
func VerifyIatSign(iat, targetSign string, iatTimout time.Duration, mp map[string]any) bool {
	return VerifyIat(iat, iatTimout) && targetSign == SignHmacSha256(mp, iat)
}
