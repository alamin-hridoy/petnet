package iremit

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
)

var OtherInfoFieldsIR = []*ppb.Input{
	{Name: "address"},
	{Name: "contact_number"},
	{Name: "desc"},
	{Name: "receiver_first_name"},
	{Name: "receiver_last_name"},
	{Name: "receiver_name"},
	{Name: "reference_number"},
	{Name: "sender_name"},
	{Name: "transaction_date"},
}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	return &ppb.InputGuideResponse{
		InputGuide: map[string]*ppb.Guide{
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsIR,
			},
		},
	}, nil
}
