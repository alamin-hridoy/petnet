package fee

import (
	"context"

	"brank.as/petnet/api/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) FeeInquiry(ctx context.Context, r core.FeeInquiryReq) (map[string]string, error) {
	f, ok := s.fee[r.RemitPartner]
	if !ok {
		return nil, status.Error(codes.NotFound, "missing fee inquiry endpoint for partner")
	}
	return f.FeeInquiry(ctx, r)
}
