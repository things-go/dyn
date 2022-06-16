package metadata

import (
	"context"
)

type ctxMetadataKey struct{}

// NewContext creates a new context with client md attached.
func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, ctxMetadataKey{}, md)
}

// FromContext returns the metadata in ctx if it exists.
func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(ctxMetadataKey{}).(Metadata)
	return md, ok
}

// MergeContext merge new metadata into ctx.
// if patchMd key exists, but value is empty, means delete it.
// overwrite flag means if key exist, it will be overwritten.
func MergeContext(ctx context.Context, patchMd Metadata, overwrite bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	md, ok := FromContext(ctx)
	if ok {
		md = md.Clone()
	} else {
		md = Metadata{}
	}
	for k, v := range patchMd {
		_, ok = md[k]
		switch {
		case !ok:
			md[k] = v
		case v == "":
			delete(md, k)
		case overwrite:
			md[k] = v
		default:
			// skip
		}
	}
	return NewContext(ctx, md)
}
