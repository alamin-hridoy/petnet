package middleware

import (
	"context"
	"fmt"

	"brank.as/petnet/serviceutil/errors"
	"brank.as/petnet/serviceutil/logging"

	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
)

// StatusCodeUnaryInterceptor checks that all non-nil errors returned by a grpc
// server have a grpc status set. This will prevent the errors from defaulting
// to the Unknown status or 500 HTTP status codes in the grpc gateway.
//
// In production, the lack of a status code is logged at the error level. In all
// other environments, the lack of the status code results in a panic.
func StatusCodeUnaryInterceptor(env string, logger *logrus.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if _, ok := err.(errors.GRPCStatuser); err != nil && !ok {
			if env != "production" {
				panic(fmt.Sprintf("grpc error is missing GRPCStatus: %+v", err))
			}
			logging.WithError(err, logger).Errorf("grpc error is missing GRPCStatus")
		}
		return resp, err
	}
}

// New creates a gRPC middleware chain
func New(env string, logger *logrus.Entry, logOpts []grpc_logrus.Option, external bool, ints ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	rh := newRecoveryHandler(logger)

	interceptors := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(
			grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		logging.UnaryServerInterceptor(logger, external),
		grpc_logrus.UnaryServerInterceptor(logger, logOpts...),
	}

	if len(ints) > 0 {
		interceptors = append(interceptors, ints...)
	}

	// Recover at the end of the chain, so other interceptors can still process.
	interceptors = append(interceptors,
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(rh.recover)))
	interceptors = append(interceptors, StatusCodeUnaryInterceptor(env, logger))
	return grpc_middleware.ChainUnaryServer(interceptors...)
}

// StatusCodeStreamServerInterceptor checks that all non-nil errors returned by a grpc
// stream server have a grpc status set. This will prevent the errors from defaulting
// to the Unknown status or 500 HTTP status codes in the grpc gateway.
//
// In production, the lack of a status code is logged at the error level. In all
// other environments, the lack of the status code results in a panic.
func StatusCodeStreamServerInterceptor(env string, logger *logrus.Entry) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if _, ok := err.(errors.GRPCStatuser); err != nil && !ok {
			if env != "production" {
				panic(fmt.Sprintf("grpc error is missing GRPCStatus: %+v", err))
			}
			logging.WithError(err, logger).Errorf("grpc error is missing GRPCStatus")
		}
		return err
	}
}

// New creates a gRPC middleware chain for StreamServer
func NewStream(env string, logger *logrus.Entry, logOpts []grpc_logrus.Option, external bool, ints ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	rh := newRecoveryHandler(logger)

	interceptors := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(
			grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
		),
		logging.StreamServerInterceptor(logger, external),
		grpc_logrus.StreamServerInterceptor(logger, logOpts...),
		grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(rh.recover)),
	}

	if len(ints) > 0 {
		interceptors = append(interceptors, ints...)
	}

	interceptors = append(interceptors, StatusCodeStreamServerInterceptor(env, logger))
	return grpc_middleware.ChainStreamServer(interceptors...)
}
