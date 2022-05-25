package proto

import (
	"testing"

	"github.com/stretchr/testify/require"

	testData "github.com/things-go/dyn/testdata/encoding"
)

func TestCodec(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		require.Equal(t, "proto", Name())
	})
	t.Run("Marshal/Unmarshal", func(t *testing.T) {
		want := testData.TestModel{
			Id:        1,
			Name:      "hello",
			Hobby:     []string{"study", "eat", "play"},
			SnakeCase: map[string]string{"1": "111", "2": "222"},
		}
		var got testData.TestModel

		tmp, err := Marshal(&want)
		require.NoError(t, err)

		err = Unmarshal(tmp, &got)
		require.NoError(t, err)
		require.Equal(t, want.Id, got.Id)
		require.Equal(t, want.Name, got.Name)
		require.Equal(t, want.Hobby, got.Hobby)
		require.Equal(t, want.SnakeCase, got.SnakeCase)
	})
}
