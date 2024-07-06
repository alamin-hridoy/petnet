package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func NewDetails(c codes.Code, msg string, details ...proto.Message) error {
	s := status.New(c, msg)
	st := s.Proto()
	st.Details = make([]*anypb.Any, len(details))
	for i, d := range details {
		a, err := anypb.New(d)
		if err != nil {
			continue
		}
		st.Details[i] = a
	}
	return status.ErrorProto(st)
}
