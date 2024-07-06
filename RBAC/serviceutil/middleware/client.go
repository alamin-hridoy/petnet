package middleware

import (
	"brank.as/rbac/serviceutil/logging"

	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
)

// New creates a gRPC middleware chain
func NewClient(env string, logger *logrus.Entry, ignoreFields []string, logOpts []grpc_logrus.Option, ints ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	interceptors := []grpc.UnaryClientInterceptor{
		logging.UnaryClientInterceptor(ignoreFields...),
		grpc_logrus.UnaryClientInterceptor(logger, logOpts...),
	}

	if len(ints) > 0 {
		interceptors = append(interceptors, ints...)
	}

	return grpc_middleware.ChainUnaryClient(interceptors...)
}
