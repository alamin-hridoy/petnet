package fee

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/static"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	fpb "brank.as/petnet/gunk/drp/v1/fee"
	ppb "brank.as/petnet/gunk/dsa/v2/partner"
)

type (
	WUVal   struct{}
	USSCVal struct{}
)

func (*WUVal) Kind() string {
	return static.WUCode
}

func (*USSCVal) Kind() string {
	return static.USSCCode
}

func NewValidators() []Validator {
	return []Validator{&WUVal{}, &USSCVal{}}
}

type Validator interface {
	FeeInquiryValidate(ctx context.Context, st RemcoStore, req *fpb.FeeInquiryRequest) (*core.FeeInquiryReq, error)
	Kind() string
}

type RemcoStore interface {
	ListPartners(context.Context, string) ([]core.Remco, error)
	SendRemitType(ctx context.Context, partner, remTyp string, code bool) (*core.SendRemitType, error)
}

type FeeStore interface {
	FeeInquiry(ctx context.Context, r core.FeeInquiryReq) (map[string]string, error)
}

type Svc struct {
	fpb.UnimplementedFeeServiceServer
	remit      RemcoStore
	validators map[string]Validator
	fee        FeeStore
	cl         ppb.PartnerServiceClient
}

// New Fee service.
func New(remit RemcoStore, fee FeeStore, vs []Validator) *Svc {
	s := &Svc{
		validators: make(map[string]Validator, len(vs)),
		remit:      remit,
		fee:        fee,
	}
	for i, v := range vs {
		switch {
		case v == nil:
			log.Fatalf("validator %d nil", i)
		case v.Kind() == "":
			log.Fatalf("validator %d missing partner type", i)
		}
		s.validators[v.Kind()] = v
	}
	return s
}

// Register the remit service.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	fpb.RegisterFeeServiceServer(srv, s)
	return nil
}

// RegisterGateway parter endpoints.
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	return fpb.RegisterFeeServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
