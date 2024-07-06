package session

import (
	"context"
	"database/sql"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	spb "brank.as/petnet/gunk/v1/session"
	c "brank.as/petnet/profile/core/session"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) SetSessionExpiry(ctx context.Context, req *spb.SetSessionExpiryRequest) (*spb.SetSessionExpiryResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.IDType, validation.Required, validation.Min(0), validation.Max(1)),
		validation.Field(&req.ID, validation.Required, validation.When(req.IDType == spb.IDType_EMAIL, is.Email).Else(is.UUID)),
		validation.Field(&req.Expiry, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := s.core.UpsertSession(ctx, &c.UpsertSessionReq{
		Type: int(req.GetIDType()),
		ID:   req.GetID(),
		Expiry: sql.NullTime{
			Time:  req.GetExpiry().AsTime(),
			Valid: req.GetExpiry().IsValid(),
		},
	}); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &spb.SetSessionExpiryResponse{}, nil
}

func (s *Svc) ExpireSession(ctx context.Context, req *spb.ExpireSessionRequest) (*spb.ExpireSessionResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.IDType, validation.Required, validation.Min(0), validation.Max(1)),
		validation.Field(&req.ID, validation.Required, validation.When(req.IDType == spb.IDType_EMAIL, is.Email).Else(is.UUID)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := s.core.UpsertSession(ctx, &c.UpsertSessionReq{
		Type: int(req.GetIDType()),
		ID:   req.GetID(),
	}); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &spb.ExpireSessionResponse{}, nil
}

func (s *Svc) GetSession(ctx context.Context, req *spb.GetSessionRequest) (*spb.GetSessionResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.IDType, validation.Required, validation.Min(0), validation.Max(1)),
		validation.Field(&req.ID, validation.Required, validation.When(req.IDType == spb.IDType_EMAIL, is.Email).Else(is.UUID)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.core.GetSession(ctx, &c.GetSessionReq{
		Type: int(req.GetIDType()),
		ID:   req.GetID(),
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &spb.GetSessionResponse{
		Expired: res.Expired,
		Expiry:  tspb.New(res.Expiry),
	}, nil
}
