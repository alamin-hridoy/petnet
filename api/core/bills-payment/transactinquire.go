package bills_payment

import (
	"context"

	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
)

func (s *Svc) BPTransactInquire(ctx context.Context, r *bp.BPTransactInquireRequest, partner string) (*bp.BPTransactInquireResponse, error) {
	return s.billers[partner].TransactInquire(ctx, r)
}