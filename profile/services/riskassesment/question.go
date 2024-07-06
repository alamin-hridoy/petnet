package riskassesment

import (
	"context"

	rat "brank.as/petnet/gunk/dsa/v1/riskassesment"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpsertQuestion(ctx context.Context, req *rat.RiskAssesmentQuestionRequest) (*rat.RiskAssesmentQuestionResponse, error) {
	if err := validation.ValidateStruct(req.Question,
		validation.Field(&req.Question.OrgID, validation.Required, is.UUIDv4),
		validation.Field(&req.Question.UserID, validation.Required, is.UUIDv4),
		validation.Field(&req.Question.QType, validation.Required),
		validation.Field(&req.Question.QID, validation.Required),
		validation.Field(&req.Question.ANS, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.core.UpsertQuestion(ctx, req)
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to store question record")
	}
	return res, nil
}

func (s *Svc) UpsertMlTfQuestion(ctx context.Context, req *rat.RiskAssesmentQuestionRequest) (*rat.RiskAssesmentQuestionResponse, error) {
	if err := validation.ValidateStruct(req.Question,
		validation.Field(&req.Question.OrgID, validation.Required, is.UUIDv4),
		validation.Field(&req.Question.UserID, validation.Required, is.UUIDv4),
		validation.Field(&req.Question.QID, validation.Required),
		validation.Field(&req.Question.CustomersTotal, validation.Required),
		validation.Field(&req.Question.HrTotal, validation.Required),
		validation.Field(&req.Question.ImpactScore, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.core.UpsertMlTfQuestion(ctx, req)
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to store question record")
	}
	return res, nil
}

func (s *Svc) ListQuestion(ctx context.Context, req *rat.ListQuestionRequest) (*rat.ListQuestionResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, validation.Required, is.UUIDv4),
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.ListQuestion(ctx, req)
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to get question record")
	}
	return res, nil
}

func (s *Svc) ListMlTfQuestion(ctx context.Context, req *rat.ListQuestionRequest) (*rat.ListQuestionResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, validation.Required, is.UUIDv4),
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.ListMlTfQuestion(ctx, req)
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to get MLTF question record")
	}
	return res, nil
}
