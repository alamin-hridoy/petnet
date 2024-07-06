package japanremit

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
)

var OtherInfoFieldsJPR = []*ppb.Input{
	{Name: "currency"},
	{Name: "destination_country"},
	{Name: "originating_country"},
	{Name: "pay_token_id"},
	{Name: "receiver_name"},
	{Name: "reference_number"},
	{Name: "sender_name"},
}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	return &ppb.InputGuideResponse{
		InputGuide: map[string]*ppb.Guide{
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsJPR,
			},
		},
	}, nil
}
