package metadata

import (
	"context"
	"reflect"
	"testing"
)

func TestMergeContext(t *testing.T) {
	type args struct {
		ctx       context.Context
		patchMd   Metadata
		overwrite bool
	}
	tests := []struct {
		name string
		args args
		want Metadata
	}{
		{
			name: "metadata",
			args: args{
				NewContext(context.Background(), Metadata{}),
				Metadata{"hello": "metadata", "env": "dev"},
				false,
			},
			want: Metadata{"hello": "metadata", "env": "dev"},
		},
		{
			name: "hello",
			args: args{
				NewContext(context.Background(), Metadata{"hi": "https://test-metadata.dev"}),
				Metadata{"hello": "metadata", "env": "dev"},
				false,
			},
			want: Metadata{"hello": "metadata", "env": "dev", "hi": "https://test-metadata.dev"},
		},
		{
			name: "delete hi",
			args: args{
				NewContext(context.Background(), Metadata{"hi": "https://test-metadata.dev"}),
				Metadata{"hello": "metadata", "env": "dev", "hi": ""},
				false,
			},
			want: Metadata{"hello": "metadata", "env": "dev"},
		},
		{
			name: "overwrite hi",
			args: args{
				NewContext(context.Background(), Metadata{"hi": "https://test-metadata.dev"}),
				Metadata{"hello": "metadata", "env": "dev", "hi": "hi"},
				true,
			},
			want: Metadata{"hello": "metadata", "env": "dev", "hi": "hi"},
		},
		{
			name: "skip hi",
			args: args{
				NewContext(context.Background(), Metadata{"hi": "https://test-metadata.dev"}),
				Metadata{"hello": "metadata", "env": "dev", "hi": "hi"},
				false,
			},
			want: Metadata{"hello": "metadata", "env": "dev", "hi": "https://test-metadata.dev"},
		},
		{
			name: "nil context",
			args: args{
				nil,
				Metadata{"hello": "metadata", "env": "dev", "hi": "hi"},
				false,
			},
			want: Metadata{"hello": "metadata", "env": "dev", "hi": "hi"},
		},
		{
			name: "empty context",
			args: args{
				context.Background(),
				Metadata{"hello": "metadata", "env": "dev", "hi": "hi"},
				false,
			},
			want: Metadata{"hello": "metadata", "env": "dev", "hi": "hi"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := MergeContext(tt.args.ctx, tt.args.patchMd, tt.args.overwrite)
			md, ok := FromContext(ctx)
			if !ok {
				t.Errorf("FromContext() = %v, want %v", ok, true)
			}
			if !reflect.DeepEqual(md, tt.want) {
				t.Errorf("metadata = %v, want %v", md, tt.want)
			}
		})
	}
}
