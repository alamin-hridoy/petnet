package mfa

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mpb "brank.as/petnet/gunk/v1/mfa"
	core "brank.as/petnet/profile/core/mfa"
)

func (s *Svc) EnableMFA(ctx context.Context, req *mpb.EnableMFARequest) (*mpb.EnableMFAResponse, error) {
	res, err := s.core.EnableMFA(ctx, core.EnableMFAReq{
		UserID: req.GetUserID(),
		Type:   int(req.GetType()),
		Source: req.GetSource(),
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to enable mfa")
	}
	return &mpb.EnableMFAResponse{
		ID:             res.ID,
		InitializeCode: res.InitializeCode,
		EventID:        res.EventID,
		Codes:          res.Codes,
	}, nil
}
