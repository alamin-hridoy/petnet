package remittoaccount

import (
	"context"

	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
)

func (s *Svc) RTAPayment(ctx context.Context, r *rta.RTAPaymentRequest, partner string) (*rta.RTAPaymentResponse, error) {
	return s.remitters[partner].Payment(ctx, r)
}
