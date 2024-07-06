package ria

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
)

var OtherInfoFieldsRIA = []*ppb.Input{
	{Name: "client_reference_no"},
	{Name: "destination_country"},
	{Name: "is_domestic"},
	{Name: "order_number"},
	{Name: "originating_country"},
	{Name: "receiver_name"},
	{Name: "sender_name"},
}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	return &ppb.InputGuideResponse{
		InputGuide: map[string]*ppb.Guide{
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsRIA,
			},
		},
	}, nil
}
