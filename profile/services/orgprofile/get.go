package profile

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
)

func (h *Svc) GetProfile(ctx context.Context, req *ppb.GetProfileRequest) (*ppb.GetProfileResponse, error) {
	log := logging.FromContext(ctx)

	if req.OrgID == "" {
		return nil, status.Error(codes.InvalidArgument, "org_id cannot be empty")
	}

	pf, err := h.ps.GetOrgProfile(ctx, req.OrgID)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "org not found")
		}
		logging.WithError(err, log).Error("getting profile")
		return nil, status.Error(codes.Internal, "failed to get profile")
	}

	ppf := storageToProto(pf)
	return &ppb.GetProfileResponse{Profile: ppf}, nil
}

func (h *Svc) GetProfileByDsaCode(ctx context.Context, req *ppb.GetProfileByDsaCodeRequest) (*ppb.GetProfileByDsaCodeResponse, error) {
	log := logging.FromContext(ctx)

	if req.DsaCode == "" {
		return nil, status.Error(codes.InvalidArgument, "dsa_code cannot be empty")
	}

	pf, err := h.ps.GetProfileByDsaCode(ctx, req.DsaCode)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "org not found")
		}
		logging.WithError(err, log).Error("getting profile")
		return nil, status.Error(codes.Internal, "failed to get profile")
	}

	ppf := storageToProto(pf)
	return &ppb.GetProfileByDsaCodeResponse{Profile: ppf}, nil
}

func storageToProto(spf *storage.OrgProfile) *ppb.OrgProfile {
	ct := tspb.New(spf.Created)
	ut := tspb.New(spf.Updated)

	ppf := &ppb.OrgProfile{
		ID:               spf.ID,
		UserID:           spf.UserID,
		OrgID:            spf.OrgID,
		OrgType:          ppb.OrgType(spf.OrgType),
		Status:           ppb.Status(spf.Status),
		RiskScore:        ppb.RiskScore(spf.RiskScore),
		TransactionTypes: spf.TransactionTypes,
		BusinessInfo: &ppb.BusinessInfo{
			CompanyName:   spf.BusinessInfo.CompanyName,
			StoreName:     spf.BusinessInfo.StoreName,
			PhoneNumber:   spf.BusinessInfo.PhoneNumber,
			FaxNumber:     spf.BusinessInfo.FaxNumber,
			Website:       spf.BusinessInfo.Website,
			CompanyEmail:  spf.BusinessInfo.CompanyEmail,
			ContactPerson: spf.BusinessInfo.ContactPerson,
			Position:      spf.BusinessInfo.Position,
			Address: &ppb.Address{
				Address1:   spf.BusinessInfo.Address.Address1,
				City:       spf.BusinessInfo.Address.City,
				State:      spf.BusinessInfo.Address.State,
				PostalCode: spf.BusinessInfo.Address.PostalCode,
			},
		},
		AccountInfo: &ppb.AccountInfo{
			Bank:                    spf.AccountInfo.Bank,
			BankAccountNumber:       spf.AccountInfo.BankAccountNumber,
			BankAccountHolder:       spf.AccountInfo.BankAccountHolder,
			AgreeTermsConditions:    ppb.Boolean(spf.AccountInfo.AgreeTermsConditions),
			AgreeOnlineSupplierForm: ppb.Boolean(spf.AccountInfo.AgreeOnlineSupplierForm),
			Currency:                ppb.Currency(spf.Currency),
		},
		ReminderSent:      ppb.Boolean(spf.ReminderSent),
		Created:           ct,
		Updated:           ut,
		DsaCode:           spf.DsaCode,
		TerminalIdOtc:     spf.TerminalIdOtc,
		TerminalIdDigital: spf.TerminalIdDigital,
		IsProvider:        spf.IsProvider,
		Partner:           spf.Partner,
	}

	if spf.Deleted.Valid {
		ppf.Deleted = tspb.New(spf.Deleted.Time)
	}
	if spf.DateApplied.Valid {
		ppf.DateApplied = tspb.New(spf.DateApplied.Time)
	}
	return ppf
}

func split(s string, sep string) []string {
	ss := strings.Split(s, sep)
	if ss[0] == "" {
		return nil
	}
	return ss
}
