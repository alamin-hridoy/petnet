package instacash

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
)

var OtherInfoFieldsIC = []*ppb.Input{
	{Name: "destination_country"},
	{Name: "originating_country"},
	{Name: "purpose"},
	{Name: "receiver_name"},
	{Name: "reference_number"},
	{Name: "sender_name"},
	{Name: "status"},
}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	return &ppb.InputGuideResponse{
		InputGuide: map[string]*ppb.Guide{
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsIC,
			},
		},
	}, nil
}
