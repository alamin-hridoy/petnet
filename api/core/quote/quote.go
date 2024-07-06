package quote

import (
	"context"
	"log"

	wum "brank.as/petnet/api/core/quote/wise"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage/postgres"
	qpb "brank.as/petnet/gunk/drp/v1/quote"
)

type Quoter interface {
	CreateQuote(ctx context.Context, req *qpb.CreateQuoteRequest) (*qpb.CreateQuoteResponse, error)
	QuoteInquiry(ctx context.Context, req *qpb.QuoteInquiryRequest) (*qpb.QuoteInquiryResponse, error)
	QuoteRequirements(ctx context.Context, req *qpb.QuoteRequirementsRequest) (*qpb.QuoteRequirementsResponse, error)
	Kind() string
}

type Svc struct {
	quoters map[string]Quoter
	st      *postgres.Storage
	ph      *perahub.Svc
}

func New(st *postgres.Storage, ph *perahub.Svc) *Svc {
	gs := []Quoter{wum.New(ph)}
	s := &Svc{
		quoters: make(map[string]Quoter, len(gs)),
		st:      st,
		ph:      ph,
	}
	for i, r := range gs {
		switch {
		case r == nil:
			log.Fatalf("quoter %d nil", i)
		case r.Kind() == "":
			log.Fatalf("quoter %d missing partner type", i)
		}
		s.quoters[r.Kind()] = r
	}
	return s
}
