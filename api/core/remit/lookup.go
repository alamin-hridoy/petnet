package remit

import (
	"context"

	"brank.as/petnet/api/core"
)

func (s *Svc) SearchRemit(ctx context.Context, r core.SearchRemit, partner string) (*core.SearchRemit, error) {
	return s.remitters[partner].Search(ctx, r)
}
