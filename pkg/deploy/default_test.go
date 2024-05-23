package deploy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Mode_Parse(t *testing.T) {
	tests := []struct {
		name string
		args string
		want Mode
	}{

		{
			name: "Invalid Mode String",
			args: "Invalid",
			want: None,
		},
		{
			name: "DEV",
			args: "DEV",
			want: Dev,
		},
		{
			name: "dev",
			args: "dev",
			want: Dev,
		},
		{
			name: "test",
			args: "test",
			want: Test,
		},
		{
			name: "uat",
			args: "uat",
			want: Uat,
		},
		{
			name: "prod",
			args: "prod",
			want: Prod,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args); got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Mode_MustSet(t *testing.T) {
	require.Panics(t, func() {
		MustSet("Invalid")
	})
	require.NotPanics(t, func() {
		MustSet("dev")
	})
}

func Test_Mode_Is(t *testing.T) {
	Set(Dev)
	m := Get()
	require.True(t, m.IsDev())
	require.True(t, IsDev())
	require.True(t, Is(Dev))

	require.True(t, Valid())
	require.True(t, IsDev())
	require.False(t, IsTest())
	require.False(t, IsUat())
	require.False(t, IsProd())
	require.True(t, IsTesting())
	require.False(t, IsRelease())
}

func Test_Mode_Invalid(t *testing.T) {
	var m Mode = 5

	require.Equal(t, "Mode(5)", m.String())
}
