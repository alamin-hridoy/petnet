package wise

import (
	"context"

	"brank.as/petnet/api/integration/perahub"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) CreateRecipient(ctx context.Context, req *ppb.CreateRecipientRequest) (*ppb.CreateRecipientResponse, error) {
	log := logging.FromContext(ctx)

	rs := []perahub.WISERequirement{}
	for _, r := range req.Requirements {
		if r.Value != "" {
			rs = append(rs, perahub.WISERequirement{
				PropName: r.Name,
				Value:    r.Value,
			})
		} else {
			rs = append(rs, perahub.WISERequirement{
				PropName: r.Name,
				Value:    r.Values,
			})
		}
	}

	res, err := s.ph.WISECreateRecipient(ctx, perahub.WISECreateRecipientReq{
		Email:             req.Email,
		Currency:          req.Currency,
		Type:              req.Type,
		OwnedByCustomer:   req.OwnedByCustomer,
		AccountHolderName: req.AccountHolderName,
		Requirements:      rs,
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub error")
		return nil, err
	}
	return &ppb.CreateRecipientResponse{RecipientID: res.RecipientID}, nil
}
