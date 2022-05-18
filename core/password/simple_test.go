package password

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		org := "hahaha"

		dst, err := GenerateSimpleFromPassword(org)
		t.Log(dst)
		require.Nil(t, err)
		require.NoError(t, CompareSimpleHashAndPassword(dst, org))
	})

	t.Run("not correct", func(t *testing.T) {
		org := "hahaha"

		dst, err := GenerateSimpleFromPassword(org)
		require.Nil(t, err)
		require.Error(t, CompareSimpleHashAndPassword(dst, "invalid"))
	})
}

func BenchmarkSimple_GenerateFromPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateSimpleFromPassword("hahaha")
	}
}

func BenchmarkSimple_CompareHashAndPassword(b *testing.B) {
	org := "hahaha"

	dst, _ := GenerateSimpleFromPassword(org)
	for i := 0; i < b.N; i++ {
		_ = CompareSimpleHashAndPassword(dst, org)
	}
}
