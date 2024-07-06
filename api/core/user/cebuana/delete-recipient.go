package cebuana

import (
	"context"

	ppb "brank.as/petnet/gunk/drp/v1/profile"
)

func (s *Svc) DeleteRecipient(ctx context.Context, req *ppb.DeleteRecipientRequest) (*ppb.DeleteRecipientResponse, error) {
	return &ppb.DeleteRecipientResponse{}, nil
}
