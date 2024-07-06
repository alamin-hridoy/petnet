package microinsurance

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
)

type iMicroInsuranceStore interface {
	Transact(ctx context.Context, req *migunk.TransactRequest) (*migunk.Insurance, error)
	GetReprint(ctx context.Context, req *migunk.GetReprintRequest) (*migunk.Insurance, error)
	RetryTransaction(ctx context.Context, req *migunk.RetryTransactionRequest) (*migunk.Insurance, error)
	GetTransactionList(ctx context.Context, req *migunk.GetTransactionListRequest) (*migunk.TransactionListResult, error)

	GetProduct(context.Context, *migunk.GetProductRequest) (*migunk.ProductResult, error)
	GetOfferProduct(ctx context.Context, req *migunk.GetOfferProductRequest) (*migunk.OfferProduct, error)
	CheckActiveProduct(ctx context.Context, req *migunk.CheckActiveProductRequest) (*migunk.ActiveProduct, error)
	GetProductList(ctx context.Context) (*migunk.GetProductListResult, error)

	GetRelationships(ctx context.Context) (*migunk.GetRelationshipsResult, error)
	GetAllCities(context.Context) (*migunk.CityListResult, error)
}

// Svc ...
type Svc struct {
	migunk.UnimplementedMicroInsuranceServiceServer
	store iMicroInsuranceStore
}

func NewMicroInsuranceSvc(store iMicroInsuranceStore) *Svc {
	return &Svc{
		store: store,
	}
}

// RegisterSvc register the microinsurance service.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	migunk.RegisterMicroInsuranceServiceServer(srv, s)
	return nil
}

// RegisterGateway microinsurance endpoints.
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	return migunk.RegisterMicroInsuranceServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
