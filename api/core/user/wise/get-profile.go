package wise

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) GetProfile(ctx context.Context, req core.GetProfileReq) (*core.GetProfileResp, error) {
	log := logging.FromContext(ctx)

	res, err := s.ph.WISEGetProfile(ctx, perahub.WISEGetProfileReq{
		Email: req.Email,
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub error")
		return nil, err
	}

	var occ string
	if len(res.Details.Occupations) > 0 {
		occ = res.Details.Occupations[0].Code
	}

	return &core.GetProfileResp{
		ID:         res.ProfileID.String(),
		Type:       res.Type,
		FirstName:  res.Details.FirstName,
		LastName:   res.Details.LastName,
		BirthDate:  res.Details.BirthDate,
		Phone:      res.Details.PhoneNumber,
		Occupation: occ,
		Address: core.Address{
			Address1:   res.Address.FirstLine,
			City:       res.Address.City,
			State:      res.Address.State,
			PostalCode: res.Address.PostCode,
			Country:    res.Address.Country,
		},
	}, nil
}
