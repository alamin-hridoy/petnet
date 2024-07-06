package wise

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) CreateProfile(ctx context.Context, req core.CreateProfileReq) (*core.CreateProfileResp, error) {
	log := logging.FromContext(ctx)

	if _, err := s.ph.WISECreateProfile(ctx, perahub.WISECreateProfileReq{
		Email: req.Email,
		Type:  req.Type,
		Details: perahub.WISECreatePFDetails{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			BirthDate:   req.BirthDate,
			PhoneNumber: "+" + req.Phone.CtyCode + req.Phone.Number,
		},
		Address: perahub.WISEPFAddress{
			Country:   req.Address.Country,
			FirstLine: req.Address.Address1,
			PostCode:  req.Address.PostalCode,
			City:      req.Address.City,
			State:     req.Address.State,
		},
		Occupation: req.Occupation,
	}); err != nil {
		logging.WithError(err, log).Error("perahub error")
		return nil, err
	}
	return &core.CreateProfileResp{}, nil
}
