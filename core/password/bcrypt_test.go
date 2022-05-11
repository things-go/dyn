package password

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBCrypt(t *testing.T) {
	org := "hahaha"

	dst, err := GenerateFromPassword(org)
	t.Log(dst)
	require.Nil(t, err)
	require.Nil(t, CompareHashAndPassword(dst, org))
}

func BenchmarkBCrypt_GenerateFromPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateFromPassword("hahaha")
	}
}

func BenchmarkBCrypt_CompareHashAndPassword(b *testing.B) {
	org := "hahaha"
	dst, _ := GenerateFromPassword(org)

	for i := 0; i < b.N; i++ {
		_ = CompareHashAndPassword(dst, org)
	}
}
