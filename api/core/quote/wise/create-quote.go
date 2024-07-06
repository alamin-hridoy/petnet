package wise

import (
	"context"
	"encoding/json"

	"brank.as/petnet/api/integration/perahub"
	qpb "brank.as/petnet/gunk/drp/v1/quote"
	"brank.as/petnet/serviceutil/logging"
	"github.com/bojanz/currency"
)

func (s *Svc) CreateQuote(ctx context.Context, req *qpb.CreateQuoteRequest) (*qpb.CreateQuoteResponse, error) {
	log := logging.FromContext(ctx)

	amt, err := currency.NewMinor(req.Amount.SourceAmount, req.Amount.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("converting source amount")
		return nil, err
	}

	res, err := s.ph.WISECreateQuote(ctx, perahub.WISECreateQuoteReq{
		Email:          req.Email,
		SourceCurrency: req.Amount.SourceCurrency,
		TargetCurrency: req.Amount.DestinationCurrency,
		SourceAmount:   json.Number(amt.Amount.Number()),
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub error")
		return nil, err
	}

	rs := make([]*qpb.Requirement, len(res.Requirements))
	for i, r := range res.Requirements {
		rs[i] = &qpb.Requirement{
			Type:      r.Type,
			Title:     r.Title,
			UsageInfo: r.UsageInfo,
		}

		fs := make([]*qpb.Field, len(r.Fields))
		for i, f := range r.Fields {
			fs[i] = &qpb.Field{
				Name: f.Name,
			}

			gs := make([]*qpb.Group, len(f.Group))
			for i, g := range f.Group {
				gs[i] = &qpb.Group{
					Key:                         g.Key,
					Name:                        g.Name,
					Type:                        g.Type,
					RefreshRequirementsOnChange: g.RefreshReqOnChange,
					Required:                    g.Required,
					Example:                     g.Example,
					DisplayFormat:               g.DisplayFormat,
					MinLength:                   string(g.MinLength),
					MaxLength:                   string(g.MaxLength),
					ValidationRegexp:            g.ValidationRegexp,
				}

				vs := make([]*qpb.Param, len(g.ValidationAsync.Params))
				for i, v := range g.ValidationAsync.Params {
					vs[i] = &qpb.Param{
						Key:           v.Key,
						ParamaterName: v.ParamName,
						Required:      v.Required,
					}
				}

				gs[i].ValidationAsync = &qpb.ValidationAsync{
					URL:    g.ValidationAsync.URL,
					Params: vs,
				}

				vas := make([]*qpb.ValueAllowed, len(g.ValuesAllowed))
				for i, v := range g.ValuesAllowed {
					vas[i] = &qpb.ValueAllowed{
						Key:  v.Key,
						Name: v.Name,
					}
				}
				gs[i].ValuesAllowed = vas
			}
			fs[i].Groups = gs
		}
		rs[i].Fields = fs
	}

	samt, err := currency.NewAmount(string(res.QuoteSummary.SourceAmount), res.QuoteSummary.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating source amount")
		return nil, err
	}
	damt, err := currency.NewAmount(string(res.QuoteSummary.TargetAmount), res.QuoteSummary.TargetCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating destination amount")
		return nil, err
	}
	twamt, err := currency.NewAmount(res.QuoteSummary.FeeBreakdown.Transferwise.String(), res.QuoteSummary.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating transferwise amount")
		return nil, err
	}
	piamt, err := currency.NewAmount(res.QuoteSummary.FeeBreakdown.PayIn.String(), res.QuoteSummary.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating payin amount")
		return nil, err
	}
	dcamt, err := currency.NewAmount(res.QuoteSummary.FeeBreakdown.Discount.String(), res.QuoteSummary.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating discount amount")
		return nil, err
	}
	tamt, err := currency.NewAmount(res.QuoteSummary.FeeBreakdown.Total.String(), res.QuoteSummary.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating total amount")
		return nil, err
	}
	tfamt, err := currency.NewAmount(res.QuoteSummary.TotalFee.String(), res.QuoteSummary.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating total fee")
		return nil, err
	}
	trfamt, err := currency.NewAmount(res.QuoteSummary.TransferAmount.String(), res.QuoteSummary.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating transer amount")
		return nil, err
	}

	qs := &qpb.QuoteSummary{
		SourceCurrency:      res.QuoteSummary.SourceCurrency,
		DestinationCurrency: res.QuoteSummary.TargetCurrency,
		SourceAmount:        currency.ToMinor(samt.Round()).Number(),
		DestinationAmount:   currency.ToMinor(damt.Round()).Number(),
		FeeBreakdown: map[string]string{
			"transferwise": currency.ToMinor(twamt.Round()).Number(),
			"payIn":        currency.ToMinor(piamt.Round()).Number(),
			"discount":     currency.ToMinor(dcamt.Round()).Number(),
			"total":        currency.ToMinor(tamt.Round()).Number(),
			"priceSetId":   string(res.QuoteSummary.FeeBreakdown.PriceSetID),
			"partner":      string(res.QuoteSummary.FeeBreakdown.Partner),
		},
		TotalFee:       currency.ToMinor(tfamt.Round()).Number(),
		TransferAmount: currency.ToMinor(trfamt.Round()).Number(),
		Payout:         res.QuoteSummary.PayOut,
		Rate:           string(res.QuoteSummary.Rate),
	}

	return &qpb.CreateQuoteResponse{
		Requirements: rs,
		QuoteSummary: qs,
	}, nil
}
