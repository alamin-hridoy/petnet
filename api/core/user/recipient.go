package user

import (
	"context"

	phmw "brank.as/petnet/api/perahub-middleware"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CreateRecipient(ctx context.Context, req *ppb.CreateRecipientRequest) (*ppb.CreateRecipientResponse, error) {
	um, ok := s.usermanagers[phmw.GetPartner(ctx)]
	if !ok {
		return nil, status.Error(codes.NotFound, "missing user management for partner")
	}
	return um.CreateRecipient(ctx, req)
}

func (s *Svc) GetRecipients(ctx context.Context, req *ppb.GetRecipientsRequest) (*ppb.GetRecipientsResponse, error) {
	um, ok := s.usermanagers[phmw.GetPartner(ctx)]
	if !ok {
		return nil, status.Error(codes.NotFound, "missing user management for partner")
	}
	return um.GetRecipients(ctx, req)
}

func (s *Svc) RefreshRecipient(ctx context.Context, req *ppb.UpdateRecipientRequest) (*ppb.UpdateRecipientResponse, error) {
	um, ok := s.usermanagers[phmw.GetPartner(ctx)]
	if !ok {
		return nil, status.Error(codes.NotFound, "missing user management for partner")
	}
	return um.RefreshRecipient(ctx, req)
}

func (s *Svc) DeleteRecipient(ctx context.Context, req *ppb.DeleteRecipientRequest) (*ppb.DeleteRecipientResponse, error) {
	um, ok := s.usermanagers[phmw.GetPartner(ctx)]
	if !ok {
		return nil, status.Error(codes.NotFound, "missing user management for partner")
	}
	return um.DeleteRecipient(ctx, req)
}
