package uniteller

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OtherInfoFieldsUNT ...
// TODO(vitthal): make it dynamic and save in DB after petnet confirmation
var OtherInfoFieldsUNT = []*ppb.Input{
	{Name: "address"},
	{Name: "city"},
	{Name: "contact_number"},
	{Name: "country"},
	{Name: "creation_date"},
	{Name: "destination_country"},
	{Name: "formatted_receiver_name"},
	{Name: "formatted_sender_name"},
	{Name: "originating_country"},
	{Name: "receiver_name"},
	{Name: "sender_name"},
	{Name: "zip_code"},
}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	log := logging.FromContext(ctx)

	ig, err := s.st.GetInputGuide(ctx, s.Kind())
	if err != nil {
		if err != storage.ErrNotFound {
			logging.WithError(err, log).Error("get input guide from storage")
			return nil, status.Error(codes.Internal, "processing")
		}

		// countries
		ips := make(storage.InputGuideData)
		ppRes, err := s.ph.UNTgetCountries(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get countries")
			return nil, err
		}
		ip := make([]storage.Input, len(ppRes.Data))
		for i, v := range ppRes.Data {
			ip[i] = storage.Input{
				Value: v.Code,
				Name:  v.Name,
			}
		}
		ips[storage.IGCountryGroup] = ip

		// currency
		gcRes, err := s.ph.UNTgetCurrencies(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get currency")
			return nil, err
		}
		ip = make([]storage.Input, len(gcRes.Data))
		for i, v := range gcRes.Data {
			ip[i] = storage.Input{
				Value: v.Code,
				Name:  v.Name,
			}
		}
		ips[storage.IGCurrencyGroup] = ip

		// id types
		idRes, err := s.ph.UNTgetIds(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get ids")
			return nil, err
		}
		ip = make([]storage.Input, len(idRes.Data))
		for i, v := range idRes.Data {
			ip[i] = storage.Input{
				Value:       v.Code,
				Name:        v.Country,
				Description: v.Description,
			}
		}
		ips[storage.IGIDsGroup] = ip

		// occupations
		ocRes, err := s.ph.UNTgetOccupations(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get occupations")
			return nil, err
		}
		ip = make([]storage.Input, len(ocRes.Data))
		for i, v := range ocRes.Data {
			ip[i] = storage.Input{
				Value:       v.Code,
				Name:        v.Country,
				Description: v.Description,
			}
		}
		ips[storage.IGOccupationsGroup] = ip

		// ph states
		gsRes, err := s.ph.UNTgetStates(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get ph states")
			return nil, err
		}
		ip = make([]storage.Input, len(gsRes.Data))
		for i, v := range gsRes.Data {
			ip[i] = storage.Input{
				Value:       v.UtlCode,
				StateName:   v.StateName,
				CountryName: v.Country,
			}
		}
		ips[storage.IGStateGroup] = ip

		// usa states
		ugsRes, err := s.ph.UNTgetUsStates(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get usa states")
			return nil, err
		}
		ip = make([]storage.Input, len(ugsRes.Data))
		for i, v := range ugsRes.Data {
			ip[i] = storage.Input{
				Value:       v.UtlCode,
				StateName:   v.StateName,
				CountryName: v.Country,
			}
		}
		ips[storage.IGUsStateGroup] = ip

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
			storage.IGIDsLabel: {
				Field:  "receiver.identification.type",
				Inputs: ig.ToGuide(storage.IGIDsGroup),
			},
			storage.IGCountryLabel: {
				Field:  "receiver.country",
				Inputs: ig.ToGuide(storage.IGCountryGroup),
			},
			storage.IGCurrencyLabel: {
				Field:  "receiver.currency",
				Inputs: ig.ToGuide(storage.IGCurrencyGroup),
			},
			storage.IGOccupsLabel: {
				Field:  "receiver.occupation",
				Inputs: ig.ToGuide(storage.IGOccupationsGroup),
			},
			storage.IGStateLabel: {
				Field:  "receiver.states",
				Inputs: ig.ToGuide(storage.IGStateGroup),
			},
			storage.IGUsStateLabel: {
				Field:  "receiver.usa-states",
				Inputs: ig.ToGuide(storage.IGUsStateGroup),
			},
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsUNT,
			},
		},
	}, nil
}
