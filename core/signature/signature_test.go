package signature

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

var aesKeySizes = []int{16, 24, 32}

func TestEncryptDecrypt(t *testing.T) {
	plainText := []byte("helloworld,this is golang language. welcome")
	for _, keySize := range aesKeySizes {
		key := make([]byte, keySize)
		_, err := io.ReadFull(rand.Reader, key)
		require.NoError(t, err)

		cipherText, err := AesCbcEncrypt(string(key), plainText)
		require.NoError(t, err)

		got, err := AesCbcDecrypt(string(key), cipherText)
		require.NoError(t, err)
		require.Equal(t, plainText, got)
	}
}

var pri = `
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAuC0yzDEWZ/ZF492cxgpSbcki2HCuhJACWi6I2m/AZHiTPU95
IiPDbI23PtvAD7pq9nMypFM3d/t5WiNRBQht+W+eDUSxFdWkHlUTZv5xNZn4H0ZS
4rO8/69DwSjrhbO3Gh/wykyajPUvkU1ispJ0jqnnLtVhp8zWPfktBKodMQMLSxSY
V9CHDBriUfvI9eUcFDtiL1IA36u/0s8sO1ril9jp/ap2gbLKNRIp/TtiO2p5oS73
sJZJyA0aBhcr8zIHm1vp42KIqyKMwDlF24OEdP5Xhgme0konTglTU4OPzbzRPTzg
FWeWgr1BoQqtjQ36kzpJgpFIbZq89sQGa6NcNwIDAQABAoIBAFJkNeNO2l0ruNxV
GXsmxvVCE4qL7iZzXfoF80b00zafEg9NbD0vjo8TXrSKDgN7O7qAklkcuSW9o+4E
+our17bMjgIVhrrj1gWTvZhHn1GBTZSAGzg7fANN3pucuLAZU4ImY8u8NS91yA/S
dKK6DdZ8f2VeI8+zPOyAefkqOODhX1qayqTDUClYKTP5FnuFx+L2AXNccyBogegm
ZNCJ2nUVBbLWikIAZGaJ3dPS6v69S0jB5xPTN0tj7m7DSv0QkzhM1xucmcZZmLMG
yd/JtB/87Rgr4jyTxd289rXg1pG8PwAy9VQYsZ78wnThRD39f/kk88d0M34gdn+y
FM4inpECgYEA4JtMJWUY7IyNsww7sFMIc6IEnLDqZD8Tf6fhk7WfiEDUnRxL8TWH
pySwbFBA7JeWqhJJ4Tj8bxb/3epjrA5QKLuKz4oaei4vChIhD3eLppE9+WibA/67
RMv54OuYhey9uq/FXVYHQ5TlZq3lMOnqhOT2y6wrm799u4h8e2G0uLMCgYEA0etD
VHpXKTkeIN0fepWoOEDGdX+Kbh0zhHZeWa4RLZt2I6Hi/ZUkveOIYoSdARcjM09r
6KUXX5W6UZ4NbBZO7cNeMnZPAzLA6LcBahzlDyPRBL5nVhpSJGuTzNfLWwcPZIqe
4MNFPYEBmSrSswIXPEWhdJz1lZuVsQtKUMK4aG0CgYBuotCdUvE2A/4AhjQYpK3z
F4miDVtHyfI23WE2Oy68FQMl6LxXsoCBiocEs3tnjzv9xkhyEnn11qRukXhLVjmR
9t9nX6WvLXSqR0fVsJMvlzep1ScWjrF8L+WEL0jQH09N2Csl0Kx/U6a0L1BICdEl
aQtQRByu+WJbr91xgS1eFQKBgBEFZIY7DUo4aWr8wwqri+JTzkPEvbLEB2NcPbZD
2Py7uE6XV9J7/2iuRGbInfpyp9YHQJaynDyR5XOsvyXegTPiPYcV9L4rpVy5ShIS
mbgqjU43KiXfKH3vgyJ+9OxCnErouo07CCg+h6SlxkPhjYTDmJ3eBEPHQ9IBOltm
DpHZAoGAKCOX7ub+E2wx2cqrjQZ7KCGOiEEZEdsRulblxks8VZa4tVbDCoUj3X16
XyvtdW6U02yy4rls18un3yrhJk7HE/21NvQPeSEJcS0bl+//61nNaUBzqbjWM4Ii
vrdhOYn2ywSTuU0Frna2ODulw/d3gS+SqKlSfoVvuG/+xQwrp9M=
-----END RSA PRIVATE KEY-----
`

var pub = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuC0yzDEWZ/ZF492cxgpS
bcki2HCuhJACWi6I2m/AZHiTPU95IiPDbI23PtvAD7pq9nMypFM3d/t5WiNRBQht
+W+eDUSxFdWkHlUTZv5xNZn4H0ZS4rO8/69DwSjrhbO3Gh/wykyajPUvkU1ispJ0
jqnnLtVhp8zWPfktBKodMQMLSxSYV9CHDBriUfvI9eUcFDtiL1IA36u/0s8sO1ri
l9jp/ap2gbLKNRIp/TtiO2p5oS73sJZJyA0aBhcr8zIHm1vp42KIqyKMwDlF24OE
dP5Xhgme0konTglTU4OPzbzRPTzgFWeWgr1BoQqtjQ36kzpJgpFIbZq89sQGa6Nc
NwIDAQAB
-----END PUBLIC KEY-----

`

func TestRsaEncryptDecrypt(t *testing.T) {
	priKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pri))
	require.NoError(t, err)
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pub))
	require.NoError(t, err)

	want := "helloworld,this is golang language. welcome"

	mid, err := RsaEncrypt(pubKey, want)
	require.NoError(t, err)

	got, err := RsaDecrypt(priKey, mid)
	require.NoError(t, err)

	require.Equal(t, want, got)
}

func TestPCKS(t *testing.T) {
	want := []byte("1234567890abcdef")

	b := PCKSPadding(want, 16)

	got, err := PCKSUnPadding(b, 16)
	require.NoError(t, err)
	require.Equal(t, want, got)
}
