package remittoaccount

import (
	"context"

	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
)

func (s *Svc) RTAInquire(ctx context.Context, r *rta.RTAInquireRequest, partner string) (*rta.RTAInquireResponse, error) {
	return s.remitters[partner].Inquire(ctx, r)
}
