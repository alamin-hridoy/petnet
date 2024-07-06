package partner

import (
	"context"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"google.golang.org/grpc/codes"
)

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	g, ok := s.guiders[req.Ptnr]
	if !ok {
		return nil, coreerror.NewCoreError(codes.NotFound, "missing input guide for partner")
	}
	return g.InputGuide(ctx, req)
}
