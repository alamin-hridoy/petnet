package meta

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/serviceutil/logging"
)

type RequestMeta interface {
	Metadata(context.Context) (context.Context, error)
}

type PublicEndpoint interface {
	PublicEndpoint(method string) bool
}

type PublicService struct{}

func (PublicService) PublicEndpoint(method string) bool { return true }

type Service struct {
	log      logrus.FieldLogger
	metaFunc []RequestMeta
}

func NewMetadata(log logrus.FieldLogger, metaLoaders ...RequestMeta) (*Service, error) {
	s := &Service{
		log:      log,
		metaFunc: make([]RequestMeta, 0, len(metaLoaders)),
	}
	for _, m := range metaLoaders {
		if m != nil {
			s.metaFunc = append(s.metaFunc, m)
		}
	}
	if len(s.metaFunc) == 0 {
		return nil, fmt.Errorf("no metadata loaders")
	}
	return s, nil
}

type MetaFunc func(context.Context) (context.Context, error)

func (m MetaFunc) Metadata(ctx context.Context) (context.Context, error) { return m(ctx) }

// Default secure
func noPub(string) bool { return false }

func (s *Service) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := s.log.WithField("full_method", info.FullMethod)
		svcPub := noPub
		if p, ok := info.Server.(PublicEndpoint); ok {
			svcPub = p.PublicEndpoint
		}

		var finalCtx context.Context
		var err error
		for _, f := range s.metaFunc {
			// exectue metaloader
			finalCtx, err = f.Metadata(logging.WithLogger(ctx, log))
			if err != nil {
				if !svcPub(info.FullMethod) {
					logging.WithError(err, log).Debug("metadata failed")
					// Private endpoints error on first failure
					return nil, status.Error(codes.PermissionDenied, "permission denied")
				}
				finalCtx = ctx
				continue
			}
			// Update context to allow chaining of metaloaders
			ctx = finalCtx
		}
		return handler(finalCtx, req)
	}
}
