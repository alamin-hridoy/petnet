package invite

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/rbac/serviceutil/logging"

	ipb "brank.as/rbac/gunk/v1/invite"
	"brank.as/rbac/usermgm/storage"
)

func (h *Svc) RetrieveInvite(ctx context.Context, req *ipb.RetrieveInviteRequest) (*ipb.RetrieveInviteResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.invite.RetrieveInvite")
	log.Trace("request received")

	err := validation.ValidateStruct(req,
		validation.Field(&req.ID, validation.Required.When(req.Code == "").Error("either code or id is required.")),
		validation.Field(&req.Code, validation.Required.When(req.ID == "").Error("either code or id is required.")),
	)
	if err != nil {
		logging.WithError(err, log).Info("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	usr := &storage.User{}
	if req.Code != "" {
		usr, err = h.store.GetUserByInvite(ctx, req.Code)
		if err != nil {
			logging.WithError(err, log).Error("failed to get user record by invite")
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if req.ID != "" {
		usr, err = h.store.GetUserByID(ctx, req.ID)
		if err != nil {
			logging.WithError(err, log).Error("failed to get user record by id")
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	org, err := h.plt.GetOrgByID(ctx, usr.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("failed to get organization record")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ipb.RetrieveInviteResponse{
		ID:           usr.ID,
		OrgID:        usr.OrgID,
		Email:        usr.Email,
		CompanyName:  org.OrgName,
		Active:       org.Active,
		InviteEmail:  usr.Email,
		InviteStatus: usr.InviteStatus,
		FirstName:    usr.FirstName,
		LastName:     usr.LastName,
		Invited:      timestamppb.New(usr.Updated),
	}, nil
}
