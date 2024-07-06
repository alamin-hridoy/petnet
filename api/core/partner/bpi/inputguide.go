package bpi

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
)

var OtherInfoFieldsBPI = []*ppb.Input{
	{Name: "Desc"},
	{Name: "address"},
	{Name: "client_reference_no"},
	{Name: "contact_number"},
	{Name: "destination_country"},
	{Name: "originating_country"},
	{Name: "receiver_name"},
	{Name: "reference_number"},
	{Name: "sender_name"},
}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	return &ppb.InputGuideResponse{
		InputGuide: map[string]*ppb.Guide{
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsBPI,
			},
		},
	}, nil
}
