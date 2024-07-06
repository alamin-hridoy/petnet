package remit

import (
	"context"

	"brank.as/petnet/api/core"
)

func (s *Svc) ProcessRemit(ctx context.Context, r core.ProcessRemit, partner string) (*core.ProcessRemit, error) {
	return s.remitters[partner].ProcessRemit(ctx, r)
}
