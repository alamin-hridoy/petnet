package token

import (
	"context"

	// grpc
	"google.golang.org/grpc"

	// grpc middleware
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

type AuthFunc func(ctx context.Context, fullMethodName string) (context.Context, error)

// ServiceAuthFuncOverride allows a given gRPC service implementation to override the global `AuthFunc`.
//
// If a service implements the AuthFuncOverride method, it takes precedence over the `AuthFunc` method,
// and will be called instead of AuthFunc for all method invocations within that service.
type ServiceAuthFuncOverride interface {
	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)
}

// UnaryAuthInterceptor returns a new unary server interceptors that performs per-request auth.
func UnaryAuthInterceptor(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
		} else {
			newCtx, err = authFunc(ctx, info.FullMethod)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptors that performs per-request auth.
func StreamServerInterceptor(authFunc AuthFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := srv.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(stream.Context(), info.FullMethod)
		} else {
			newCtx, err = authFunc(stream.Context(), info.FullMethod)
		}
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}

// NoopAuthOverride is a drop-in no-op authentication function.
type NoopAuthOverride struct{}

func (NoopAuthOverride) AuthFuncOverride(ctx context.Context, _ string) (context.Context, error) {
	return ctx, nil
}
