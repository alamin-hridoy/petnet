package user

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/storage"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	uapb "brank.as/rbac/gunk/v1/user"
)

func (h *Handler) ListUsers(ctx context.Context, req *uapb.ListUsersRequest) (*uapb.ListUsersResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.user.ListAccounts")
	log.Trace("request received")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, is.UUIDv4),
		validation.Field(&req.ID, validation.Each(validation.Required, is.UUIDv4)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if req.OrgID == "" {
		req.OrgID = hydra.OrgID(ctx)
	}

	var sts []string
	for _, st := range req.GetStatus() {
		s := uapb.Status_name[int32(st.Number())]
		if s == "InviteSent" {
			sts = append(sts, "'Invite Sent'")
		} else {
			sts = append(sts, "'"+s+"'")
		}
	}

	usrs, err := h.acct.GetUsers(ctx, storage.FilterList{
		ID:           req.GetID(),
		OrgID:        req.GetOrgID(),
		Name:         req.GetName(),
		SortBy:       req.GetSortBy().String(),
		SortByColumn: req.GetSortByColumn().String(),
		Status:       sts,
		Limit:        req.GetLimit(),
		Offset:       req.GetOffset(),
	})
	if err != nil {
		logging.WithError(err, log).Error("no users found")
		return nil, status.Error(codes.NotFound, "no users found")
	}
	o, err := h.acct.GetOrgByID(ctx, req.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("no org found")
		return nil, status.Error(codes.NotFound, "user doesn't exist")
	}

	list := make([]*uapb.User, len(usrs))
	mp := make(map[string]*uapb.User, len(usrs))
	for i, u := range usrs {
		usr := &uapb.User{
			ID:           u.ID,
			OrgID:        u.OrgID,
			OrgName:      o.OrgName,
			FirstName:    u.FirstName,
			LastName:     u.LastName,
			Email:        u.Email,
			InviteStatus: u.InviteStatus,
			Created:      tspb.New(u.Created),
			Updated:      tspb.New(u.Updated),
			Deleted:      tspb.New(u.Deleted.Time),
		}
		list[i] = usr
		mp[u.ID] = usr
	}

	var total int32
	if len(usrs) > 0 {
		total = int32(usrs[0].Count)
	}
	return &uapb.ListUsersResponse{
		Users: list,
		Total: total,
		User:  mp,
	}, nil
}
