package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/things-go/dyn/transport"
)

var _ transport.Transporter = (*Transport)(nil)

// Transport is a gRPC transport.
type Transport struct {
	fullMethod     string
	clientIp       string
	requestHeader  header
	responseHeader header
}

// Kind returns the transport kind.
func (tr *Transport) Kind() transport.Kind {
	return transport.GRPC
}

// FullPath Service full method or path
func (tr *Transport) FullPath() string {
	return tr.fullMethod
}

// ClientIp client ip
func (tr *Transport) ClientIp() string {
	return tr.clientIp
}

// RequestHeader returns the request header.
func (tr *Transport) RequestHeader() transport.Header {
	return tr.requestHeader
}

// ResponseHeader returns the reply header.
func (tr *Transport) ResponseHeader() transport.Header {
	return tr.responseHeader
}

type header metadata.MD

// Len returns the number of items in header.
func (h header) Len() int { return metadata.MD(h).Len() }

// Get returns the value associated with the passed key.
func (h header) Get(key string) string {
	vals := metadata.MD(h).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// Add adds the key, value pair to the header.
func (h header) Add(key, value string) { metadata.MD(h).Append(key, value) }

// Set stores the key-value pair.
func (h header) Set(key string, value string) { metadata.MD(h).Set(key, value) }

// Append adds the values to key k, not overwriting what was already stored at
// that key.
//
// k is converted to lowercase before storing in header.
func (h header) Append(key string, vals ...string) {
	metadata.MD(h).Append(key, vals...)
}

// Delete removes the values for a given key k which is converted to lowercase
// before removing it from header.
func (h header) Delete(key string) { metadata.MD(h).Delete(key) }

// Keys lists the keys stored in this carrier.
func (h header) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range metadata.MD(h) {
		keys = append(keys, k)
	}
	return keys
}

// Clone returns a copy of h or nil if h is nil.
func (h header) Clone() transport.Header { return transport.Header(header(metadata.MD(h).Copy())) }

// UnaryServerInterceptor is a gRPC unary server interceptor
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		p, _ := peer.FromContext(ctx)
		clientIp := ""
		if p != nil {
			clientIp = p.Addr.String()
		}
		responseHeader := metadata.MD{}
		ctx = transport.WithValueTransporter(ctx, &Transport{
			info.FullMethod,
			clientIp,
			header(md),
			header(responseHeader),
		})
		reply, err := handler(ctx, req)
		if len(responseHeader) > 0 {
			_ = grpc.SetHeader(ctx, responseHeader)
		}
		return reply, err
	}
}

// wrappedStream is rewrite grpc stream's context
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func NewWrappedStream(ctx context.Context, stream grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{
		stream,
		ctx,
	}
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// StreamServerInterceptor is a gRPC stream server interceptor
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, _ := metadata.FromIncomingContext(ss.Context())
		p, _ := peer.FromContext(ss.Context())
		clientIp := ""
		if p != nil {
			clientIp = p.Addr.String()
		}
		responseHeader := metadata.MD{}
		ctx := transport.WithValueTransporter(ss.Context(), &Transport{
			info.FullMethod,
			clientIp,
			header(md),
			header(responseHeader),
		})

		ws := NewWrappedStream(ctx, ss)
		err := handler(srv, ws)
		if len(responseHeader) > 0 {
			_ = grpc.SetHeader(ctx, responseHeader)
		}
		return err
	}
}
