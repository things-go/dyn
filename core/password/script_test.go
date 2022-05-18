package password

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSCrypt(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		org := "hahaha"

		dst, err := GenerateSCryptFromPassword(org)
		t.Log(dst)
		require.Nil(t, err)
		require.NoError(t, CompareSCryptHashAndPassword(dst, org))
	})

	t.Run("not correct", func(t *testing.T) {
		org := "hahaha"

		dst, err := GenerateSCryptFromPassword(org)
		require.Nil(t, err)
		require.Error(t, CompareSCryptHashAndPassword(dst, "invalid"))
	})
}

func BenchmarkSCrypt_GenerateFromPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateSCryptFromPassword("hahaha")
	}
}

func BenchmarkSCrypt_CompareHashAndPassword(b *testing.B) {
	org := "hahaha"
	dst, _ := GenerateSCryptFromPassword(org)

	for i := 0; i < b.N; i++ {
		_ = CompareSCryptHashAndPassword(dst, org)
	}
}
