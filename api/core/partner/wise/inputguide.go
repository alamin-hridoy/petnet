package wise

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OtherInfoFieldsWISE ...
// TODO(vitthal): make it dynamic and save in DB after petnet confirmation
var OtherInfoFieldsWISE = []*ppb.Input{}

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
		cntryRes, err := s.ph.WISEgetCountries(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get countries")
			return nil, err
		}
		ip := make([]storage.Input, len(*cntryRes))
		for i, v := range *cntryRes {
			ip[i] = storage.Input{
				Value: v.Key,
				Name:  v.Name,
			}
		}
		ips[storage.IGCountryGroup] = ip

		// get States
		stateRes, err := s.ph.WISEgetStates(ctx, req.CtryCode)
		if err != nil {
			logging.WithError(err, log).Error("get states")
			return nil, err
		}
		ip = make([]storage.Input, len(*stateRes))
		for i, v := range *stateRes {
			ip[i] = storage.Input{
				Value: v.Key,
				Name:  v.Name,
			}
		}
		ips[storage.IGStateGroup+"-"+req.CtryCode] = ip

		// get Currencies
		currRes, err := s.ph.WISEgetCurrencies(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get Currencies")
			return nil, err
		}
		ip = make([]storage.Input, len(*currRes))
		for i, v := range *currRes {
			ip[i] = storage.Input{
				Value:       v.Currency,
				Description: v.Description,
			}
		}
		ips[storage.IGCurrencyGroup] = ip

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

	if _, ok := ig.Data[storage.IGStateGroup+"-"+req.CtryCode]; !ok {
		stateRes, err := s.ph.WISEgetStates(ctx, req.CtryCode)
		if err != nil {
			logging.WithError(err, log).Error("get states")
			return nil, err
		}
		ip := make([]storage.Input, len(*stateRes))
		for i, v := range *stateRes {
			ip[i] = storage.Input{
				Value: v.Key,
				Name:  v.Name,
			}
		}

		ig.Data[storage.IGStateGroup+"-"+req.CtryCode] = ip
		_, err = s.st.UpdateInputGuide(ctx, *ig)
		if err != nil {
			logging.WithError(err, log).Error("create input guide")
			return nil, err
		}
	}

	return &ppb.InputGuideResponse{
		InputGuide: map[string]*ppb.Guide{
			storage.IGCountryLabel: {
				Field:  "receiver.country",
				Inputs: ig.ToGuide(storage.IGCountryGroup),
			},
			storage.IGStateLabel: {
				Field:  "receiver.state",
				Inputs: ig.ToGuide(storage.IGStateGroup + "-" + req.CtryCode),
			},
			storage.IGCurrencyLabel: {
				Field:  "receiver.currency",
				Inputs: ig.ToGuide(storage.IGCurrencyGroup),
			},
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsWISE,
			},
		},
	}, nil
}
