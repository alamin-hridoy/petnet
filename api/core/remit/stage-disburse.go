package remit

import (
	"context"

	"brank.as/petnet/api/core"
)

func (s *Svc) StageDisburseRemit(ctx context.Context, r core.Remittance, partner string) (*core.Remittance, error) {
	return s.remitters[partner].StageDisburseRemit(ctx, r)
}
