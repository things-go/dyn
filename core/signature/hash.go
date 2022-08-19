package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

// HmacSha1 hmac sha1 with base64 encoded.
func HmacSha1(key, str string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// HmacSha256 hmac sha256 with base64 encoded.
func HmacSha256(key, str string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Hmac512 hmac sha1 with base64 encoded.
func Hmac512(key, str string) string {
	h := hmac.New(sha512.New, []byte(key))
	h.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// HexSha1 sha1 with hex encoded.
func HexSha1(str string) string {
	bs := sha1.Sum([]byte(str))
	return hex.EncodeToString(bs[:])
}

// HexSha256 sha256 with hex encoded.
func HexSha256(str string) string {
	bs := sha256.Sum256([]byte(str))
	return hex.EncodeToString(bs[:])
}

// HexSha512 sha512 with hex encoded.
func HexSha512(str string) string {
	bs := sha512.Sum512([]byte(str))
	return hex.EncodeToString(bs[:])
}
