package signup

import (
	"context"

	"brank.as/petnet/profile/storage/postgres"
	"brank.as/petnet/serviceutil/logging"
	ipb "brank.as/rbac/gunk/v1/invite"
	opb "brank.as/rbac/gunk/v1/organization"
	spb "brank.as/rbac/gunk/v1/user"
)

type SignupReq struct {
	Username   string
	FirstName  string
	LastName   string
	Email      string
	Password   string
	InviteCode string
	OrgID      string
}

type SignupResp struct {
	UserID string
	OrgID  string
}

type Svc struct {
	st              *postgres.Storage
	scl             spb.SignupClient
	ocl             opb.OrganizationServiceClient
	icl             ipb.InviteServiceClient
	disableLoginMFA bool
}

func New(st *postgres.Storage, scl spb.SignupClient, ocl opb.OrganizationServiceClient, icl ipb.InviteServiceClient, dsbLoginMFA bool) *Svc {
	return &Svc{
		st:              st,
		scl:             scl,
		ocl:             ocl,
		icl:             icl,
		disableLoginMFA: dsbLoginMFA,
	}
}

func (s *Svc) Signup(ctx context.Context, r SignupReq) (*SignupResp, error) {
	log := logging.FromContext(ctx)
	res, err := s.scl.Signup(ctx, &spb.SignupRequest{
		Username:   r.Email,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		Email:      r.Email,
		Password:   r.Password,
		InviteCode: r.InviteCode,
	})
	if err != nil {
		logging.WithError(err, log).Error("signing up")
		return nil, err
	}

	o, err := s.ocl.GetOrganization(ctx, &opb.GetOrganizationRequest{ID: res.GetOrgID()})
	if err != nil {
		logging.WithError(err, log).Error("get org")
	}

	if !o.Organization[0].Active {
		if _, err := s.icl.Approve(ctx, &ipb.ApproveRequest{
			ID: res.UserID,
		}); err != nil {
			logging.WithError(err, log).Error("approve invitation")
		}
	}

	if !s.disableLoginMFA {
		if org, err := s.ocl.UpdateOrganization(ctx, &opb.UpdateOrganizationRequest{
			OrganizationID: res.GetOrgID(),
			LoginMFA:       opb.EnableOpt_Enable,
		}); err != nil {
			logging.WithError(err, log).Error("update org")
		} else if org.MFAEventID != "" {
			log.WithField("event_id", org.MFAEventID).Error("auto-enable failed")
		}
	}

	return &SignupResp{
		UserID: res.GetUserID(),
		OrgID:  res.GetOrgID(),
	}, nil
}
