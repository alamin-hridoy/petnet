package svcaccount

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/mw"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	sapb "brank.as/rbac/gunk/v1/serviceaccount"
)

func (h *Svc) ListAccounts(ctx context.Context, req *sapb.ListAccountsRequest) (*sapb.ListAccountsResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.svcaccount.ListAccounts")
	log.Trace("request received")
	clID := hydra.ClientID(ctx)
	orgID := mw.GetOrg(ctx)
	if orgID == "" {
		log.WithField("id", clID).Error("no org found")
		return nil, status.Error(codes.PermissionDenied, "invalid organization")
	}
	accts, err := h.get.GetSvcAccountByOrgID(ctx, orgID)
	if err != nil {
		logging.WithError(err, log).Error("get accounts from store")
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "no service accounts found")
		}
		return nil, status.Error(codes.Internal, "failed to list service accounts")
	}
	resp := &sapb.ListAccountsResponse{
		Accounts: make([]*sapb.ServiceAccount, len(accts)),
	}
	for i, a := range accts {
		cr := tspb.New(a.Created)
		if !cr.IsValid() {
			logging.WithError(cr.CheckValid(), log).WithField("created", a.Created.String()).
				Error("create timestamp conversion")
			continue
		}
		resp.Accounts[i] = &sapb.ServiceAccount{
			Env:      a.Environment,
			ClientID: a.ClientID,
			Creator:  a.CreateUserID,
			Created:  cr,
		}
		if a.Disabled.Valid {
			ds := tspb.New(a.Disabled.Time)
			if !ds.IsValid() {
				logging.WithError(ds.CheckValid(), log).
					WithField("disabled", a.Disabled.Time.String()).
					Error("disabled timestamp conversion")
				continue
			}
			resp.Accounts[i].Disabled = ds
		}
	}
	return resp, nil
}
