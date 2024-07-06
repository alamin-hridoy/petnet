package cebuanaint

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
)

var OtherInfoFieldsCEBINT = []*ppb.Input{
	{Name: "beneficiary_id"},
	{Name: "birth_date"},
	{Name: "client_reference_no"},
	{Name: "is_domestic"},
	{Name: "log_id"},
	{Name: "message_id"},
	{Name: "receiver_name"},
	{Name: "remittance_status_description"},
	{Name: "remittance_status_id"},
	{Name: "sender_name"},
	{Name: "service_charge"},
}

func (s *Svc) InputGuide(ctx context.Context, req core.InputGuideRequest) (*ppb.InputGuideResponse, error) {
	return &ppb.InputGuideResponse{
		InputGuide: map[string]*ppb.Guide{
			storage.IGOtherInfoLabel: {
				OtherInfo: OtherInfoFieldsCEBINT,
			},
		},
	}, nil
}
