package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"

	"brank.as/petnet/serviceutil/logging"
)

// UserLogin to authenticate with perahub gateway.
func (s *Svc) UserLogin(ctx context.Context, user, pass string) (*core.User, error) {
	log := logging.FromContext(ctx)

	c, err := s.p.Login(ctx, perahub.LoginRequest{
		Username: user,
		Password: pass,
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub login")
		return nil, status.Error(codes.NotFound, "login invalid")
	}

	if _, err := s.st.UpsertSession(ctx, storage.Session{
		CustomerCode: c.CustomerCode,
		Customer:     storage.Customer(*c),
	}); err != nil {
		logging.WithError(err, log).Error("perahub login")
		return nil, status.Error(codes.Internal, "login failed")
	}

	return &core.User{
		FrgnRefNo:     c.FrgnRefNo,
		CustNo:        c.CustomerCode,
		LastName:      c.LastName,
		FirstName:     c.FirstName,
		Birthdate:     c.Birthdate,
		Nationality:   c.Nationality,
		Address:       c.PresentAddress,
		Occupation:    c.Occupation,
		Employer:      c.EmployerName,
		ValidIdnt:     c.CustomerIDNo,
		WUCardNo:      c.Wucardno,
		DebitCardNo:   c.Debitcardno,
		LoyaltyCardNo: c.Loyaltycardno,
	}, nil
}
