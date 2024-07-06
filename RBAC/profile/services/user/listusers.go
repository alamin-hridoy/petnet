package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	ipupb "brank.as/rbac/gunk/v1/user"

	uapb "brank.as/rbac/profile/gunk/v1/useraccount"
)

func (h *Handler) ListUsers(ctx context.Context, req *uapb.ListUsersRequest) (*uapb.ListUsersResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.user.ListUsers")
	log.Info("request received")

	res, err := h.cl.ListUsers(ctx, &ipupb.ListUsersRequest{
		OrgID: req.OrgID,
	})
	if err != nil {
		st, ok := status.FromError(err)
		logging.WithError(err, log).Error("fetch users from identity-provider")
		if ok {
			return nil, status.Error(st.Code(), "no users found")
		}
		return nil, status.Error(codes.Internal, "no users found")
	}

	sus, err := h.acct.GetUsersByOrg(ctx, req.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("fetch users from storage")
		return nil, status.Error(codes.NotFound, "no users found")
	}

	us := res.GetUsers()
	list := make([]*uapb.User, len(us))
	for _, u := range us {
		for i, su := range sus {
			if u.ID == su.ID {
				// TODO: this is where profile will decorate with user info
				list[i] = &uapb.User{
					// Logo: su.Logo       example
					ID:           u.ID,
					OrgID:        u.OrgID,
					OrgName:      u.OrgName,
					FirstName:    u.FirstName,
					LastName:     u.LastName,
					Email:        u.Email,
					InviteStatus: u.InviteStatus,
					Created:      u.Created,
					Updated:      u.Updated,
				}
			}
		}
	}
	return &uapb.ListUsersResponse{Users: list}, nil
}
