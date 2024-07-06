package cebuana

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) GetUser(ctx context.Context, req core.GetUserRequest) (*core.GetUserResponse, error) {
	log := logging.FromContext(ctx)

	res, err := s.ph.CebFindClient(ctx, perahub.CebFindClientRequest{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BirthDate:    req.BirthDate,
		ClientNumber: req.ClientNumber,
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub error")
		return nil, err
	}

	return &core.GetUserResponse{
		Code:    res.Code,
		Message: res.Message,
		Result: core.GUResult{
			Client: core.Client{
				ClientID:     res.Result.Client.ClientID,
				ClientNumber: res.Result.Client.ClientNumber,
				FirstName:    res.Result.Client.FirstName,
				MiddleName:   res.Result.Client.MiddleName,
				LastName:     res.Result.Client.LastName,
				BirthDate:    res.Result.Client.BirthDate,
				CPCountry: core.CrtyID{
					CountryID: res.Result.Client.CPCountry.CountryID,
				},
				TPCountry: core.CrtyID{
					CountryID: res.Result.Client.TPCountry.CountryID,
				},
				CtryAddress: core.CrtyID{
					CountryID: res.Result.Client.CtryAddress.CountryID,
				},
				CSOfFund: core.CSOfFund{
					SourceOfFundID: res.Result.Client.CSOfFund.SourceOfFundID,
				},
			},
		},
		RemcoID: res.RemcoID,
	}, nil
}
