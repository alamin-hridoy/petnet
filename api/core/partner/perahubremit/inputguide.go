package perahubremit

import (
	"context"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
)

// OtherInfoFieldsUNT ...
// TODO(vitthal): make it dynamic and save in DB after petnet confirmation
var OtherInfoFieldsPerahubRemit = []*ppb.Input{
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
			return nil, coreerror.NewCoreError(codes.Internal, "processing")
		}

		// Get ProvincesCity List
		ips := make(storage.InputGuideData)
		pcRes, err := s.ph.GetProvincesCityList(ctx, perahub.GetProvincesCityListRequest{
			ID:          req.ID,
			PartnerCode: req.Ptnr,
		})
		if err != nil {
			logging.WithError(err, log).Error("get provinces city")
			return nil, coreerror.ToCoreError(err)
		}
		pclCnt := 0
		for _, v := range pcRes.Result {
			pclCnt += len(v.CityList)
		}
		ip := make([]storage.Input, pclCnt)
		cnt := 0
		for _, v := range pcRes.Result {
			for _, cl := range v.CityList {
				ip[cnt] = storage.Input{
					Value:     cl,
					Name:      cl,
					StateName: v.Province,
				}
				cnt++
			}
		}
		ips[storage.IGProvincesCityGroup] = ip
		// Get Brgy List
		blRes, err := s.ph.GetBrgyList(ctx, perahub.GetBrgyListRequest{
			City: req.City,
		})
		if err != nil {
			logging.WithError(err, log).Error("get brgy")
			return nil, coreerror.ToCoreError(err)
		}
		ip = make([]storage.Input, len(blRes.Result))
		for i, v := range blRes.Result {
			ip[i] = storage.Input{
				Value: v.Zipcode,
				Name:  v.Barangay,
			}
		}
		ips[storage.IGBrgyGroup] = ip

		// Get Purpose List
		upRes, err := s.ph.GetUtilityPurpose(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get purpose")
			return nil, coreerror.ToCoreError(err)
		}
		ip = make([]storage.Input, len(upRes.Result))
		for i, v := range upRes.Result {
			ip[i] = storage.Input{
				Value: v,
				Name:  v,
			}
		}
		ips[storage.IGPurposesGroup] = ip

		// Get Utility Relationship
		urRes, err := s.ph.GetUtilityRelationship(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get relationship")
			return nil, coreerror.ToCoreError(err)
		}
		ip = make([]storage.Input, len(urRes.Result))
		for i, v := range urRes.Result {
			ip[i] = storage.Input{
				Value: v,
				Name:  v,
			}
		}
		ips[storage.IGRelationsGroup] = ip

		// Get Utility partner
		uplRes, err := s.ph.GetUtilityPartner(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get partner")
			return nil, coreerror.ToCoreError(err)
		}
		ip = make([]storage.Input, len(uplRes.Result))
		for i, v := range uplRes.Result {
			ip[i] = storage.Input{
				Value: v.PartnerCode,
				Name:  v.PartnerName,
			}
		}
		ips[storage.IGPartnerGroup] = ip

		// Get Utility occupation
		uolRes, err := s.ph.GetUtilityOccupation(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get occupation")
			return nil, coreerror.ToCoreError(err)
		}
		ip = make([]storage.Input, len(uolRes.Result))
		for i, v := range uolRes.Result {
			ip[i] = storage.Input{
				Value: v,
				Name:  v,
			}
		}
		ips[storage.IGOccupationsGroup] = ip

		// Get Utility employement
		uelRes, err := s.ph.GetUtilityEmployment(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get employement")
			return nil, coreerror.ToCoreError(err)
		}
		ip = make([]storage.Input, len(uelRes.Result))
		for i, v := range uelRes.Result {
			ip[i] = storage.Input{
				Value: v,
				Name:  v,
			}
		}
		ips[storage.IGEmploymentGroup] = ip

		// Get Utility sourcefund
		usflRes, err := s.ph.GetUtilitySourceFund(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get sourcefund")
			return nil, coreerror.ToCoreError(err)
		}
		ip = make([]storage.Input, len(usflRes.Result))
		for i, v := range usflRes.Result {
			ip[i] = storage.Input{
				Value: v,
				Name:  v,
			}
		}
		ips[storage.IGFundsGroup] = ip

		ig = &storage.InputGuide{
			Partner: s.Kind(),
			Data:    ips,
		}
		_, err = s.st.CreateInputGuide(ctx, *ig)
		if err != nil {
			logging.WithError(err, log).Error("create input guide")
			return nil, coreerror.ToCoreError(err)
		}
	}

	return &ppb.InputGuideResponse{
		InputGuide: map[string]*ppb.Guide{
			storage.IGProvincesCityLabel: {
				Inputs: ig.ToGuide(storage.IGProvincesCityGroup),
			},
			storage.IGBrgyLabel: {
				Inputs: ig.ToGuide(storage.IGBrgyGroup),
			},
			storage.IGPurposesLabel: {
				Inputs: ig.ToGuide(storage.IGPurposesGroup),
			},
			storage.IGRelationsLabel: {
				Inputs: ig.ToGuide(storage.IGRelationsGroup),
			},
			storage.IGPartnerLabel: {
				Inputs: ig.ToGuide(storage.IGPartnerGroup),
			},
			storage.IGOccupsLabel: {
				Inputs: ig.ToGuide(storage.IGOccupationsGroup),
			},
			storage.IGFundsLabel: {
				Inputs: ig.ToGuide(storage.IGFundsGroup),
			},
			storage.IGEmploymentLabel: {
				Inputs: ig.ToGuide(storage.IGEmploymentGroup),
			},
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsPerahubRemit,
			},
		},
	}, nil
}
