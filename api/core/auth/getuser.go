package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	phmw "brank.as/petnet/api/perahub-middleware"

	"brank.as/petnet/serviceutil/logging"
)

// UserLogin to authenticate with perahub gateway.
func (s *Svc) GetUser(ctx context.Context, id string) (*core.User, error) {
	log := logging.FromContext(ctx)

	c, err := s.p.MyWUSearch(ctx, perahub.WUSearchRequest{
		OperatorID: phmw.GetOperatorID(ctx),
		TerminalID: phmw.GetTerminalID(ctx),
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub login")
		return nil, status.Error(codes.NotFound, "login invalid")
	}

	return &core.User{
		FrgnRefNo:     c.ForRefNo,
		LastName:      c.Surname,
		FirstName:     c.GivenName,
		Birthdate:     c.Birthdate,
		Nationality:   c.Nationality,
		Address:       c.PresentAddress,
		Occupation:    c.Occupation,
		Employer:      c.NameOfEmployer,
		ValidIdnt:     c.ValidIdentification,
		WUCardNo:      c.WuCardNo,
		DebitCardNo:   c.DebitCardNo,
		LoyaltyCardNo: c.LoyaltyCardNo,
	}, nil
}
