package event

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tpb "brank.as/petnet/gunk/dsa/v2/temp"
)

func (s *Svc) GetEventData(ctx context.Context, req *tpb.GetEventDataRequest) (*tpb.GetEventDataResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.EventID, validation.Required, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	e, err := s.core.GetEventData(ctx, req.GetEventID())
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to get event data record")
	}
	return &tpb.GetEventDataResponse{
		EventData: &tpb.EventData{
			EventID:  e.EventID,
			Resource: e.Resource,
			Action:   tpb.ActionType(tpb.ActionType_value[e.Action]),
			Data:     e.Data,
		},
	}, nil
}
