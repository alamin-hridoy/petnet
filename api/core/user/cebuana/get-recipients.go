package cebuana

import (
	"context"

	"brank.as/petnet/api/integration/perahub"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) GetRecipients(ctx context.Context, req *ppb.GetRecipientsRequest) (*ppb.GetRecipientsResponse, error) {
	log := logging.FromContext(ctx)

	res, err := s.ph.CebFindBF(ctx, perahub.CebFindBFReq{
		SenderClientId: req.GetSenderUserID(),
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub error")
		return nil, err
	}
	var recipients []*ppb.Recipient
	for _, v := range res.Result.Beneficiary {
		recipients = append(recipients, &ppb.Recipient{
			RecipientID:    v.BeneficiaryID.String(),
			FirstName:      v.FirstName,
			MiddleName:     v.MiddleName,
			LastName:       v.LastName,
			BirthDate:      v.BirthDate,
			StateIDAddress: v.StateIDAddress.String(),
			MobileCountry:  int32(v.CPCountry.CountryID),
			PhoneCountry:   int32(v.TPCountry.CountryID),
			CountryAddress: int32(v.CtryAddress.CountryID),
			BirthCountry:   int32(v.BirthCountry.CountryID),
		})
	}
	return &ppb.GetRecipientsResponse{
		Recipients: recipients,
	}, nil
}
