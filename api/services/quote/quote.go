package quote

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"brank.as/petnet/api/core/static"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	qpb "brank.as/petnet/gunk/drp/v1/quote"
)

type (
	WISEVal struct{}
)

func (*WISEVal) Kind() string {
	return static.WISECode
}

type Validator interface {
	CreateQuoteValidate(ctx context.Context, req *qpb.CreateQuoteRequest) error
	QuoteInquiryValidate(ctx context.Context, req *qpb.QuoteInquiryRequest) error
	QuoteRequirementsValidate(ctx context.Context, req *qpb.QuoteRequirementsRequest) error
	Kind() string
}

func NewValidators() []Validator {
	return []Validator{&WISEVal{}}
}

type QuoteStore interface {
	CreateQuote(ctx context.Context, req *qpb.CreateQuoteRequest) (*qpb.CreateQuoteResponse, error)
	QuoteInquiry(ctx context.Context, req *qpb.QuoteInquiryRequest) (*qpb.QuoteInquiryResponse, error)
	QuoteRequirements(ctx context.Context, req *qpb.QuoteRequirementsRequest) (*qpb.QuoteRequirementsResponse, error)
}

type Svc struct {
	qpb.UnimplementedQuoteServiceServer
	quote      QuoteStore
	validators map[string]Validator
}

func New(st QuoteStore, vs []Validator) *Svc {
	s := &Svc{
		validators: make(map[string]Validator, len(vs)),
		quote:      st,
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
	qpb.RegisterQuoteServiceServer(srv, s)
	return nil
}

// RegisterGateway parter endpoints.
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	return qpb.RegisterQuoteServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
