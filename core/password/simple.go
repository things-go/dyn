package password

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"unsafe"
)

// Simple password encryption
type Simple struct{}

// GenerateFromPassword generate password hash encryption 加盐法
func GenerateSimpleFromPassword(password string) (string, error) {
	unencodedSalt := make([]byte, maxSaltSize)
	_, err := rand.Read(unencodedSalt)
	if err != nil {
		return "", err
	}

	pwd := simpleHash(password, unencodedSalt)
	out := make([]byte, base64.StdEncoding.EncodedLen(len(pwd)))
	base64.StdEncoding.Encode(out, pwd)
	return *(*string)(unsafe.Pointer(&out)), nil
}

// CompareHashAndPassword Compare password hash verification
func CompareSimpleHashAndPassword(hashedPassword, password string) error {
	orgRb, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return err
	}

	if len(orgRb) < maxSaltSize {
		return ErrCompareFailed
	}
	unencodedSalt := orgRb[:maxSaltSize]

	pwd := simpleHash(password, unencodedSalt)
	if !bytes.Equal(orgRb, pwd) {
		return ErrCompareFailed
	}
	return nil
}

func simpleHash(password string, salt []byte) []byte {
	bs := make([]byte, 0, len(salt)+sha256.BlockSize)

	bs = append(bs, salt...)
	bs = append(bs, password...)

	mdv := hmac.New(sha256.New, salt).Sum(bs)

	bs = bs[:len(salt)]
	bs = append(bs, mdv...)
	return bs
}
