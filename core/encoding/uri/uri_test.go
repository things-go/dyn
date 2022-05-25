package uri

import (
	"testing"

	"github.com/things-go/dyn/testdata/encoding"
)

func TestEncode(t *testing.T) {
	type args struct {
		pathTemplate string
		msg          interface{}
		needQuery    bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"",
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
			"",
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
			"",
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
			"",
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
			"",
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
			"",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encode(tt.args.pathTemplate, tt.args.msg, tt.args.needQuery); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
