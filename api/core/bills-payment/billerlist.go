package bills_payment

import (
	"context"

	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
)

func (s *Svc) BPBillerList(ctx context.Context, r *bp.BPBillerListRequest, partner string) (*bp.BPBillerListResponse, error) {
	return s.billers[partner].BillerList(ctx, r)
}
