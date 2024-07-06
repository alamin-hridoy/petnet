package profile

import (
	"context"
	"database/sql"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Svc) UpsertProfile(ctx context.Context, req *ppb.UpsertProfileRequest) (*ppb.UpsertProfileResponse, error) {
	if err := validateProfile(ctx, req.Profile); err != nil {
		return nil, err
	}

	spf := profileToStorage(req.Profile, req.Profile.OrgID)
	pid, err := h.core.CreateOrgProfile(ctx, *spf)
	if err != nil {
		if err != storage.Conflict {
			return nil, status.Error(codes.Internal, "failed to create profile")
		}
		pid, err = h.core.UpdateOrgProfile(ctx, *spf)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to update profile")
		}
	}
	return &ppb.UpsertProfileResponse{ID: pid}, nil
}

func validateProfile(ctx context.Context, pf *ppb.OrgProfile) error {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(pf,
		// validation.Field(&pf.UserID, validation.Required, is.UUID),
		validation.Field(&pf.OrgID, validation.Required, is.UUID),
		validation.Field(&pf.BusinessInfo, validation.By(func(interface{}) error {
			bi := pf.GetBusinessInfo()
			if bi == nil {
				return nil
			}
			if err := validation.ValidateStruct(pf.BusinessInfo,
				validation.Field(&bi.Website, is.URL),
				validation.Field(&bi.CompanyEmail, is.Email),
			); err != nil {
				logging.WithError(err, log).Error("invalid request")
				return status.Error(codes.InvalidArgument, err.Error())
			}
			return nil
		})),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

func profileToStorage(ppf *ppb.OrgProfile, oid string) *storage.OrgProfile {
	spf := &storage.OrgProfile{
		OrgID:             oid,
		UserID:            ppf.GetUserID(),
		TransactionTypes:  ppf.GetTransactionTypes(),
		OrgType:           int(ppf.GetOrgType()),
		Status:            int(ppf.GetStatus()),
		RiskScore:         int(ppf.GetRiskScore()),
		ReminderSent:      int(ppf.GetReminderSent()),
		DateApplied:       sql.NullTime{Time: ppf.GetDateApplied().AsTime(), Valid: ppf.GetDateApplied().IsValid()},
		Deleted:           sql.NullTime{Time: ppf.GetDeleted().AsTime(), Valid: ppf.GetDeleted().IsValid()},
		DsaCode:           ppf.GetDsaCode(),
		TerminalIdOtc:     ppf.GetTerminalIdOtc(),
		TerminalIdDigital: ppf.GetTerminalIdDigital(),
		Partner:           ppf.GetPartner(),
		IsProvider:        ppf.GetIsProvider(),
	}

	pbi := ppf.GetBusinessInfo()
	spf.BusinessInfo = storage.BusinessInfo{
		CompanyName:   pbi.GetCompanyName(),
		StoreName:     pbi.GetStoreName(),
		PhoneNumber:   pbi.GetPhoneNumber(),
		FaxNumber:     pbi.GetFaxNumber(),
		Website:       pbi.GetWebsite(),
		CompanyEmail:  pbi.GetCompanyEmail(),
		ContactPerson: pbi.GetContactPerson(),
		Position:      pbi.GetPosition(),
		Address: storage.Address{
			Address1:   pbi.GetAddress().GetAddress1(),
			City:       pbi.GetAddress().GetCity(),
			State:      pbi.GetAddress().GetState(),
			PostalCode: pbi.GetAddress().GetPostalCode(),
		},
	}

	aip := ppf.GetAccountInfo()
	spf.AccountInfo = storage.AccountInfo{
		Bank:                    aip.GetBank(),
		BankAccountNumber:       aip.GetBankAccountNumber(),
		BankAccountHolder:       aip.GetBankAccountHolder(),
		AgreeTermsConditions:    int(aip.GetAgreeTermsConditions()),
		AgreeOnlineSupplierForm: int(aip.GetAgreeOnlineSupplierForm()),
		Currency:                int(aip.GetCurrency()),
	}
	return spf
}

func protoToNullTime(ts *tspb.Timestamp) sql.NullTime {
	if ts.IsValid() {
		return sql.NullTime{
			Time:  ts.AsTime(),
			Valid: true,
		}
	}
	return sql.NullTime{}
}

func (s *Svc) UpdateOrgProfileUserID(ctx context.Context, req *ppb.UpdateOrgProfileUserIDRequest) (*ppb.UpdateOrgProfileUserIDResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OldOrgID, validation.Required, is.UUID),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.UpdateOrgProfileUserID(ctx, storage.UpdateOrgProfileOrgIDUserID{
		OldOrgID: req.GetOldOrgID(),
		NewOrgID: req.GetNewOrgID(),
		UserID:   req.GetUserID(),
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to Update org profile")
	}
	return &ppb.UpdateOrgProfileUserIDResponse{ID: res}, nil
}
