package bills_payment

import (
	"context"

	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
)

func (s *Svc) BPTransact(ctx context.Context, r *bp.BPTransactRequest, partner string) (*bp.BPTransactResponse, error) {
	return s.billers[partner].Transact(ctx, r)
}
