package wise

import (
	"context"
	"encoding/json"

	"brank.as/petnet/api/integration/perahub"
	qpb "brank.as/petnet/gunk/drp/v1/quote"
	"brank.as/petnet/serviceutil/logging"
	"github.com/bojanz/currency"
)

func (s *Svc) QuoteInquiry(ctx context.Context, req *qpb.QuoteInquiryRequest) (*qpb.QuoteInquiryResponse, error) {
	log := logging.FromContext(ctx)

	amt, err := currency.NewMinor(req.Amount.SourceAmount, req.Amount.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("converting source amount")
		return nil, err
	}

	res, err := s.ph.WISEQuoteInquiry(ctx, perahub.WISEQuoteInquiryReq{
		SourceCurrency: req.Amount.SourceCurrency,
		TargetCurrency: req.Amount.DestinationCurrency,
		SourceAmount:   json.Number(amt.Amount.Number()),
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub error")
		return nil, err
	}

	samt, err := currency.NewAmount(string(res.SourceAmount), res.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating source amount")
		return nil, err
	}
	twamt, err := currency.NewAmount(res.FeeBreakdown.Transferwise.String(), res.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating transferwise amount")
		return nil, err
	}
	piamt, err := currency.NewAmount(res.FeeBreakdown.PayIn.String(), res.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating payin amount")
		return nil, err
	}
	damt, err := currency.NewAmount(res.FeeBreakdown.Discount.String(), res.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating discount amount")
		return nil, err
	}
	tamt, err := currency.NewAmount(res.FeeBreakdown.Total.String(), res.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating total amount")
		return nil, err
	}
	tfamt, err := currency.NewAmount(res.TotalFee.String(), res.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating total fee")
		return nil, err
	}
	trfamt, err := currency.NewAmount(res.TransferAmount.String(), res.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating transer amount")
		return nil, err
	}

	return &qpb.QuoteInquiryResponse{
		Amount: &qpb.QuoteAmount{
			SourceAmount:        currency.ToMinor(samt.Round()).Number(),
			SourceCurrency:      res.SourceCurrency,
			DestinationCurrency: res.TargetCurrency,
		},
		FeeBreakdown: map[string]string{
			"transferwise": currency.ToMinor(twamt.Round()).Number(),
			"payIn":        currency.ToMinor(piamt.Round()).Number(),
			"discount":     currency.ToMinor(damt.Round()).Number(),
			"total":        currency.ToMinor(tamt.Round()).Number(),
			"priceSetId":   string(res.FeeBreakdown.PriceSetID),
			"partner":      string(res.FeeBreakdown.Partner),
		},
		TotalFee:       currency.ToMinor(tfamt.Round()).Number(),
		TransferAmount: currency.ToMinor(trfamt.Round()).Number(),
	}, nil
}
