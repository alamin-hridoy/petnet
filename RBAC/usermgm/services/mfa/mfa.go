package mfa

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/usermgm/core"

	mpb "brank.as/rbac/gunk/v1/mfa"
)

type Svc struct {
	mpb.UnsafeMFAServiceServer
	mpb.UnsafeMFAAuthServiceServer
	a Auth
}

type Auth interface {
	RegisterMFA(context.Context, core.MFA) (*core.MFA, error)
	DisableMFA(context.Context, core.MFA) (*core.MFA, error)
	ListMFA(context.Context, string) ([]core.MFA, error)
	InitiateMFA(context.Context, core.MFAChallenge) (*core.MFAChallenge, error)
	ExternalMFA(context.Context, core.MFAChallenge) (*core.MFAChallenge, error)
	RestartMFA(context.Context, core.MFAChallenge) (*core.MFAChallenge, error)
	MFAuth(context.Context, core.MFAChallenge) (*core.MFAChallenge, error)
}

// New mfa service.
func New(a Auth) *Svc {
	return &Svc{a: a}
}

// RegisterService with grpc server.
func (s *Svc) RegisterSvc(srv *grpc.Server) error { mpb.RegisterMFAServiceServer(srv, s); return nil }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return mpb.RegisterMFAServiceHandlerFromEndpoint(ctx, mux, address, options)
}

// RegisterService with grpc server.
func (s *Svc) RegisterInternal(srv *grpc.Server) error {
	mpb.RegisterMFAAuthServiceServer(srv, s)
	return nil
}

// RegisterGateway grpcgw
func (s *Svc) RegisterGatewayInternal(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return mpb.RegisterMFAAuthServiceHandlerFromEndpoint(ctx, mux, address, options)
}

type resAct struct {
	res, act string
	pub      bool
}

func (s *Svc) Permission(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"GetRegisteredMFA": {res: "ACCOUNT:mfa", act: "view", pub: true},
		"InitiateMFA":      {res: "ACCOUNT:mfa", act: "validate", pub: true},
		"EnableMFA":        {res: "ACCOUNT:mfa", act: "create", pub: true},
		"DisableMFA":       {res: "ACCOUNT:mfa", act: "create", pub: true},
		"ValidateMFA":      {res: "ACCOUNT:mfa", act: "validate", pub: true},
	}
	return p[mthd].res, p[mthd].act, p[mthd].pub
}
