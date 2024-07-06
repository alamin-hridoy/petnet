package mfa

import (
	"context"

	"brank.as/petnet/profile/storage/postgres"
	"brank.as/petnet/serviceutil/logging"
	mpb "brank.as/rbac/gunk/v1/mfa"
)

type EnableMFAReq struct {
	UserID string
	Type   int
	Source string
}

type EnableMFAResp struct {
	ID             string
	InitializeCode string
	EventID        string
	Codes          []string
}

type Svc struct {
	st  *postgres.Storage
	mcl mpb.MFAServiceClient
}

func New(st *postgres.Storage, mcl mpb.MFAServiceClient) *Svc {
	return &Svc{
		st:  st,
		mcl: mcl,
	}
}

func (s *Svc) EnableMFA(ctx context.Context, r EnableMFAReq) (*EnableMFAResp, error) {
	log := logging.FromContext(ctx)
	res, err := s.mcl.EnableMFA(ctx, &mpb.EnableMFARequest{
		UserID: r.UserID,
		Type:   mpb.MFA(r.Type),
		Source: r.Source,
	})
	if err != nil {
		logging.WithError(err, log).Error("signing up")
		return nil, err
	}

	return &EnableMFAResp{
		ID:             res.GetID(),
		InitializeCode: res.GetInitializeCode(),
		EventID:        res.GetEventID(),
		Codes:          res.GetCodes(),
	}, nil
}
