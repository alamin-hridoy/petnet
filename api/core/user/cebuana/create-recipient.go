package cebuana

import (
	"context"
	"strconv"

	"brank.as/petnet/api/integration/perahub"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) CreateRecipient(ctx context.Context, req *ppb.CreateRecipientRequest) (*ppb.CreateRecipientResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.ph.CebAddBf(ctx, perahub.CebAddBftReq{
		FirstName:          req.GetFirstName(),
		MiddleName:         req.GetMiddleName(),
		LastName:           req.GetLastName(),
		SenderClientID:     int(req.GetSenderUserID()),
		BirthDate:          req.GetBirthDate(),
		CellphoneCountryID: req.GetPhoneCountryID(),
		ContactNumber:      req.GetContactNumber(),
		TelephoneCountryID: req.GetPhoneCountryID(),
		TelephoneAreaCode:  req.GetPhoneAreaCode(),
		TelephoneNumber:    req.GetPhoneNumber(),
		CountryAddressID:   req.GetCountryAddressID(),
		BirthCountryID:     req.GetBirthCountryID(),
		ProvinceAddress:    req.GetProvinceAddress(),
		Address:            req.GetAddress(),
		UserID:             int(req.GetUserID()),
		Occupation:         req.Occupation,
		ZipCode:            req.GetPostalCode(),
		StateIDAddress:     req.GetStateIDAddress(),
		Tin:                req.GetTin(),
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub error")
		return nil, err
	}
	return &ppb.CreateRecipientResponse{
		RecipientID: strconv.Itoa(res.Result.BeneficiaryID),
	}, nil
}
