package riskassesment

import (
	"context"

	rat "brank.as/petnet/gunk/dsa/v1/riskassesment"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) UpsertQuestion(ctx context.Context, req *rat.RiskAssesmentQuestionRequest) (*rat.RiskAssesmentQuestionResponse, error) {
	log := logging.FromContext(ctx)
	qes, err := s.st.UpsertQuestion(ctx, &storage.Question{
		ID:      req.GetQuestion().GetID(),
		OrgID:   req.GetQuestion().GetOrgID(),
		UserID:  req.GetQuestion().GetUserID(),
		QID:     req.GetQuestion().GetQID(),
		ANS:     req.GetQuestion().GetANS(),
		QType:   req.GetQuestion().GetQType(),
		Created: req.GetQuestion().GetCreated().AsTime(),
		Updated: req.GetQuestion().GetUpdated().AsTime(),
	})
	if err != nil {
		logging.WithError(err, log).Error("Upsert Question error")
		return nil, err
	}
	return &rat.RiskAssesmentQuestionResponse{
		Question: &rat.Question{
			ID:      qes.ID,
			OrgID:   qes.OrgID,
			UserID:  qes.UserID,
			QID:     qes.QID,
			ANS:     qes.ANS,
			QType:   qes.QType,
			Created: timestamppb.New(qes.Created),
			Updated: timestamppb.New(qes.Updated),
		},
	}, nil
}

func (s *Svc) UpsertMlTfQuestion(ctx context.Context, req *rat.RiskAssesmentQuestionRequest) (*rat.RiskAssesmentQuestionResponse, error) {
	log := logging.FromContext(ctx)
	qes, err := s.st.UpsertMlTfQuestion(ctx, &storage.Question{
		ID:             req.GetQuestion().GetID(),
		OrgID:          req.GetQuestion().GetOrgID(),
		UserID:         req.GetQuestion().GetUserID(),
		QID:            req.GetQuestion().GetQID(),
		QType:          req.GetQuestion().GetQType(),
		CustomersTotal: req.GetQuestion().GetCustomersTotal(),
		HrTotal:        req.GetQuestion().GetHrTotal(),
		ImpactScore:    req.GetQuestion().GetImpactScore(),
		Created:        req.GetQuestion().GetCreated().AsTime(),
		Updated:        req.GetQuestion().GetUpdated().AsTime(),
	})
	if err != nil {
		logging.WithError(err, log).Error("Upsert MLTF Question error")
		return nil, err
	}
	return &rat.RiskAssesmentQuestionResponse{
		Question: &rat.Question{
			ID:             qes.ID,
			OrgID:          qes.OrgID,
			UserID:         qes.UserID,
			QID:            qes.QID,
			QType:          qes.QType,
			CustomersTotal: qes.CustomersTotal,
			HrTotal:        qes.HrTotal,
			ImpactScore:    qes.ImpactScore,
			Created:        timestamppb.New(qes.Created),
			Updated:        timestamppb.New(qes.Updated),
		},
	}, nil
}

func (s *Svc) ListQuestion(ctx context.Context, req *rat.ListQuestionRequest) (*rat.ListQuestionResponse, error) {
	log := logging.FromContext(ctx)
	qes, err := s.st.ListQuestion(ctx, &storage.Question{
		OrgID:  req.GetOrgID(),
		UserID: req.GetUserID(),
		QID:    req.GetID(),
	})
	if err != nil {
		logging.WithError(err, log).Error("Question not Found")
		return nil, err
	}
	var question []*rat.Question
	for _, v := range qes {
		question = append(question, &rat.Question{
			ID:      v.ID,
			OrgID:   v.OrgID,
			UserID:  v.UserID,
			QID:     v.QID,
			ANS:     v.ANS,
			QType:   v.QType,
			Created: timestamppb.New(v.Created),
			Updated: timestamppb.New(v.Updated),
		})
	}
	return &rat.ListQuestionResponse{
		Question: question,
	}, nil
}

func (s *Svc) ListMlTfQuestion(ctx context.Context, req *rat.ListQuestionRequest) (*rat.ListQuestionResponse, error) {
	log := logging.FromContext(ctx)
	qes, err := s.st.ListMlTfQuestion(ctx, &storage.Question{
		OrgID:  req.GetOrgID(),
		UserID: req.GetUserID(),
		QID:    req.GetID(),
	})
	if err != nil {
		logging.WithError(err, log).Error("MLTF Question not Found")
		return nil, err
	}
	var question []*rat.Question
	for _, v := range qes {
		question = append(question, &rat.Question{
			ID:             v.ID,
			OrgID:          v.OrgID,
			UserID:         v.UserID,
			QID:            v.QID,
			QType:          v.QType,
			CustomersTotal: v.CustomersTotal,
			HrTotal:        v.HrTotal,
			ImpactScore:    v.ImpactScore,
			Created:        timestamppb.New(v.Created),
			Updated:        timestamppb.New(v.Updated),
		})
	}
	return &rat.ListQuestionResponse{
		Question: question,
	}, nil
}
