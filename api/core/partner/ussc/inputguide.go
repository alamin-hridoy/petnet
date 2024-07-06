package ussc

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
)

var OtherInfoFieldsUSSC = []*ppb.Input{
	{Name: "contact_number"},
	{Name: "purpose_transaction"},
	{Name: "receiver_first_name"},
	{Name: "receiver_last_name"},
	{Name: "receiver_middle_name"},
	{Name: "receiver_name"},
	{Name: "reference_number"},
	{Name: "relation_to"},
	{Name: "sender_first_name"},
	{Name: "sender_last_name"},
	{Name: "sender_middle_name"},
	{Name: "sender_name"},
	{Name: "service_charge"},
	{Name: "total_amount"},
	{Name: "trx_date"},
}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	return &ppb.InputGuideResponse{
		InputGuide: map[string]*ppb.Guide{
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsUSSC,
			},
		},
	}, nil
}
