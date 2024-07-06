package terminal

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/util"
	"brank.as/petnet/serviceutil/logging"

	tpb "brank.as/petnet/gunk/drp/v1/terminal"
)

func (S *Svc) ConfirmRemit(ctx context.Context, req *tpb.ConfirmRemitRequest) (*tpb.ConfirmRemitResponse, error) {
	log := logging.FromContext(ctx)

	pn, err := S.remit.GetPartnerByTxnID(ctx, req.TransactionID)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	md := metautils.ExtractIncoming(ctx)
	ctx = md.Set(phmw.Partner, pn).ToIncoming(ctx)

	r, err := S.validators[pn].ConfirmRemitValidate(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("failed validating request")
		return nil, util.HandleServiceErr(err)
	}
	remit, err := S.remit.ProcessRemit(ctx, *r, pn)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return &tpb.ConfirmRemitResponse{ControlNumber: remit.ControlNumber}, nil
}

func (s *WUVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.AuthCode, validation.Required),
		validation.Field(&req.AuthSource, validation.Required),
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
		AuthSource:    req.TransactionID,
		AuthCode:      req.AuthCode,
	}
	return rReq, nil
}

func (s *IRVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *TFVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *RMVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *RIAVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *MBVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *BPIVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *USSCVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *ICVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *JPRVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *WISEVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *UNTVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *CEBVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *CEBIVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *AYAVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *IEVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}

func (s *PHUBVal) ConfirmRemitValidate(ctx context.Context, req *tpb.ConfirmRemitRequest) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TransactionID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.ProcessRemit{
		TransactionID: req.TransactionID,
	}
	return rReq, nil
}
