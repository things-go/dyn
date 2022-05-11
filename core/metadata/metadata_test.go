package metadata

import (
	"context"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		mds []map[string]string
	}
	tests := []struct {
		name string
		args args
		want Metadata
	}{
		{
			name: "hello",
			args: args{[]map[string]string{{"hello": "metadata"}, {"hello2": "test-metadata"}}},
			want: Metadata{"hello": "metadata", "hello2": "test-metadata"},
		},
		{
			name: "hi",
			args: args{[]map[string]string{{"hi": "metadata"}, {"hi2": "test-metadata"}}},
			want: Metadata{"hi": "metadata", "hi2": "test-metadata"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.mds...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    Metadata
		args args
		want string
	}{
		{
			name: "metadata",
			m:    Metadata{"metadata": "value", "env": "dev"},
			args: args{key: "metadata"},
			want: "value",
		},
		{
			name: "env",
			m:    Metadata{"metadata": "value", "env": "dev"},
			args: args{key: "env"},
			want: "dev",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Get(tt.args.key); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		m    Metadata
		args args
		want Metadata
	}{
		{
			name: "ignore",
			m:    Metadata{},
			args: args{key: "hello", value: ""},
			want: Metadata{},
		},
		{
			name: "metadata",
			m:    Metadata{},
			args: args{key: "hello", value: "metadata"},
			want: Metadata{"hello": "metadata"},
		},
		{
			name: "env",
			m:    Metadata{"hello": "metadata"},
			args: args{key: "env", value: "pro"},
			want: Metadata{"hello": "metadata", "env": "pro"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Set(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("Set() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestMetadata_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    Metadata
		args args
		want Metadata
	}{
		{
			name: "metadata",
			m:    Metadata{"metadata": "value", "env": "dev"},
			args: args{key: "metadata"},
			want: Metadata{"env": "dev"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.m.Delete(tt.args.key); !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("Delete() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestMetadata_Exist(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    Metadata
		args args
		want bool
	}{
		{
			name: "exist",
			m:    Metadata{"metadata": "value"},
			args: args{key: "metadata"},
			want: true,
		},
		{
			name: "not exist",
			m:    Metadata{"metadata": "value"},
			args: args{key: "env"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Exist(tt.args.key); got != tt.want {
				t.Errorf("Exist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_Range(t *testing.T) {
	md := Metadata{"metadata": "metadata", "https://test-metadata.dev": "https://test-metadata.dev", "test-metadata": "test-metadata"}
	tmp := Metadata{}
	md.Range(func(k, v string) bool {
		if k == "https://test-metadata.dev" || k == "metadata" {
			tmp[k] = v
		}
		if len(tmp) == 2 {
			return false
		}
		return true
	})
	if !reflect.DeepEqual(tmp, Metadata{"https://test-metadata.dev": "https://test-metadata.dev", "metadata": "metadata"}) {
		t.Errorf("metadata = %v, want %v", tmp, Metadata{"metadata": "metadata"})
	}
}

func TestMetadata_Clone(t *testing.T) {
	tests := []struct {
		name string
		m    Metadata
		want Metadata
	}{
		{
			name: "metadata",
			m:    Metadata{"metadata": "metadata", "https://test-metadata.dev": "https://test-metadata.dev", "test-metadata": "test-metadata"},
			want: Metadata{"metadata": "metadata", "https://test-metadata.dev": "https://test-metadata.dev", "test-metadata": "test-metadata"},
		},
		{
			name: "go",
			m:    Metadata{"language": "golang"},
			want: Metadata{"language": "golang"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Clone()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
			got["metadata"] = "go"
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("want got != want got %v want %v", got, tt.want)
			}
		})
	}
}

func TestContext(t *testing.T) {
	type args struct {
		ctx context.Context
		md  Metadata
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "metadata",
			args: args{context.Background(), Metadata{"hello": "metadata", "metadata": "https://test-metadata.dev"}},
		},
		{
			name: "hello",
			args: args{context.Background(), Metadata{"hello": "metadata", "hello2": "https://test-metadata.dev"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext(tt.args.ctx, tt.args.md)
			m, ok := FromContext(ctx)
			if !ok {
				t.Errorf("FromContext() = %v, want %v", ok, true)
			}

			if !reflect.DeepEqual(m, tt.args.md) {
				t.Errorf("meta = %v, want %v", m, tt.args.md)
			}
		})
	}
}
