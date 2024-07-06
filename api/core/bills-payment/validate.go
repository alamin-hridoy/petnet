package bills_payment

import (
	"context"

	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
)

func (s *Svc) BPValidate(ctx context.Context, r *bp.BPValidateRequest, partner string) (*bp.BPValidateResponse, error) {
	return s.billers[partner].Validate(ctx, r)
}
