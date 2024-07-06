package event

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) CreateEventData(ctx context.Context, req *tpb.CreateEventDataRequest) (*tpb.CreateEventDataResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.EventData, validation.Required, validation.By(func(interface{}) error {
			e := req.GetEventData()
			if err := validation.ValidateStruct(e,
				validation.Field(&e.EventID, validation.Required, is.UUIDv4),
				validation.Field(&e.Resource, validation.Required),
				validation.Field(&e.Action, validation.Required),
				validation.Field(&e.Data, validation.Required, is.JSON),
			); err != nil {
				return status.Error(codes.InvalidArgument, err.Error())
			}
			return nil
		})),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	e := req.GetEventData()
	if err := s.core.CreateEventData(ctx, &storage.EventData{
		EventID:  e.EventID,
		Resource: e.Resource,
		Action:   e.Action.String(),
		Data:     e.Data,
	}); err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to store event data record")
	}
	return &tpb.CreateEventDataResponse{}, nil
}
