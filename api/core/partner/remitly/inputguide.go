package remitly

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OtherInfoFieldsRM ...
// TODO(vitthal): make it dynamic and save in DB after petnet confirmation
var OtherInfoFieldsRM = []*ppb.Input{
	{Name: "address"},
	{Name: "city"},
	{Name: "contact_number"},
	{Name: "country"},
	{Name: "receiver_name"},
	{Name: "sender_name"},
}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	log := logging.FromContext(ctx)

	ig, err := s.st.GetInputGuide(ctx, s.Kind())
	if err != nil {
		if err != storage.ErrNotFound {
			logging.WithError(err, log).Error("get input guide from storage")
			return nil, status.Error(codes.Internal, "processing")

		}
		res, err := s.ph.RMIDs(ctx)
		if err != nil {
			logging.WithError(err, log).Error("get ids")
			return nil, err
		}

		ip := make([]storage.Input, len(res.Result))
		for i, v := range res.Result {
			ip[i] = storage.Input{
				Value: v.Value,
				Name:  v.Name,
			}
		}

		ig = &storage.InputGuide{
			Partner: s.Kind(),
			Data: storage.InputGuideData{
				storage.IGIDsGroup: ip,
			},
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
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsRM,
			},
		},
	}, nil
}
