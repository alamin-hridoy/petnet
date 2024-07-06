package oauth2

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/usermgm/core"

	oapb "brank.as/rbac/gunk/v1/oauth2"
)

type Svc struct {
	oapb.UnimplementedAuthClientServiceServer
	cl   Store
	envs []interface{}
}

type Store interface {
	CreateClient(context.Context, core.AuthCodeClient) (*core.AuthCodeClient, error)
	GetClient(ctx context.Context, org, client string, listDisable bool) ([]core.AuthCodeClient, error)
	UpdateClient(context.Context, core.AuthCodeClient) (*core.AuthCodeClient, error)
	DeleteClient(context.Context, core.AuthCodeClient) (*core.AuthCodeClient, error)
}

func New(st Store, envs []string) *Svc {
	e := make([]interface{}, len(envs))
	for i, ev := range envs {
		e[i] = ev
	}
	return &Svc{cl: st, envs: e}
}

// RegisterService with grpc server.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	oapb.RegisterAuthClientServiceServer(srv, s)
	return nil
}

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return oapb.RegisterAuthClientServiceHandlerFromEndpoint(ctx, mux, address, options)
}

type resAct struct {
	res, act string
	pub      bool
}

func (s *Svc) Permission(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"CreateClient":  {res: "ACCOUNT:service", act: "create"},
		"UpdateClient":  {res: "ACCOUNT:service", act: "create"},
		"ListClients":   {res: "ACCOUNT:service", act: "view"},
		"DisableClient": {res: "ACCOUNT:service", act: "create"},
	}
	return p[mthd].res, p[mthd].act, p[mthd].pub
}
