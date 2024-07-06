package cebuana

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) RegisterUser(ctx context.Context, req core.RegisterUserReq) (*core.RegisterUserResp, error) {
	log := logging.FromContext(ctx)

	res, err := s.ph.CebAddClient(ctx, perahub.CebAddClientReq{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: req.BirthDate,
		CpCtryID:  req.CpCtryID,
		ContactNo: req.ContactNo,
		TpCtryID:  req.TpCtryID,
		TpArCode:  req.TpArCode,
		CrtyAdID:  req.CrtyAdID,
		PAdd:      req.PAdd,
		CAdd:      req.CAdd,
		UserID:    req.UserID,
		SOFID:     req.SOFID,
		Tin:       req.Tin,
		TpNo:      req.TpNo,
		AgentCode: req.AgentCode,
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub error")
		return nil, err
	}
	return &core.RegisterUserResp{
		Code:    res.Code,
		Message: res.Message,
		Result: core.RUResult{
			ResultStatus: res.Result.ResultStatus,
			MessageID:    res.Result.MessageID,
			LogID:        res.Result.LogID,
			ClientID:     res.Result.ClientID,
			ClientNo:     res.Result.ClientNo,
		},
		RemcoID: res.RemcoID,
	}, nil
}
