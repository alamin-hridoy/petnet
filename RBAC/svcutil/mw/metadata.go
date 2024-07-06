package mw

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
)

type RequestMeta interface {
	Metadata(context.Context) (context.Context, error)
}

type PublicEndpoint interface {
	PublicEndpoint(method string) bool
}

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

func ResetMD(keys ...string) MetaFunc {
	keys = append([]string{OrgIDKey, EnvKey, IDNameKey}, keys...)
	return func(ctx context.Context) (context.Context, error) {
		md := metautils.ExtractIncoming(ctx)
		for _, k := range keys {
			md.Del(k)
		}
		return md.ToIncoming(ctx), nil
	}
}

func ValidateMD(keys ...string) MetaFunc {
	return func(ctx context.Context) (context.Context, error) {
		md := metautils.ExtractIncoming(ctx)
		for _, k := range keys {
			if md.Get(k) == "" {
				return nil, fmt.Errorf("missing %q metadata key", k)
			}
		}
		return ctx, nil
	}
}

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
			// check public endpoints of metaloader
			metaPub := noPub
			if p, ok := f.(PublicEndpoint); ok {
				metaPub = p.PublicEndpoint
			}
			// exectue metaloader
			finalCtx, err = f.Metadata(logging.WithLogger(ctx, log))
			if err != nil {
				logging.WithError(err, log).Debug("metadata failed")
				if !metaPub(info.FullMethod) && !svcPub(info.FullMethod) {
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
