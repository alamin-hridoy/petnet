package cebuana

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OtherInfoFieldsCEB ...
// TODO(vitthal): make it dynamic and save in DB after petnet confirmation
var OtherInfoFieldsCEB = []*ppb.Input{
	{Name: "beneficiary_id"},
	{Name: "birth_date"},
	{Name: "client_reference_no"},
	{Name: "log_id"},
	{Name: "message_id"},
	{Name: "receiver_name"},
	{Name: "remittance_status_description"},
	{Name: "remittance_status_id"},
	{Name: "sender_name"},
	{Name: "service_charge"},
}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	log := logging.FromContext(ctx)

	ig, err := s.st.GetInputGuide(ctx, s.Kind())
	if err != nil {
		if err != storage.ErrNotFound {
			logging.WithError(err, log).Error("get input guide from storage")
			return nil, status.Error(codes.Internal, "processing")
		}
		// get Countries
		ips := make(storage.InputGuideData)
		cntryRes, err := s.ph.CEBgetCountries(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get countries")
			return nil, err
		}
		ip := make([]storage.Input, len(cntryRes.Result.Country))
		for i, v := range cntryRes.Result.Country {
			ip[i] = storage.Input{
				Value:       string(v.CountryID),
				Name:        v.CountryName,
				CountryCode: v.CountryCodeAplha2,
			}
		}
		ips[storage.IGCountryGroup] = ip

		// get Currencies
		stateRes, err := s.ph.CEBgetCurrencies(ctx, &perahub.CEBCurrencyReq{
			AgentCode: req.AgentCode,
		}) // agent code
		if err != nil {
			logging.WithError(err, log).Error("get states")
			return nil, err
		}
		ip = make([]storage.Input, 1)
		ip[0] = storage.Input{
			Value:        string(stateRes.Result.Currency.CurrencyID),
			Description:  stateRes.Result.Currency.Description,
			CurrencyCode: stateRes.Result.Currency.Code,
		}
		ips[storage.IGCurrencyGroup] = ip

		// get Source of fund
		fundsRes, err := s.ph.CEBgetSourceFunds(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get Currencies")
			return nil, err
		}
		ip = make([]storage.Input, len(fundsRes.Result.ClientSourceOfFund))
		for i, v := range fundsRes.Result.ClientSourceOfFund {
			ip[i] = storage.Input{
				Value: string(v.SourceOfFundID),
				Name:  v.SourceOfFund,
			}
		}
		ips[storage.IGFundsGroup] = ip

		// get identification types
		idTypesRes, err := s.ph.CEBgetIdTypes(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get Currencies")
			return nil, err
		}
		ip = make([]storage.Input, len(*idTypesRes))
		for i, v := range *idTypesRes {
			ip[i] = storage.Input{
				Value:       string(v.IdentificationTypeID),
				Name:        v.SmsCode,
				Description: v.Description,
			}
		}
		ips[storage.IGIDsGroup] = ip

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

	return &ppb.InputGuideResponse{
		InputGuide: map[string]*ppb.Guide{
			storage.IGCountryLabel: {
				Field:  "*country",
				Inputs: ig.ToGuide(storage.IGCountryGroup),
			},
			storage.IGCurrencyLabel: {
				Field:  "*currency",
				Inputs: ig.ToGuide(storage.IGCurrencyGroup),
			},
			storage.IGFundsLabel: {
				Field:  "receiver.source_funds,sender.source_funds",
				Inputs: ig.ToGuide(storage.IGFundsGroup),
			},
			storage.IGIDsLabel: {
				Field:  "receiver.identification.type,sender.identification.type",
				Inputs: ig.ToGuide(storage.IGIDsGroup),
			},
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsCEB,
			},
		},
	}, nil
}
