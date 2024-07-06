package remittoaccount

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"brank.as/petnet/api/core/static"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pariz/gountries"

	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
)

type (
	MBRtaVal struct {
		q *gountries.Query
	}
	UBRtaVal struct {
		q *gountries.Query
	}
	BPIRtaVal struct {
		q *gountries.Query
	}
)

func (*MBRtaVal) Kind() string {
	return static.MBRTA
}

func (*UBRtaVal) Kind() string {
	return static.UBRTA
}

func (*BPIRtaVal) Kind() string {
	return static.BPIRTA
}

type Validator interface {
	RTAInquireValidate(context.Context, *rta.RTAInquireRequest) (*rta.RTAInquireRequest, error)
	RTAPaymentValidate(context.Context, *rta.RTAPaymentRequest) (*rta.RTAPaymentRequest, error)
	RTARetryValidate(context.Context, *rta.RTARetryRequest) (*rta.RTARetryRequest, error)
	Kind() string
}

func NewValidators(q *gountries.Query) []Validator {
	return []Validator{&MBRtaVal{}, &UBRtaVal{q}, &BPIRtaVal{q}}
}

type RtaStore interface {
	RTAInquire(context.Context, *rta.RTAInquireRequest, string) (*rta.RTAInquireResponse, error)
	RTAPayment(context.Context, *rta.RTAPaymentRequest, string) (*rta.RTAPaymentResponse, error)
	RTARetry(context.Context, *rta.RTARetryRequest, string) (*rta.RTARetryResponse, error)
}

type Svc struct {
	rta.UnimplementedRemitToAccountServiceServer
	rtaStore   RtaStore
	validators map[string]Validator
}

// New Remit service.
func New(remit RtaStore, vs []Validator) (*Svc, error) {
	s := &Svc{
		validators: make(map[string]Validator, len(vs)),
		rtaStore:   remit,
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
	rta.RegisterRemitToAccountServiceServer(srv, s)
	return nil
}

// RegisterGateway ...
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	return rta.RegisterRemitToAccountServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
