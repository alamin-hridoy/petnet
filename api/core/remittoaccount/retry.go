package remittoaccount

import (
	"context"

	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
)

func (s *Svc) RTARetry(ctx context.Context, r *rta.RTARetryRequest, partner string) (*rta.RTARetryResponse, error) {
	return s.remitters[partner].Retry(ctx, r)
}
