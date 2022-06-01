package transport

import (
	"context"
)

type Propagator interface {
	// Kind transporter
	// grpc
	// http
	Kind() Kind
	// FullMethod Service full method or path
	FullMethod() string
	// Route Service full route
	Route() string
}

// Kind defines the type of Transport
type Kind string

func (k Kind) String() string { return string(k) }

// Defines a set of transport kind
const (
	KindGRPC Kind = "grpc"
	KindHTTP Kind = "http"
)

type ctxPropagatorKey struct{}

// WithValuePropagator returns a new Context that carries value.
func WithValuePropagator(ctx context.Context, p Propagator) context.Context {
	return context.WithValue(ctx, ctxPropagatorKey{}, p)
}

// FromPropagator returns the Propagator value stored in ctx, if any.
func FromPropagator(ctx context.Context) (p Propagator, ok bool) {
	p, ok = ctx.Value(ctxPropagatorKey{}).(Propagator)
	return
}

// MustFromPropagator returns the Propagator value stored in ctx.
func MustFromPropagator(ctx context.Context) Propagator {
	p, ok := ctx.Value(ctxPropagatorKey{}).(Propagator)
	if !ok {
		panic("transport: must be set Propagator into context but it is not!!!")
	}
	return p
}
