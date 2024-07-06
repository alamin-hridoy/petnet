package signup

import (
	"context"
	"errors"

	ipb "brank.as/rbac/gunk/v1/invite"

	upb "brank.as/petnet/gunk/dsa/v1/user"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) RetrieveInvite(ctx context.Context, code string) (*upb.RetrieveInviteResponse, error) {
	log := logging.FromContext(ctx)
	usr, err := s.icl.RetrieveInvite(ctx, &ipb.RetrieveInviteRequest{Code: code})
	if err != nil {
		logging.WithError(err, log).Error("failed to get user record by invite")
		return nil, errors.New("failed to get user record by invite")
	}
	return &upb.RetrieveInviteResponse{
		ID:           usr.GetID(),
		OrgID:        usr.GetOrgID(),
		Email:        usr.GetEmail(),
		CountryCode:  code,
		Phone:        usr.GetPhone(),
		CompanyName:  usr.GetCompanyName(),
		Active:       usr.GetActive(),
		InviteEmail:  usr.GetInviteEmail(),
		InviteStatus: usr.GetInviteStatus(),
		Invited:      usr.GetInvited(),
		FirstName:    usr.GetFirstName(),
		LastName:     usr.GetLastName(),
	}, nil
}
