package ayannah

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
)

var OtherInfoFieldsAYA = []*ppb.Input{
	{Name: "address"},
	{Name: "city"},
	{Name: "contact_number"},
	{Name: "country"},
	{Name: "creation_date"},
	{Name: "destination_country"},
	{Name: "originating_country"},
	{Name: "receiver_name"},
	{Name: "reference_number"},
	{Name: "response_message"},
	{Name: "sender_name"},
	{Name: "zip_code"},
}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	return &ppb.InputGuideResponse{
		InputGuide: map[string]*ppb.Guide{
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsAYA,
			},
		},
	}, nil
}
