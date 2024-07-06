package bills_payment

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"brank.as/petnet/api/core/static"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pariz/gountries"

	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
)

type (
	ECPBPVal struct {
		q *gountries.Query
	}
	BYCBPVal struct {
		q *gountries.Query
	}
	MLPBPVal struct {
		q *gountries.Query
	}
)

func (*ECPBPVal) Kind() string {
	return static.ECPBP
}

func (*BYCBPVal) Kind() string {
	return static.BYCBP
}

func (*MLPBPVal) Kind() string {
	return static.MLPBP
}

type Validator interface {
	BPTransactValidate(context.Context, *bp.BPTransactRequest) (*bp.BPTransactRequest, error)
	BPValidateValidate(context.Context, *bp.BPValidateRequest) (*bp.BPValidateRequest, error)
	BPTransactInquireValidate(context.Context, *bp.BPTransactInquireRequest) (*bp.BPTransactInquireRequest, error)
	BPBillerListValidate(context.Context, *bp.BPBillerListRequest) (*bp.BPBillerListRequest, error)
	Kind() string
}

func NewValidators(q *gountries.Query) []Validator {
	return []Validator{&ECPBPVal{}, &BYCBPVal{q}, &MLPBPVal{q}}
}

type BPStore interface {
	BPTransact(context.Context, *bp.BPTransactRequest, string) (*bp.BPTransactResponse, error)
	BPValidate(context.Context, *bp.BPValidateRequest, string) (*bp.BPValidateResponse, error)
	BPTransactInquire(context.Context, *bp.BPTransactInquireRequest, string) (*bp.BPTransactInquireResponse, error)
	BPBillerList(context.Context, *bp.BPBillerListRequest, string) (*bp.BPBillerListResponse, error)
	BillsPaymentTransactList(context.Context, *bp.BillsPaymentTransactListRequest) (*bp.BillsPaymentTransactListResponse, error)
}

type Svc struct {
	bp.UnimplementedBillspaymentServiceServer
	bpStore    BPStore
	validators map[string]Validator
}

// New Remit service.
func New(billspay BPStore, vs []Validator) (*Svc, error) {
	s := &Svc{
		validators: make(map[string]Validator, len(vs)),
		bpStore:    billspay,
	}
	for i, v := range vs {
		switch {
		case v == nil:
			return nil, fmt.Errorf("validator %d nil", i)
		case v.Kind() == "":
			return nil, fmt.Errorf("validator %d missing partner type", i)
		}
		s.validators[v.Kind()] = v
	}
	return s, nil
}

// Register the remit service.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	bp.RegisterBillspaymentServiceServer(srv, s)
	return nil
}

// RegisterGateway ...
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	return bp.RegisterBillspaymentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
