package organization

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/usermgm/storage"

	opb "brank.as/rbac/gunk/v1/organization"
)

type Svc struct {
	opb.UnsafeOrganizationServiceServer
	org OrgStore
}

type OrgStore interface {
	GetOrg(ctx context.Context, id string) (*storage.Organization, error)
	UpdateOrg(ctx context.Context, org storage.Organization) (*storage.Organization, error)
}

func New(org OrgStore) *Svc {
	return &Svc{
		org: org,
	}
}

// RegisterService with grpc server.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	opb.RegisterOrganizationServiceServer(srv, s)
	return nil
}

// RegisterGateway grpcgw.
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return opb.RegisterOrganizationServiceHandlerFromEndpoint(ctx, mux, address, options)
}

func (s *Svc) ConfirmUpdate(context.Context, *opb.ConfirmUpdateRequest) (*opb.ConfirmUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmUpdate not implemented")
}

type resAct struct {
	res, act string
	pub      bool
}

func (h *Svc) Permission(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"GetOrganization":    {res: "ACCOUNT:org", act: "view"},
		"ConfirmUpdate":      {res: "ACCOUNT:org", act: "edit", pub: true},
		"UpdateOrganization": {res: "ACCOUNT:org", act: "edit", pub: true},
	}
	return p[mthd].res, p[mthd].act, p[mthd].pub
}
