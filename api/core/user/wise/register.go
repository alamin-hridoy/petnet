package wise

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) RegisterUser(ctx context.Context, req core.RegisterUserReq) (*core.RegisterUserResp, error) {
	log := logging.FromContext(ctx)

	if _, err := s.ph.WISECreateUser(ctx, perahub.WISECreateUserReq{
		Email: req.Email,
	}); err != nil {
		logging.WithError(err, log).Error("perahub error")
		return nil, err
	}
	return &core.RegisterUserResp{}, nil
}
