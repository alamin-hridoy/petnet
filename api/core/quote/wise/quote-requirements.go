package wise

import (
	"context"
	"encoding/json"

	"brank.as/petnet/api/integration/perahub"
	qpb "brank.as/petnet/gunk/drp/v1/quote"
	"brank.as/petnet/serviceutil/logging"
	"github.com/bojanz/currency"
)

func (s *Svc) QuoteRequirements(ctx context.Context, req *qpb.QuoteRequirementsRequest) (*qpb.QuoteRequirementsResponse, error) {
	log := logging.FromContext(ctx)

	amt, err := currency.NewMinor(req.Amount.SourceAmount, req.Amount.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("converting source amount")
		return nil, err
	}

	res, err := s.ph.WISEGetQuoteRequirements(ctx, perahub.WISEGetQuoteRequirementsReq{
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

	samt, err := currency.NewAmount(res.Quote.SourceAmount.String(), res.Quote.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("creating source amount")
		return nil, err
	}

	q := &qpb.Quote{
		SourceCurrency:      res.Quote.SourceCurrency,
		DestinationCurrency: res.Quote.TargetCurrency,
		SourceAmount:        currency.ToMinor(samt.Round()).Number(),
	}
	return &qpb.QuoteRequirementsResponse{
		Requirements: rs,
		Quote:        q,
	}, nil
}
