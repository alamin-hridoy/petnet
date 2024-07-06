package userremit

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) StageCreateRemit(ctx context.Context, r core.Remittance, partner string) (*core.RemitResponse, error) {
	log := logging.FromContext(ctx)

	sess, err := s.st.GetSession(ctx, r.UserID)
	if err != nil {
		logging.WithError(err, log).Error("context")
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.Unauthenticated, "user not found")
		}
		return nil, err
	}
	c := sess.Customer
	bday, err := toDate(c.Birthdate)
	if err != nil {
		logging.WithError(err, log).Error("birthdate")
		return nil, status.Error(codes.InvalidArgument,
			"profile incomplete, missing birthdate. Please update and try again.")
	}
	issue, err := toDate(c.IDIssueDate)
	if err != nil {
		logging.WithError(err, log).Error("id issue date")
		return nil, status.Error(codes.InvalidArgument,
			"invalid registered ID. Please update and try again.")
	}
	expiry, err := toDate(c.IDExpiration)
	if err != nil && c.IDExpiration != "" {
		logging.WithError(err, log).Error("id expiry date")
		return nil, status.Error(codes.InvalidArgument,
			"invalid registered ID. Please update and try again.")
	}

	u := core.UserKYC{
		PartnerMemberID: c.Wucardno,
		FName:           c.FirstName,
		MdName:          c.MiddleName,
		LName:           c.LastName,
		Gender:          c.Gender,
		Address: core.Address{
			Address1:   c.PresentAddress,
			City:       c.City,
			State:      c.State,
			PostalCode: c.PostalCode,
			Country:    c.PresentCountry,
		},
		Phone:       core.PhoneNumber{Number: c.Phone},
		Mobile:      core.PhoneNumber{Number: c.Mobile},
		SourceFunds: r.Remitter.SourceFunds,
		Employment: core.Employment{
			Employer:   c.EmployerName,
			Occupation: c.Occupation,
			// PositionLevel: c.posi,
		},
		ReceiverRelation: r.Remitter.ReceiverRelation,
		PrimaryID: core.Identification{
			IDType:  c.IDType,
			Number:  c.CustomerIDNo,
			Country: c.CountryIDIssue,
			Issued:  *issue,
			Expiry:  expiry,
		},
		AlternateID:  r.Remitter.AlternateID,
		Email:        r.Remitter.Email,
		BirthDate:    *bday,
		BirthCountry: c.BirthCountry,
		Nationality:  c.Nationality,
	}
	r.Remitter = u
	return s.rmt.StageCreateRemit(ctx, r, partner)
}
