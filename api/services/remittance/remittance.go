package remittance

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

type RemittanceCoreStore interface {
	ValidateSendMoney(ctx context.Context, req *bpa.ValidateSendMoneyRequest) (*bpa.ValidateSendMoneyResponse, error)
	ConfirmSendMoney(ctx context.Context, req *bpa.ConfirmSendMoneyRequest) (*bpa.ConfirmSendMoneyResponse, error)
	CancelSendMoney(ctx context.Context, req *bpa.CancelSendMoneyRequest) (*bpa.CancelSendMoneyResponse, error)
	ValidateReceiveMoney(ctx context.Context, req *bpa.ValidateReceiveMoneyRequest) (*bpa.ValidateReceiveMoneyResponse, error)
	Inquire(ctx context.Context, req *bpa.InquireRequest) (*bpa.InquireResponse, error)
	ConfirmReceiveMoney(ctx context.Context, req *bpa.ConfirmReceiveMoneyRequest) (*bpa.ConfirmReceiveMoneyResponse, error)
	PartnersGrid(ctx context.Context) (*bpa.PartnersGridResponse, error)
	PartnersCreate(ctx context.Context, req *bpa.PartnersCreateRequest) (*bpa.PartnersCreateResponse, error)
	PurposeOfRemittanceGrid(ctx context.Context) (*bpa.PurposeOfRemittanceGridResponse, error)
	PurposeOfRemittanceGet(ctx context.Context, req *bpa.PurposeOfRemittanceGetRequest) (*bpa.PurposeOfRemittanceGetResponse, error)
	PurposeOfRemittanceUpdate(ctx context.Context, req *bpa.PurposeOfRemittanceUpdateRequest) (*bpa.PurposeOfRemittanceUpdateResponse, error)
	PurposeOfRemittanceCreate(ctx context.Context, req *bpa.PurposeOfRemittanceCreateRequest) (*bpa.PurposeOfRemittanceCreateResponse, error)
	SourceOfFundGrid(ctx context.Context) (*bpa.SourceOfFundGridResponse, error)
	SourceOfFundCreate(ctx context.Context, req *bpa.SourceOfFundCreateRequest) (*bpa.SourceOfFundCreateResponse, error)
	SourceOfFundGet(ctx context.Context, req *bpa.SourceOfFundGetRequest) (*bpa.SourceOfFundGetResponse, error)
	RemittanceEmploymentCreate(ctx context.Context, req *bpa.RemittanceEmploymentCreateRequest) (*bpa.RemittanceEmploymentCreateResponse, error)
	EmploymentGrid(ctx context.Context) (*bpa.EmploymentGridResponse, error)
	OccupationGrid(ctx context.Context) (*bpa.OccupationGridResponse, error)
	OccupationGet(ctx context.Context, req *bpa.OccupationGetRequest) (*bpa.OccupationGetResponse, error)
	OccupationCreate(ctx context.Context, req *bpa.OccupationCreateRequest) (*bpa.OccupationCreateResponse, error)
	OccupationUpdate(ctx context.Context, req *bpa.OccupationUpdateRequest) (*bpa.OccupationUpdateResponse, error)
	OccupationDelete(ctx context.Context, req *bpa.OccupationDeleteRequest) (*bpa.OccupationDeleteResponse, error)
	RelationshipGet(ctx context.Context, req *bpa.RelationshipGetRequest) (*bpa.RelationshipGetResponse, error)
	SourceOfFundUpdate(ctx context.Context, req *bpa.SourceOfFundUpdateRequest) (*bpa.SourceOfFundUpdateResponse, error)
	SourceOfFundDelete(ctx context.Context, req *bpa.SourceOfFundDeleteRequest) (*bpa.SourceOfFundDeleteResponse, error)
	PurposeOfRemittanceDelete(ctx context.Context, req *bpa.PurposeOfRemittanceDeleteRequest) (*bpa.PurposeOfRemittanceDeleteResponse, error)
	RelationshipDelete(ctx context.Context, req *bpa.RelationshipDeleteRequest) (*bpa.RelationshipDeleteResponse, error)
	EmploymentGet(ctx context.Context, req *bpa.EmploymentGetRequest) (*bpa.EmploymentGetResponse, error)
	RemittanceEmploymentUpdate(ctx context.Context, req *bpa.RemittanceEmploymentUpdateRequest) (*bpa.RemittanceEmploymentUpdateResponse, error)
	RemittanceEmploymentDelete(ctx context.Context, req *bpa.RemittanceEmploymentDeleteRequest) (*bpa.RemittanceEmploymentDeleteResponse, error)
	RelationshipGrid(ctx context.Context, em *emptypb.Empty) (*bpa.RelationshipGridResponse, error)
	RelationshipUpdate(ctx context.Context, req *bpa.RelationshipUpdateRequest) (*bpa.RelationshipUpdateResponse, error)
	PartnersDelete(ctx context.Context, req *bpa.PartnersDeleteRequest) (*bpa.PartnersDeleteResponse, error)
	PartnersGet(ctx context.Context, req *bpa.PartnersGetRequest) (*bpa.PartnersGetResponse, error)
	PartnersUpdate(ctx context.Context, req *bpa.PartnersUpdateRequest) (*bpa.PartnersUpdateResponse, error)
	RelationshipCreate(ctx context.Context, req *bpa.RelationshipCreateRequest) (*bpa.RelationshipCreateResponse, error)
}

type Svc struct {
	bpa.UnimplementedRemittanceServiceServer
	remittanceStore RemittanceCoreStore
}

func New(st RemittanceCoreStore) *Svc {
	s := &Svc{
		remittanceStore: st,
	}
	return s
}

// RegisterSvc the remit service.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	bpa.RegisterRemittanceServiceServer(srv, s)
	return nil
}

// RegisterGateway parter endpoints.
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	return bpa.RegisterRemittanceServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
