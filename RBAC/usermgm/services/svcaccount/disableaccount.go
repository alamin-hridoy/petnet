package svcaccount

import (
	"context"
	"database/sql"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	sapb "brank.as/rbac/gunk/v1/serviceaccount"
	"brank.as/rbac/svcutil/mw"
)

func (h *Svc) DisableAccount(ctx context.Context, req *sapb.DisableAccountRequest) (*sapb.DisableAccountResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.svcaccount.DisableAccount")
	log.Trace("request received")
	clID := hydra.ClientID(ctx)
	orgID := mw.GetOrg(ctx)
	if orgID == "" {
		log.WithField("id", clID).Error("no org found")
		return nil, status.Error(codes.PermissionDenied, "invalid organization")
	}
	err := validation.ValidateStruct(req,
		validation.Field(&req.Name, validation.Required, validation.Length(1, 0)))
	if err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	log = log.WithField("client_id", req.Name)
	acct, err := h.get.GetSvcAccountByID(ctx, req.Name)
	if err != nil {
		logging.WithError(err, log).Error("get accounts from store")
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "service account does not exist")
		}
		return nil, status.Error(codes.Internal, "failed to list service accounts")
	}
	if acct.OrgID != orgID {
		log.WithFields(logrus.Fields{
			"user":            clID,
			"platform":        orgID,
			"target_platform": acct.OrgID,
		}).Error("attempt to disable service account from a different platform")
		return nil, status.Error(codes.NotFound, "service account does not exist")
	}
	if acct.Disabled.Valid {
		return &sapb.DisableAccountResponse{}, nil
	}
	acct.DisableUserID = hydra.ClientID(ctx)
	tm, err := h.store.DisableSvcAccount(logging.WithLogger(ctx, log), *acct)
	if err != nil {
		logging.WithError(err, log).Error("disable service account in hydra")
		return nil, status.Error(codes.Internal, "failed to disable service account")
	}
	log.WithField("disabled", tm.String()).Info("service account disabled")
	return &sapb.DisableAccountResponse{}, nil
}
