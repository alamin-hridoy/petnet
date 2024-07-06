package remit

import (
	"context"

	"brank.as/petnet/api/core"
)

func (s *Svc) StageCreateRemit(ctx context.Context, r core.Remittance, partner string) (*core.RemitResponse, error) {
	return s.remitters[partner].StageCreateRemit(ctx, r)
}
