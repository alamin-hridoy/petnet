package wise

import (
	"context"

	"brank.as/petnet/api/integration/perahub"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) GetRecipients(ctx context.Context, req *ppb.GetRecipientsRequest) (*ppb.GetRecipientsResponse, error) {
	log := logging.FromContext(ctx)

	res, err := s.ph.WISEGetRecipients(ctx, perahub.WISEGetRecipientsReq{
		Email:    req.SenderUserEmail,
		Currency: req.Currency,
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub error")
		return nil, err
	}

	rs := make([]*ppb.Recipient, len(res.Recipients))
	for i, r := range res.Recipients {
		rs[i] = &ppb.Recipient{
			RecipientID: r.RecipientID.String(),
			Details: &ppb.Details{
				AccountNumber:        r.Details.AccountNumber,
				SortCode:             r.Details.SortCode,
				HashByLooseAlgorithm: r.Details.HashedByLooseAlg,
			},
			AccountSummary:     r.AccountSummary,
			LongAccountSummary: r.LongAccountSummary,
			FullName:           r.FullName,
			Currency:           r.Currency,
			Country:            r.Country,
		}

		dfs := make([]*ppb.DisplayField, len(r.DisplayFields))
		for i, df := range r.DisplayFields {
			dfs[i] = &ppb.DisplayField{
				Label: df.Label,
				Value: df.Value,
			}
		}
		rs[i].DisplayFields = dfs
	}

	return &ppb.GetRecipientsResponse{
		Recipients: rs,
	}, nil
}
