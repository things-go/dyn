package form

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/go-playground/form/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/things-go/dyn/testdata/encoding"
)

type LoginRequest struct {
	Username string `json:"username,omitempty" form:"uname"`
	Password string `json:"password,omitempty" form:"passwd"`
}

type TestModel struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func init() {
	encoder := form.NewEncoder()
	decoder := form.NewDecoder()
	codec := New(WithEncoder(encoder), WithDecoder(decoder), WithTagName("json"))
	ReplaceDefaultCodec(codec)
}

func TestNew(t *testing.T) {
	encoder := form.NewEncoder()
	decoder := form.NewDecoder()
	codec := New(WithEncoder(encoder), WithDecoder(decoder), WithTagName("form"))
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
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	tests := []struct {
		name string
		args any
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
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	p1 := TestDecode{}
	type args struct {
		vars   url.Values
		target any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    any
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

type NoProtoSub struct {
	Name string `json:"name"`
}

type NoProtoHello struct {
	Name string      `json:"name"`
	Sub  *NoProtoSub `json:"sub"`
}

func TestEncodeURL(t *testing.T) {
	type args struct {
		pathTemplate string
		msg          any
		needQuery    bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"proto: no any param",
			args{
				"http://helloworld.dev/sub",
				&encoding.HelloRequest{
					Name: "test",
					Sub:  &encoding.Sub{Name: "2233!!!"},
				},
				false,
			},
			`http://helloworld.dev/sub`,
		},
		{
			"proto: param",
			args{
				"http://helloworld.dev/{name}/sub/{sub.name}",
				&encoding.HelloRequest{
					Name: "test",
					Sub:  &encoding.Sub{Name: "2233!!!"},
				},
				false,
			},
			`http://helloworld.dev/test/sub/2233!!!`,
		},
		{
			"proto: param with proto json_name=naming",
			args{
				"http://helloworld.dev/{name}/sub/{sub.naming}",
				&encoding.HelloRequest{
					Name: "test",
					Sub:  &encoding.Sub{Name: "5566!!!"},
				},
				false,
			},
			`http://helloworld.dev/test/sub/5566!!!`,
		},
		{
			"proto: param with empty",
			args{
				"http://helloworld.dev/{name}/sub/{sub.name}",
				&encoding.HelloRequest{
					Name: "test",
				},
				false,
			},
			`http://helloworld.dev/test/sub/`,
		},
		{
			"proto: param not match",
			args{
				"http://helloworld.dev/{name}/sub/{sub.name33}",
				&encoding.HelloRequest{
					Name: "test",
				},
				false,
			},
			`http://helloworld.dev/test/sub/{sub.name33}`,
		},
		{
			"proto: param with query",
			args{
				"http://helloworld.dev/{name}/sub",
				&encoding.HelloRequest{
					Name: "go",
					Sub:  &encoding.Sub{Name: "golang"},
				},
				true,
			},
			`http://helloworld.dev/go/sub?sub.naming=golang`,
		},

		{
			"no proto: no any param",
			args{
				"http://helloworld.dev/sub",
				&NoProtoHello{
					Name: "test",
					Sub:  &NoProtoSub{Name: "2233!!!"},
				},
				false,
			},
			`http://helloworld.dev/sub`,
		},
		{
			"no proto: param",
			args{
				"http://helloworld.dev/{name}/sub/{sub.name}",
				&NoProtoHello{
					Name: "test",
					Sub:  &NoProtoSub{Name: "2233!!!"},
				},
				false,
			},
			`http://helloworld.dev/test/sub/2233!!!`,
		},
		{
			"no proto: param with empty",
			args{
				"http://helloworld.dev/{name}/sub/{sub.name}",
				&NoProtoHello{
					Name: "test",
				},
				false,
			},
			`http://helloworld.dev/test/sub/`,
		},
		{
			"no proto: param not match",
			args{
				"http://helloworld.dev/{name}/sub/{sub.name33}",
				&NoProtoHello{
					Name: "test",
				},
				false,
			},
			`http://helloworld.dev/test/sub/{sub.name33}`,
		},
		{
			"no proto: param with query",
			args{
				"http://helloworld.dev/{name}/sub",
				&NoProtoHello{
					Name: "go",
					Sub:  &NoProtoSub{Name: "golang"},
				},
				true,
			},
			`http://helloworld.dev/go/sub?sub.name=golang`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeURL(tt.args.pathTemplate, tt.args.msg, tt.args.needQuery); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
