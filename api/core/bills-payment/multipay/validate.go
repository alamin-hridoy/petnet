package multipay

import (
	"context"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
)

func (s *Svc) Validate(ctx context.Context, r *bp.BPValidateRequest) (*bp.BPValidateResponse, error) {
	return nil, coreerror.NewCoreError(codes.Unavailable, "Multipay partner validate is not implemented yet")
}
