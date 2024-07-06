package transfast

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OtherInfoFieldsTF ...
// TODO(vitthal): make it dynamic and save in DB after petnet confirmation
var OtherInfoFieldsTF = []*ppb.Input{
	{Name: "Desc"},
	{Name: "address"},
	{Name: "contact_number"},
	{Name: "destination_country"},
	{Name: "id_type"},
	{Name: "is_domestic"},
	{Name: "originating_country"},
	{Name: "purpose_of_remittance_id"},
	{Name: "receiver_city_id"},
	{Name: "receiver_city_name"},
	{Name: "receiver_country_iso_code"},
	{Name: "receiver_id_type"},
	{Name: "receiver_is_individual"},
	{Name: "receiver_last_name"},
	{Name: "receiver_name"},
	{Name: "receiver_state_id"},
	{Name: "reference_number"},
	{Name: "sender_name"},
	{Name: "transaction_date"},
}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	log := logging.FromContext(ctx)

	ig, err := s.st.GetInputGuide(ctx, s.Kind())
	if err != nil {
		if err != storage.ErrNotFound {
			logging.WithError(err, log).Error("get input guide from storage")
			return nil, status.Error(codes.Internal, "processing")
		}

		// id types
		ips := make(storage.InputGuideData)
		idRes, err := s.ph.TFIDs(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get ids")
			return nil, err
		}
		ip := make([]storage.Input, len(idRes.Result.IDs))
		for i, v := range idRes.Result.IDs {
			ip[i] = storage.Input{
				Value:         v.ID.String(),
				Name:          v.Name,
				HasIssueDate:  v.RequiredIssueDate,
				HasExpiration: v.RequiredExpirationDate,
				CountryCode:   v.CountryIsoCode,
			}
		}
		ips[storage.IGIDsGroup] = ip

		// relationships
		rlRes, err := s.ph.TFRelations(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get relationships")
			return nil, err
		}
		ip = make([]storage.Input, len(rlRes.Result.Relations))
		for i, v := range rlRes.Result.Relations {
			ip[i] = storage.Input{
				Value: v.ID.String(),
				Name:  v.Name,
			}
		}
		ips[storage.IGRelationsGroup] = ip

		// occupations
		ocRes, err := s.ph.TFOccupations(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get occupations")
			return nil, err
		}
		ip = make([]storage.Input, len(ocRes.Result.Occups))
		for i, v := range ocRes.Result.Occups {
			ip[i] = storage.Input{
				Value: v.ID.String(),
				Name:  v.Name,
			}
		}
		ips[storage.IGOccupationsGroup] = ip

		// purposes
		ppRes, err := s.ph.TFPrps(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get purposes")
			return nil, err
		}
		ip = make([]storage.Input, len(ppRes.Result.Prps))
		for i, v := range ppRes.Result.Prps {
			ip[i] = storage.Input{
				Value:       v.ID.String(),
				Name:        v.Name,
				CountryCode: v.CountryIsoCode,
			}
		}
		ips[storage.IGPurposesGroup] = ip

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
			storage.IGRelationsLabel: {
				Field:  "receiver.relationship",
				Inputs: ig.ToGuide(storage.IGRelationsGroup),
			},
			storage.IGOccupsLabel: {
				Field:  "receiver.employment.occupation_id,receiver.employment.occupation",
				Inputs: ig.ToGuide(storage.IGOccupationsGroup),
			},
			storage.IGPurposesLabel: {
				Field:  "receiver.transaction_purpose",
				Inputs: ig.ToGuide(storage.IGPurposesGroup),
			},
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsTF,
			},
		},
	}, nil
}
