package svcaccount

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/svcutil/mw"
	"brank.as/rbac/usermgm/storage"

	sapb "brank.as/rbac/gunk/v1/serviceaccount"
)

func (h *Svc) CreateAccount(ctx context.Context, req *sapb.CreateAccountRequest) (*sapb.CreateAccountResponse, error) {
	log := logging.FromContext(ctx).WithField("service", "svc.svcaccount.CreateAccount")
	log.Trace("request received")
	clID := hydra.ClientID(ctx)
	orgID := mw.GetOrg(ctx)
	if orgID == "" {
		log.WithField("id", clID).Error("no org found")
		return nil, status.Error(codes.PermissionDenied, "invalid organization")
	}
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Name, validation.Required, validation.Length(1, 0)),
		validation.Field(&req.Env, validation.In(h.env...)),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// r, a, _ := h.Permission(ctx, "CreateAccount")
	// _, err := h.val.ValidatePermission(ctx, &ppb.ValidatePermissionRequest{
	// 	ID:          clID,
	// 	Action:      a,
	// 	Resource:    r,
	// 	OrgID:       orgID,
	// 	Environment: req.GetEnv(),
	// })
	// if err != nil {
	// 	return nil, err
	// }

	atyp := storage.OAuth2
	if req.AuthType == sapb.AuthType_APIKey {
		atyp = storage.APIKey
	}

	id, sec, err := h.store.CreateSvcAccount(logging.WithLogger(ctx, log), storage.SvcAccount{
		AuthType:     atyp,
		OrgID:        orgID,
		Environment:  req.Env,
		ClientName:   req.Name,
		CreateUserID: clID,
	})
	if err != nil {
		return nil, err
	}
	return &sapb.CreateAccountResponse{
		ClientID: id,
		Secret:   sec,
	}, nil
}
