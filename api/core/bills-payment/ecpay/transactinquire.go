package ecpay

import (
	"context"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
)

func (s *Svc) TransactInquire(ctx context.Context, r *bp.BPTransactInquireRequest) (*bp.BPTransactInquireResponse, error) {
	return nil, coreerror.NewCoreError(codes.Unavailable, "Ecpay partner transact inquire is not implemented yet")
}
