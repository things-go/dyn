package form

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/go-playground/form/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type LoginRequest struct {
	Username string `form:"username,omitempty" json:"uname"`
	Password string `form:"password,omitempty" json:"passwd"`
}

type TestModel struct {
	ID   int32  `form:"id"`
	Name string `form:"name"`
}

func TestNew(t *testing.T) {
	encoder := form.NewEncoder()
	encoder.SetTagName("json")
	decoder := form.NewDecoder()
	decoder.SetTagName("json")
	codec := New(WithEncoder(encoder), WithDecoder(decoder))
	req := &LoginRequest{
		Username: "username",
		Password: "password",
	}
	content, err := codec.Marshal(req)
	require.NoError(t, err)
	require.Equal(t, []byte("passwd=password&uname=username"), content)
}

func TestFormCodec(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		require.Equal(t, "x-www-form-urlencoded", Name())
	})
	t.Run("Marshal", func(t *testing.T) {
		req := &LoginRequest{
			Username: "username",
			Password: "password",
		}
		content, err := Marshal(req)
		require.NoError(t, err)
		require.Equal(t, []byte("password=password&username=username"), content)

		req = &LoginRequest{
			Username: "username",
			Password: "",
		}
		content, err = Marshal(req)
		require.NoError(t, err)
		require.Equal(t, []byte("username=username"), content)

		m := &TestModel{
			ID:   1,
			Name: "username",
		}
		content, err = Marshal(m)
		require.NoError(t, err)
		require.Equal(t, []byte("id=1&name=username"), content)
	})
	t.Run("Unmarshal", func(t *testing.T) {
		want := &LoginRequest{
			Username: "username",
			Password: "password",
		}
		got := new(LoginRequest)
		err := Unmarshal([]byte(`password=password&username=username`), got)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})
	t.Run("Marshal/Unmarshal", func(t *testing.T) {
		want := &LoginRequest{
			Username: "username",
			Password: "password",
		}
		content, err := Marshal(want)
		require.NoError(t, err)

		got := new(LoginRequest)
		err = Unmarshal(content, got)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})
}
func TestEncode(t *testing.T) {
	type TestEncode struct {
		Name string `form:"name"`
		URL  string `form:"url"`
	}
	tests := []struct {
		name string
		args interface{}
		want url.Values
	}{
		{
			"test",
			TestEncode{
				Name: "test",
				URL:  "https://go.dev",
			},
			url.Values{
				"name": []string{"test"},
				"url":  []string{"https://go.dev"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "Encode(%v)", tt.args)
		})
	}
}
func TestDecode(t *testing.T) {
	type TestDecode struct {
		Name string `form:"name"`
		URL  string `form:"url"`
	}
	p1 := TestDecode{}
	type args struct {
		vars   url.Values
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			name: "test",
			args: args{
				vars:   map[string][]string{"name": {"golang"}, "url": {"https://go.dev"}},
				target: &p1,
			},
			wantErr: false,
			want:    &TestDecode{"golang", "https://go.dev"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.vars, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("Decode() target = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}
