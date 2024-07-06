package auth

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	authpb "brank.as/rbac/gunk/v1/authenticate"
	osapb "brank.as/rbac/gunk/v1/oauth2"
)

func (s *Svc) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.Session, error) {
	log := logging.FromContext(ctx)
	log.WithField("username", req.Username).Debug("received")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.Username, validation.Required),
		validation.Field(&req.Password, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	u, err := s.src.UserLogin(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.NotFound, "login failed")
	}
	r, err := phmw.UserSession(ctx, *u)
	if err != nil {
		logging.WithError(err, log).Error("session conversion")
		return nil, status.Error(codes.Internal, "login failed")
	}

	res, err := s.cl.ListClients(ctx, &osapb.ListClientsRequest{
		ClientID: req.ClientID,
		OrgID:    mw.GetOrgID(ctx),
	})
	if err != nil {
		logging.WithError(err, log).Error("listing clients")
		return nil, status.Error(codes.NotFound, "login failed")
	}

	r.Session[phmw.UsrName] = req.Username
	r.Session["orgtype"] = ppb.OrgType_DSA.String()
	r.Session["environment"] = res.GetClients()[0].GetEnv()
	return r, nil
}
