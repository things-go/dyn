package signature

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"unsafe"
)

// error defined
var (
	ErrInputInvalidLength     = errors.New("encoded message length must be more than zero")
	ErrInputNotMoreABlock     = errors.New("decoded message length must be more than a block size")
	ErrInputNotMultipleBlocks = errors.New("decoded message length must be multiple of block size")
	ErrInvalidIvSize          = errors.New("iv length must equal block size")
	ErrUnPaddingOutOfRange    = errors.New("unPadding out of range")
)

// Signature hmac sha256的base64
func Signature(key, str string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// RsaEncrypt rsa PKCS #1 v1.5. and base64 encoded.
func RsaEncrypt(pub *rsa.PublicKey, rawText string) (string, error) {
	b, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(rawText))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// RsaDecrypt  base64 decoded and rsa PKCS #1 v1.5.
func RsaDecrypt(pri *rsa.PrivateKey, ciphertext string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	bb, err := rsa.DecryptPKCS1v15(rand.Reader, pri, b)
	if err != nil {
		return "", err
	}
	return *(*string)(unsafe.Pointer(&bb)), nil
}

// Encrypt aes cbc, iv + ciphertext base64 encoded.
// key must 16, 24, 32
func Encrypt(key string, rawText []byte) (string, error) {
	bsKey := []byte(key)
	cip, err := aes.NewCipher(bsKey)
	if err != nil {
		return "", err
	}
	blockSize := cip.BlockSize()

	orig := PCKSPadding(rawText, blockSize)
	cipherText := make([]byte, blockSize+len(orig))

	// 生成随机iv
	_, err = rand.Read(cipherText[:blockSize])
	if err != nil {
		return "", err
	}
	iv := cipherText[:blockSize]
	cipher.NewCBCEncrypter(cip, iv).CryptBlocks(cipherText[blockSize:], orig)
	// iv + ciphertext 一起进行 sha256
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt aes cbc,  base64 decoded iv + ciphertext.
// key must 16, 24, 32
func Decrypt(key, cipherText string) ([]byte, error) {
	body, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	cip, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	blockSize := cip.BlockSize()
	if len(body) == 0 || len(body)%blockSize != 0 {
		return nil, ErrInputNotMultipleBlocks
	}
	iv, msg := body[:blockSize], body[blockSize:]
	cipher.NewCBCDecrypter(cip, iv).CryptBlocks(msg, msg)
	return PCKSUnPadding(msg, blockSize)
}

// PCKSPadding PKCS#5和PKCS#7 填充
func PCKSPadding(origData []byte, blockSize int) []byte {
	padSize := blockSize - len(origData)%blockSize
	padText := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(origData, padText...)
}

// PCKSUnPadding PKCS#5和PKCS#7 解填充
func PCKSUnPadding(origData []byte, blockSize int) ([]byte, error) {
	orgLen := len(origData)
	if orgLen == 0 {
		return nil, ErrUnPaddingOutOfRange
	}
	unPadSize := int(origData[orgLen-1])
	if unPadSize > blockSize || unPadSize > orgLen {
		return nil, ErrUnPaddingOutOfRange
	}
	for _, v := range origData[orgLen-unPadSize:] {
		if v != byte(unPadSize) {
			return nil, ErrUnPaddingOutOfRange
		}
	}
	return origData[:(orgLen - unPadSize)], nil
}
