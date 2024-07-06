package apitransactiontype

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	ppb "brank.as/petnet/gunk/dsa/v2/transactiontype"
)

type Svc struct {
	ppb.UnimplementedTransactionTypeServiceServer
	core ApiTransactionTypeStore
}

type ApiTransactionTypeStore interface {
	CreateApiKeyTransactionType(ctx context.Context, req *ppb.CreateApiKeyTransactionTypeRequest) (*ppb.CreateApiKeyTransactionTypeResponse, error)
	GetAPITransactionType(ctx context.Context, req *ppb.GetAPITransactionTypeRequest) (*ppb.ApiKeyTransactionType, error)
	ListUserAPIKeyTransactionType(ctx context.Context, req *ppb.ListUserAPIKeyTransactionTypeRequest) (*ppb.ListUserAPIKeyTransactionTypeResponse, error)
	GetTransactionTypeByClientId(ctx context.Context, req *ppb.GetTransactionTypeByClientIdRequest) (*ppb.GetTransactionTypeByClientIdResponse, error)
}

func New(core ApiTransactionTypeStore) *Svc {
	h := &Svc{
		core: core,
	}
	return h
}

// RegisterService with grpc server.
func (h *Svc) Register(srv *grpc.Server) { ppb.RegisterTransactionTypeServiceServer(srv, h) }

// RegisterGateway grpcgw
func (h *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return ppb.RegisterTransactionTypeServiceHandlerFromEndpoint(ctx, mux, address, options)
}
