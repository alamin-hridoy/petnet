package wise

import (
	"context"

	"brank.as/petnet/api/integration/perahub"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) DeleteRecipient(ctx context.Context, req *ppb.DeleteRecipientRequest) (*ppb.DeleteRecipientResponse, error) {
	log := logging.FromContext(ctx)

	if _, err := s.ph.WISEDeleteRecipient(ctx, perahub.WISEDeleteRecipientReq{
		RecipientID: req.RecipientID,
		Email:       req.Email,
	}); err != nil {
		logging.WithError(err, log).Error("perahub error")
		return nil, err
	}
	return &ppb.DeleteRecipientResponse{}, nil
}
