package quote

import (
	"context"

	phmw "brank.as/petnet/api/perahub-middleware"
	qpb "brank.as/petnet/gunk/drp/v1/quote"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CreateQuote(ctx context.Context, req *qpb.CreateQuoteRequest) (*qpb.CreateQuoteResponse, error) {
	q, ok := s.quoters[phmw.GetPartner(ctx)]
	if !ok {
		return nil, status.Error(codes.NotFound, "missing quoter for partner")
	}
	return q.CreateQuote(ctx, req)
}

func (s *Svc) QuoteInquiry(ctx context.Context, req *qpb.QuoteInquiryRequest) (*qpb.QuoteInquiryResponse, error) {
	q, ok := s.quoters[phmw.GetPartner(ctx)]
	if !ok {
		return nil, status.Error(codes.NotFound, "missing quoter for partner")
	}
	return q.QuoteInquiry(ctx, req)
}

func (s *Svc) QuoteRequirements(ctx context.Context, req *qpb.QuoteRequirementsRequest) (*qpb.QuoteRequirementsResponse, error) {
	q, ok := s.quoters[phmw.GetPartner(ctx)]
	if !ok {
		return nil, status.Error(codes.NotFound, "missing quoter for partner")
	}
	return q.QuoteRequirements(ctx, req)
}
