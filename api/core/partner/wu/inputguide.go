package wu

import (
	"context"
	"fmt"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OtherInfoFieldsWU ...
// TODO(vitthal): make it dynamic and save in DB after petnet confirmation
var OtherInfoFieldsWU = []*ppb.Input{}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	log := logging.FromContext(ctx)

	ig, err := s.st.GetInputGuide(ctx, s.Kind())
	if err != nil {
		if err != storage.ErrNotFound {
			logging.WithError(err, log).Error("get input guide from storage")
			return nil, status.Error(codes.Internal, "processing")
		}

		res, err := s.ph.SDQs(ctx, perahub.SDQsRequest{
			Remco:   s.Kind(),
			SDQType: "all",
		})
		if err != nil {
			logging.WithError(err, log).Error("getting sdq ids")
			return nil, err
		}

		ips := make(storage.InputGuideData)

		ip := make([]storage.Input, len(res.ID))
		for i, v := range res.ID {
			ip[i] = storage.Input{
				Value:         v.TemplateValue,
				Name:          v.DocumentType,
				HasIssueDate:  v.HasIssueDate,
				HasExpiration: v.HasExpiration,
				Description:   v.DocDescEng,
			}
		}
		ips[storage.IGIDsGroup] = ip

		oc := make([]storage.Input, len(res.Occupation))
		for i, v := range res.Occupation {
			oc[i] = storage.Input{
				Value: v.OccupationValue,
				Name:  v.Occupation,
			}
		}
		ips[storage.IGOccupationsGroup] = oc

		po := make([]storage.Input, len(res.Position))
		for i, v := range res.Position {
			po[i] = storage.Input{
				Value: v.PositionValue,
				Name:  v.Position,
			}
		}
		ips[storage.IGPositionGroup] = po

		pu := make([]storage.Input, len(res.Purpose))
		for i, v := range res.Purpose {
			pu[i] = storage.Input{
				Value: v.PurposeValue,
				Name:  v.Purpose,
			}
		}
		ips[storage.IGPurposesGroup] = pu

		re := make([]storage.Input, len(res.Relationship))
		for i, v := range res.Relationship {
			re[i] = storage.Input{
				Value: v.RelationshipValue,
				Name:  v.Relationship,
			}
		}
		ips[storage.IGRelationsGroup] = re

		sf := make([]storage.Input, len(res.SourceOfFund))
		for i, v := range res.SourceOfFund {
			sf[i] = storage.Input{
				Value: v.SourceOfFundValue,
				Name:  v.SourceOfFund,
			}
		}
		ips[storage.IGFundsGroup] = sf

		ig = &storage.InputGuide{
			Partner: s.Kind(),
			Data:    ips,
		}
		_, err = s.st.CreateInputGuide(ctx, *ig)
		if err != nil {
			logging.WithError(err, log).Error("create input guide")
			return nil, err
		}
	}

	ccs, err := s.getCountryCurrency(ctx, req.SrcCtry, req.SrcCncy)
	if err != nil {
		return nil, status.Error(codes.Internal, "retrieving input guide")
	}

	return &ppb.InputGuideResponse{
		InputGuide: map[string]*ppb.Guide{
			storage.IGIDsLabel: {
				Field:  "sender.identification.type",
				Inputs: ig.ToGuide(storage.IGIDsGroup),
			},
			"Destination Countries": {
				Field:  "amount.destination_country",
				Inputs: countryToGuide(ccs),
			},
			storage.IGOccupsLabel: {
				Field:  "*.employment.occupation",
				Inputs: ig.ToGuide(storage.IGOccupationsGroup),
			},
			storage.IGPositionLabel: {
				Field:  "*.employment.position_level",
				Inputs: ig.ToGuide(storage.IGPositionGroup),
			},
			storage.IGRelationsLabel: {
				Field:  "*.relationship",
				Inputs: ig.ToGuide(storage.IGRelationsGroup),
			},
			storage.IGPurposesLabel: {
				Field:  "*.transaction_purpose",
				Inputs: ig.ToGuide(storage.IGPurposesGroup),
			},
			storage.IGFundsLabel: {
				Field:  "sender.source_funds",
				Inputs: ig.ToGuide(storage.IGFundsGroup),
			},
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsWU,
			},
		},
	}, nil
}

type CombinedCountryCurrency struct {
	CountryCd  string
	Name       string
	Currencies []Currency
}

type Currency struct {
	CurrencyCd string
	Name       string
}

var CountryCurrencies map[string]CombinedCountryCurrency

// todo(robin): refactor once integrating another partner with country and currency codes
func (s *Svc) getCountryCurrency(ctx context.Context, country, currency string) (map[string]CombinedCountryCurrency, error) {
	log := logging.FromContext(ctx)

	ctyCur := fmt.Sprintf("%s %s", country, currency)
	ccs, err := s.ph.CurrencyCodesRaw(ctx, ctyCur)
	if err != nil {
		logging.WithError(err, log).Error("getting countries and currencies")
		return nil, err
	}

	res := map[string]CombinedCountryCurrency{}
	for _, cc := range ccs {
		c, ok := res[cc.CountryCd]
		if ok {
			c.Currencies = append(c.Currencies, Currency{
				CurrencyCd: cc.CurrencyCd,
				Name:       cc.CurrencyName,
			})
			res[cc.CountryCd] = c
			continue
		}
		res[cc.CountryCd] = CombinedCountryCurrency{
			CountryCd: cc.CountryCd,
			Name:      cc.CountryName,
			Currencies: []Currency{
				{
					CurrencyCd: cc.CurrencyCd,
					Name:       cc.CurrencyName,
				},
			},
		}
	}
	return res, nil
}

func countryToGuide(ccs map[string]CombinedCountryCurrency) []*ppb.Input {
	if ccs == nil {
		return nil
	}

	g := []*ppb.Input{}
	for _, c := range ccs {
		cgs := make([]*ppb.CurrencyGuide, len(c.Currencies))
		for i, cg := range c.Currencies {
			cgs[i] = &ppb.CurrencyGuide{
				Code:         cg.CurrencyCd,
				CurrencyName: cg.Name,
			}
		}
		g = append(g, &ppb.Input{
			Value:       c.CountryCd,
			CountryName: c.Name,
			Currencies:  cgs,
		})
	}
	return g
}
