package terminal

import (
	"context"
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/util"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/random"
)

func (s *Svc) CreateRemit(ctx context.Context, req *tpb.CreateRemitRequest) (*tpb.CreateRemitResponse, error) {
	log := logging.FromContext(ctx)
	pn := req.GetRemitPartner()
	orgType := phmw.GetOrgInfo(ctx)

	if orgType == phmw.Provider {
		pn = static.PerahubRemit
	}
	remType, err := s.lk.SendRemitType(ctx, pn, req.GetRemitType(), false)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	r, err := s.validators[pn].CreateRemitValidate(ctx, req, remType)
	if err != nil {
		logging.WithError(err, log).Error("validate request")
		return nil, util.HandleServiceErr(err)
	}

	rmt, err := s.remit.StageCreateRemit(ctx, *r, pn)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	chg := make(map[string]*tpb.Amount, len(rmt.Charges))
	for k, v := range rmt.Charges {
		chg[k] = &tpb.Amount{
			Amount:   v.Number(),
			Currency: v.CurrencyCode(),
		}
	}

	tx := make(map[string]*tpb.Amount, len(rmt.Charges))
	for k, v := range rmt.Taxes {
		chg[k] = &tpb.Amount{
			Amount:   v.Number(),
			Currency: v.CurrencyCode(),
		}
	}

	return &tpb.CreateRemitResponse{
		PrincipalAmount: &tpb.Amount{
			Amount:   rmt.PrincipalAmount.Number(),
			Currency: rmt.PrincipalAmount.CurrencyCode(),
		},
		RemitAmount: &tpb.Amount{
			Amount:   rmt.RemitAmount.Number(),
			Currency: rmt.RemitAmount.CurrencyCode(),
		},
		Taxes: tx,
		Tax: &tpb.Amount{
			Amount:   rmt.Tax.Number(),
			Currency: rmt.Tax.CurrencyCode(),
		},
		Charges: chg,
		TotalCharges: &tpb.Amount{
			Amount:   rmt.Charge.Number(),
			Currency: rmt.Charge.CurrencyCode(),
		},
		GrossTotal: &tpb.Amount{
			Amount:   rmt.GrossTotal.Number(),
			Currency: rmt.GrossTotal.CurrencyCode(),
		},
		PromoDescription: rmt.PromoDescription,
		PromoMessage:     rmt.PromoMessage,
		TransactionID:    rmt.TransactionID,
	}, nil
}

func (s *WUVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Message, validation.By(func(interface{}) error {
			return validation.Validate(strings.Split(req.Message, "\\n"),
				validation.Length(0, 14),                  // 14 lines max
				validation.Each(validation.Length(1, 79))) // 79 char max per line
		})),
		validation.Field(&req.OrderID, required, is.ASCII),

		validation.Field(&req.Remitter, required, validation.By(func(interface{}) error {
			r := req.Remitter
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withMobile, withPhoneCtryCode, withState)),
				validation.Field(&r.PartnerMemberID, is.Alphanumeric),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.Nationality, required, is.CountryCode2),
				validation.Field(&r.SourceFunds, required),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.SendingReason),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, validation.When(remType.Business,
					required, validation.By(func(interface{}) error {
						r := r.Identification
						return validation.ValidateStruct(r,
							validation.Field(&r.Type, validation.Required),
							validation.Field(&r.Number, validation.Required),
							validation.Field(&r.Country, is.CountryCode2, validation.Required),
							validation.Field(&r.City, validation.Required),
							validation.Field(&r.Expiration, validateDate(r.GetExpiration())),
							validation.Field(&r.Issued, validateDate(r.GetIssued())),
						)
					})),
				),
				validation.Field(&r.AlternateID, validation.Each(
					validation.By(func(id interface{}) error {
						i, _ := id.(*ppb.Identification)
						return valID(i, withCtry, withIss).Validate(nil)
					}),
				)),
				validation.Field(&r.Email, required, is.EmailFormat),
			)
		})),
		validation.Field(&req.Receiver, validation.When(remType.Receiver,
			validation.By(func(interface{}) error {
				r := req.Receiver
				dst := r.GetContactInfo().GetAddress().GetCountry() == "US" ||
					r.GetContactInfo().GetAddress().GetCountry() == "MX"
				return validation.ValidateStruct(r,
					validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withMobile, withPhoneCtryCode, withState, withNameOnly)),
					validation.Field(&r.DestinationCity, validation.When(dst, required)),
					validation.Field(&r.DestinationState, validation.When(dst, required)))
			})),
		),
		validation.Field(&req.Buiness, validation.When(remType.Business,
			required, validation.By(func(interface{}) error {
				r := req.Buiness
				return validation.ValidateStruct(r,
					validation.Field(&r.CompanyName, required, is.ASCII),
					validation.Field(&r.AccountCode, required, is.ASCII),
					validation.Field(&r.ControlNumber, required, is.ASCII),
					validation.Field(&r.Country, required, is.CountryCode2),
				)
			})),
		),

		validation.Field(&req.Account,
			validation.When(remType.BankAccount, required, validateAccount(req.Account))),
		validation.Field(&req.Amount, required, validation.By(func(interface{}) error {
			r := req.Amount
			return validation.ValidateStruct(r,
				validation.Field(&r.SourceCurrency, required, is.CurrencyCode),
				validation.Field(&r.DestinationCurrency, required, is.CurrencyCode),
				validation.Field(&r.Amount, required, is.Digit),
				validation.Field(&r.DestinationCountry, required, is.CountryCode2),
			)
		})),
		validation.Field(&req.Agent, required, validation.By(func(interface{}) error {
			r := req.Agent
			return validation.ValidateStruct(r,
				validation.Field(&r.UserID, required))
		})),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rReq := &core.Remittance{
		DsaID:         phmw.GetDSA(ctx),
		DsaOrderID:    req.OrderID,
		UserID:        hydra.ClientID(ctx),
		MyWUNumber:    phmw.GetWUNo(ctx),
		RemitPartner:  req.RemitPartner,
		SendRemitType: *remType,
		SendReason:    req.Remitter.SendingReason,
		TxnPurpose:    req.Remitter.TransactionPurpose,
		Remitter: core.UserKYC{
			PartnerMemberID: req.Remitter.PartnerMemberID,
			FName:           req.Remitter.ContactInfo.FirstName,
			MdName:          req.Remitter.ContactInfo.MiddleName,
			LName:           req.Remitter.ContactInfo.LastName,
			Gender:          req.Remitter.Gender.String(),
			Address:         core.ToAddr(req.Remitter.ContactInfo.Address),
			Phone:           core.ToPhone(req.Remitter.ContactInfo.Phone),
			Mobile:          core.ToPhone(req.Remitter.ContactInfo.Mobile),
			SourceFunds:     req.Remitter.SourceFunds,
			Employment: core.Employment{
				Occupation:    req.Remitter.Employment.Occupation,
				PositionLevel: req.Remitter.Employment.PositionLevel,
			},
			ReceiverRelation: req.Remitter.ReceiverRelation,
			PrimaryID:        *ToID(req.Remitter.Identification, ""),
			AlternateID: func() []core.Identification {
				ids := make([]core.Identification, len(req.Remitter.AlternateID))
				for i, id := range req.Remitter.AlternateID {
					ids[i] = *ToID(id, "")
				}
				return ids
			}(),
			Email:          req.Remitter.Email,
			BirthDate:      *core.ToDate(req.Remitter.Birthdate),
			BirthCountry:   req.Remitter.BirthCountry,
			Nationality:    req.Remitter.Nationality,
			CurrentAddress: core.ToAddr(req.Remitter.ContactInfo.Address),
		},
		Receiver: core.UserKYC{
			FName:   req.Receiver.ContactInfo.FirstName,
			MdName:  req.Receiver.ContactInfo.MiddleName,
			LName:   req.Receiver.ContactInfo.LastName,
			Address: core.ToAddr(req.Receiver.ContactInfo.Address),
			Phone:   core.ToPhone(req.Receiver.ContactInfo.Phone),
			Mobile:  core.ToPhone(req.Receiver.ContactInfo.Mobile),
		},
		SourceAmount: core.MustMinor(req.Amount.Amount, req.Amount.SourceCurrency),
		DestAmount:   core.MustMinor(req.Amount.Amount, req.Amount.DestinationCurrency),
		TargetDest:   req.Amount.DestinationAmount,
		TransactionDetails: core.TransactionDetails{
			DestCtry: req.Amount.DestinationCountry,
		},
		DestState: req.Receiver.DestinationState,
		DestCity:  req.Receiver.DestinationCity,
		DestAccount: core.Account{
			BIC:     req.GetAccount().GetBIC(),
			AcctNo:  req.GetAccount().GetAccountNumber(),
			AcctSfx: req.GetAccount().GetAccountSuffix(),
		},
		Promo:   req.Promo,
		Message: req.Message,
		Agent: core.Agent{
			UserID: int(req.Agent.UserID),
		},
	}
	return rReq, nil
}

func (s *WISEVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrderID, required, is.ASCII),
		validation.Field(&req.Remitter, required, validation.By(func(interface{}) error {
			r := req.Remitter
			return validation.ValidateStruct(r,
				validation.Field(&r.Email, required, is.EmailFormat),
			)
		})),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.RecipientID, required, is.Digit),
				validation.Field(&r.AccountHolderName, required, is.ASCII),
				validation.Field(&r.SourceAccountNumber, required, is.ASCII),
			)
		})),
		validation.Field(&req.Message, required, is.ASCII),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rReq := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       hydra.ClientID(ctx),
		DsaOrderID:   req.OrderID,
		RemitPartner: req.RemitPartner,
		Remitter: core.UserKYC{
			Email: req.Remitter.Email,
		},
		Receiver: core.UserKYC{
			RecipientID:         req.Receiver.RecipientID,
			AccountHolderName:   req.Receiver.AccountHolderName,
			SourceAccountNumber: req.Receiver.SourceAccountNumber,
		},
		Message: req.Message,
	}
	return rReq, nil
}

func (s *IRVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	return nil, fmt.Errorf("service not available for iremit")
}

func (s *TFVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	return nil, fmt.Errorf("service not available for transfast")
}

func (s *RMVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	return nil, fmt.Errorf("service not available for remitly")
}

func (s *RIAVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	return nil, fmt.Errorf("service not available for ria")
}

func (s *MBVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	return nil, fmt.Errorf("service not available for metrobank")
}

func (s *BPIVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	return nil, fmt.Errorf("service not available for BPI")
}

func (s *USSCVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required),
		validation.Field(&req.Amount, required, validation.By(func(interface{}) error {
			r := req.Amount
			return validation.ValidateStruct(r,
				validation.Field(&r.SourceCurrency, required, is.CurrencyCode),
				validation.Field(&r.Amount, required, is.Digit),
			)
		})),
		validation.Field(&req.Remitter, required, validation.By(func(interface{}) error {
			r := req.Remitter
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Gender, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.KYCVerified, required),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, validation.By(func(interface{}) error {
					r := r.ContactInfo
					return validation.ValidateStruct(r,
						validation.Field(&r.FirstName, required),
						validation.Field(&r.LastName, required),
						validation.Field(&r.MiddleName),
						validation.Field(&r.Phone, required, valPhone(r.Phone, false)),
					)
				})),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sendctry, err := s.q.FindCountryByAlpha(req.Remitter.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid sender contactinfo country")
	}
	bctry, err := s.q.FindCountryByAlpha(req.Remitter.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	rm := &core.Remittance{
		DsaID:         phmw.GetDSA(ctx),
		DsaOrderID:    req.OrderID,
		UserID:        hydra.ClientID(ctx),
		MyWUNumber:    phmw.GetWUNo(ctx),
		RemitPartner:  req.RemitPartner,
		SendRemitType: *remType,
		TxnPurpose:    req.Remitter.TransactionPurpose,
		Remitter: core.UserKYC{
			PartnerMemberID: req.Remitter.PartnerMemberID,
			FName:           req.Remitter.ContactInfo.FirstName,
			MdName:          req.Remitter.ContactInfo.MiddleName,
			LName:           req.Remitter.ContactInfo.LastName,
			Gender:          req.Remitter.Gender.String(),
			Address: core.Address{
				Address1:   req.Remitter.ContactInfo.Address.Address1,
				Address2:   req.Remitter.ContactInfo.Address.Address2,
				City:       req.Remitter.ContactInfo.Address.City,
				Province:   req.Remitter.ContactInfo.Address.Province,
				PostalCode: req.Remitter.ContactInfo.Address.PostalCode,
				Country:    sendctry.Name.Common,
				Zone:       req.Remitter.ContactInfo.Address.Zone,
			},
			Phone:       core.ToPhone(req.Remitter.ContactInfo.Phone),
			Mobile:      core.ToPhone(req.Remitter.ContactInfo.Mobile),
			SourceFunds: req.Remitter.SourceFunds,
			Employment: core.Employment{
				Occupation:    req.Remitter.Employment.Occupation,
				PositionLevel: req.Remitter.Employment.PositionLevel,
			},
			ReceiverRelation: req.Remitter.ReceiverRelation,
			PrimaryID:        *ToID(req.Remitter.Identification, ""),
			BirthDate:        *core.ToDate(req.Remitter.Birthdate),
			BirthPlace:       req.Remitter.BirthPlace,
			BirthCountry:     bctry.Name.Common,
			CurrentAddress:   core.ToAddr(req.Remitter.ContactInfo.Address),
		},
		Receiver: core.UserKYC{
			FName:  req.Receiver.ContactInfo.FirstName,
			MdName: req.Receiver.ContactInfo.MiddleName,
			LName:  req.Receiver.ContactInfo.LastName,
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
		},
		SourceAmount: core.MustMinor(req.Amount.Amount, req.Amount.SourceCurrency),
		DestAmount:   core.MustMinor("0", req.Amount.SourceCurrency),
		GrossTotal:   core.MustMinor("0", req.Amount.SourceCurrency),
		TransactionDetails: core.TransactionDetails{
			// should be static
			SrcCtry:    "Philippines",
			DestCtry:   "PH",
			IsDomestic: perahub.IsDomestic(req.Amount.SourceCountry, req.Amount.DestinationCountry),
		},
		Agent: core.Agent{
			IPAddress: req.Agent.GetIPAddress(),
			UserID:    int(req.Agent.UserID),
		},
	}
	return rm, nil
}

func (s *ICVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	return nil, fmt.Errorf("service not available for InstaCash")
}

func (s *JPRVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	return nil, fmt.Errorf("service not available for JapanRemit")
}

func (s *UNTVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	return nil, fmt.Errorf("service not available for Uniteller")
}

func (s *CEBVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required),
		validation.Field(&req.Amount, required, validation.By(func(interface{}) error {
			r := req.Amount
			return validation.ValidateStruct(r,
				validation.Field(&r.SourceCurrency, required, is.CurrencyCode),
				validation.Field(&r.Amount, required, is.Digit),
				validation.Field(&r.DestinationCountry, required, is.CountryCode2),
			)
		})),
		validation.Field(&req.Remitter, required, validation.By(func(interface{}) error {
			r := req.Remitter
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.RecipientID, required),
				validation.Field(&r.ContactInfo, required, validation.By(func(interface{}) error {
					r := r.ContactInfo
					return validation.ValidateStruct(r,
						validation.Field(&r.FirstName, required),
						validation.Field(&r.LastName, required),
						validation.Field(&r.MiddleName),
					)
				})),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctry, err := s.q.FindCountryByAlpha(req.Remitter.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}
	bctry, err := s.q.FindCountryByAlpha(req.Remitter.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	srcctry, err := s.q.FindCountryByAlpha(req.Amount.SourceCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid amount source country")
	}
	rm := &core.Remittance{
		DsaID:         phmw.GetDSA(ctx),
		DsaOrderID:    req.OrderID,
		UserID:        hydra.ClientID(ctx),
		MyWUNumber:    phmw.GetWUNo(ctx),
		RemitPartner:  req.RemitPartner,
		SendRemitType: *remType,
		TxnPurpose:    req.Remitter.TransactionPurpose,
		Remitter: core.UserKYC{
			PartnerMemberID: req.Remitter.PartnerMemberID,
			FName:           req.Remitter.ContactInfo.FirstName,
			MdName:          req.Remitter.ContactInfo.MiddleName,
			LName:           req.Remitter.ContactInfo.LastName,
			Address: core.Address{
				Address1:   req.Remitter.ContactInfo.Address.Address1,
				Address2:   req.Remitter.ContactInfo.Address.Address2,
				City:       req.Remitter.ContactInfo.Address.City,
				Province:   req.Remitter.ContactInfo.Address.Province,
				PostalCode: req.Remitter.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				Zone:       req.Remitter.ContactInfo.Address.Zone,
			},
			Phone:       core.ToPhone(req.Remitter.ContactInfo.Phone),
			SourceFunds: req.Remitter.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Remitter.Employment.Occupation,
			},
			ReceiverRelation: req.Remitter.ReceiverRelation,
			PrimaryID:        *ToID(req.Remitter.Identification, ""),
			BirthDate:        *core.ToDate(req.Remitter.Birthdate),
			BirthPlace:       req.Remitter.BirthPlace,
			BirthCountry:     bctry.Name.Common,
		},
		Receiver: core.UserKYC{
			FName:       req.Receiver.ContactInfo.FirstName,
			MdName:      req.Receiver.ContactInfo.MiddleName,
			LName:       req.Receiver.ContactInfo.LastName,
			RecipientID: req.Receiver.RecipientID,
		},
		SourceAmount: core.MustMinor(req.Amount.Amount, req.Amount.SourceCurrency),
		DestAmount:   core.MustMinor("0", req.Amount.SourceCurrency),
		GrossTotal:   core.MustMinor("0", req.Amount.SourceCurrency),
		TransactionDetails: core.TransactionDetails{
			SrcCtry:    srcctry.Name.Common,
			DestCtry:   req.Amount.DestinationCountry,
			IsDomestic: perahub.IsDomestic(req.Amount.SourceCountry, req.Amount.DestinationCountry),
		},
		Agent: core.Agent{
			IPAddress: req.Agent.GetIPAddress(),
			UserID:    int(req.Agent.GetUserID()),
		},
	}
	return rm, nil
}

func (s *CEBIVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	return nil, fmt.Errorf("service not available for Cebuana Intl")
}

func (s *AYAVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required),
		validation.Field(&req.Amount, required, validation.By(func(interface{}) error {
			r := req.Amount
			return validation.ValidateStruct(r,
				validation.Field(&r.SourceCurrency, required, is.CurrencyCode),
				validation.Field(&r.Amount, required, is.Digit),
				validation.Field(&r.DestinationCountry, required, is.CountryCode2),
			)
		})),
		validation.Field(&req.Remitter, required, validation.By(func(interface{}) error {
			r := req.Remitter
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, validation.By(func(interface{}) error {
					r := r.ContactInfo
					return validation.ValidateStruct(r,
						validation.Field(&r.FirstName, required),
						validation.Field(&r.LastName, required),
						validation.Field(&r.MiddleName),
					)
				})),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	bctry, err := s.q.FindCountryByAlpha(req.Remitter.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	srccountry := ""
	if req.Amount.SourceCountry != "" {
		srcctry, err := s.q.FindCountryByAlpha(req.Amount.SourceCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid amount source country")
		}
		srccountry = srcctry.Name.Common
	}
	destctry, err := s.q.FindCountryByAlpha(req.Amount.DestinationCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid amount destination country")
	}

	rm := &core.Remittance{
		DsaID:         phmw.GetDSA(ctx),
		DsaOrderID:    req.OrderID,
		UserID:        hydra.ClientID(ctx),
		MyWUNumber:    phmw.GetWUNo(ctx),
		RemitPartner:  req.RemitPartner,
		SendRemitType: *remType,
		ControlNo:     time.Now().Local().Format("20060102") + "PHB" + random.NumberString(9),
		SendReason:    req.Remitter.SendingReason,
		TxnPurpose:    req.Remitter.TransactionPurpose,
		Remitter: core.UserKYC{
			PartnerMemberID: req.Remitter.PartnerMemberID,
			FName:           req.Remitter.ContactInfo.FirstName,
			MdName:          req.Remitter.ContactInfo.MiddleName,
			LName:           req.Remitter.ContactInfo.LastName,
			Gender:          req.Remitter.Gender.String(),
			Address: core.Address{
				Address1:   req.Remitter.ContactInfo.Address.Address1,
				Address2:   req.Remitter.ContactInfo.Address.Address2,
				City:       req.Remitter.ContactInfo.Address.City,
				Province:   req.Remitter.ContactInfo.Address.Province,
				PostalCode: req.Remitter.ContactInfo.Address.PostalCode,
				Country:    req.Remitter.ContactInfo.Address.Country,
				Zone:       req.Remitter.ContactInfo.Address.Zone,
			},
			Phone:       core.ToPhone(req.Remitter.ContactInfo.Phone),
			Mobile:      core.ToPhone(req.Remitter.ContactInfo.Mobile),
			SourceFunds: req.Remitter.SourceFunds,
			Employment: core.Employment{
				Occupation:    req.Remitter.Employment.Occupation,
				PositionLevel: req.Remitter.Employment.PositionLevel,
			},
			ReceiverRelation: req.Remitter.ReceiverRelation,
			PrimaryID:        *ToID(req.Remitter.Identification, ""),
			AlternateID: func() []core.Identification {
				ids := make([]core.Identification, len(req.Remitter.AlternateID))
				for i, id := range req.Remitter.AlternateID {
					ids[i] = *ToID(id, "")
				}
				return ids
			}(),
			Email:          req.Remitter.Email,
			BirthDate:      *core.ToDate(req.Remitter.Birthdate),
			BirthPlace:     req.Remitter.BirthPlace,
			BirthCountry:   bctry.Name.Common,
			Nationality:    req.Remitter.Nationality,
			CurrentAddress: core.ToAddr(req.Remitter.ContactInfo.Address),
		},
		Receiver: core.UserKYC{
			FName:  req.Receiver.ContactInfo.FirstName,
			MdName: req.Receiver.ContactInfo.MiddleName,
			LName:  req.Receiver.ContactInfo.LastName,
		},
		SourceAmount: core.MustMinor(req.Amount.Amount, req.Amount.SourceCurrency),
		DestAmount:   core.MustMinor("0", req.Amount.SourceCurrency),
		GrossTotal:   core.MustMinor("0", req.Amount.SourceCurrency),
		TargetDest:   req.Amount.DestinationAmount,
		DestState:    req.Receiver.DestinationState,
		DestCity:     req.Receiver.DestinationCity,
		DestAccount: core.Account{
			BIC:     req.GetAccount().GetBIC(),
			AcctNo:  req.GetAccount().GetAccountNumber(),
			AcctSfx: req.GetAccount().GetAccountSuffix(),
		},
		TransactionDetails: core.TransactionDetails{
			SrcCtry:    srccountry,
			DestCtry:   destctry.Name.Common,
			IsDomestic: perahub.IsDomestic(req.Amount.SourceCountry, req.Amount.DestinationCountry),
		},
		Agent: core.Agent{
			IPAddress: req.Agent.GetIPAddress(),
			UserID:    int(req.Agent.UserID),
		},
		Promo:   req.Promo,
		Message: req.Message,
	}
	return rm, nil
}

func (s *IEVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required),
		validation.Field(&req.Amount, required, validation.By(func(interface{}) error {
			r := req.Amount
			return validation.ValidateStruct(r,
				validation.Field(&r.SourceCurrency, required, is.CurrencyCode),
				validation.Field(&r.Amount, required, is.Digit),
				validation.Field(&r.DestinationCountry, required, is.CountryCode2),
			)
		})),
		validation.Field(&req.Remitter, required, validation.By(func(interface{}) error {
			r := req.Remitter
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.Identification, required, validation.By(func(interface{}) error {
					r := r.Identification
					return validation.ValidateStruct(r,
						validation.Field(&r.Number, required),
					)
				})),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	recctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid receiver contactinfo country")
	}
	bctry, err := s.q.FindCountryByAlpha(req.Remitter.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}

	rm := &core.Remittance{
		DsaID:         phmw.GetDSA(ctx),
		DsaOrderID:    req.OrderID,
		UserID:        hydra.ClientID(ctx),
		MyWUNumber:    phmw.GetWUNo(ctx),
		RemitPartner:  req.RemitPartner,
		SendRemitType: *remType,
		SendReason:    req.Remitter.SendingReason,
		TxnPurpose:    req.Remitter.TransactionPurpose,
		Remitter: core.UserKYC{
			PartnerMemberID: req.Remitter.PartnerMemberID,
			FName:           req.Remitter.ContactInfo.FirstName,
			MdName:          req.Remitter.ContactInfo.MiddleName,
			LName:           req.Remitter.ContactInfo.LastName,
			Gender:          req.Remitter.Gender.String(),
			Address: core.Address{
				Address1:   req.Remitter.ContactInfo.Address.Address1,
				Address2:   req.Remitter.ContactInfo.Address.Address2,
				City:       req.Remitter.ContactInfo.Address.City,
				Province:   req.Remitter.ContactInfo.Address.Province,
				PostalCode: req.Remitter.ContactInfo.Address.PostalCode,
				Country:    req.Remitter.ContactInfo.Address.Country,
				Zone:       req.Remitter.ContactInfo.Address.Zone,
			},
			Phone:       core.ToPhone(req.Remitter.ContactInfo.Phone),
			Mobile:      core.ToPhone(req.Remitter.ContactInfo.Mobile),
			SourceFunds: req.Remitter.SourceFunds,
			Employment: core.Employment{
				Occupation:    req.Remitter.Employment.Occupation,
				PositionLevel: req.Remitter.Employment.PositionLevel,
			},
			ReceiverRelation: req.Remitter.ReceiverRelation,
			PrimaryID:        *ToID(req.Remitter.Identification, ""),
			AlternateID: func() []core.Identification {
				ids := make([]core.Identification, len(req.Remitter.AlternateID))
				for i, id := range req.Remitter.AlternateID {
					ids[i] = *ToID(id, "")
				}
				return ids
			}(),
			Email:          req.Remitter.Email,
			BirthDate:      *core.ToDate(req.Remitter.Birthdate),
			BirthPlace:     req.Remitter.BirthPlace,
			BirthCountry:   bctry.Name.Common,
			Nationality:    req.Remitter.Nationality,
			CurrentAddress: core.ToAddr(req.Remitter.ContactInfo.Address),
		},
		Receiver: core.UserKYC{
			FName:  req.Receiver.ContactInfo.FirstName,
			MdName: req.Receiver.ContactInfo.MiddleName,
			LName:  req.Receiver.ContactInfo.LastName,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.Receiver.ContactInfo.Address.Address2,
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    recctry.Name.Common,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				CtyCode: req.Receiver.ContactInfo.Phone.CountryCode,
				Number:  req.Receiver.ContactInfo.Phone.Number,
			},
			PrimaryID: core.Identification{
				Number: req.Receiver.Identification.Number,
			},
		},
		SourceAmount: core.MustMinor(req.Amount.Amount, req.Amount.SourceCurrency),
		DestAmount:   core.MustMinor("0", req.Amount.SourceCurrency),
		GrossTotal:   core.MustMinor("0", req.Amount.SourceCurrency),
		TargetDest:   req.Amount.DestinationAmount,
		DestState:    req.Receiver.DestinationState,
		DestCity:     req.Receiver.DestinationCity,
		DestAccount: core.Account{
			BIC:     req.GetAccount().GetBIC(),
			AcctNo:  req.GetAccount().GetAccountNumber(),
			AcctSfx: req.GetAccount().GetAccountSuffix(),
		},
		TransactionDetails: core.TransactionDetails{
			SrcCtry:    req.Amount.SourceCountry,
			DestCtry:   req.Amount.DestinationCountry,
			IsDomestic: perahub.IsDomestic(req.Amount.SourceCountry, req.Amount.DestinationCountry),
		},
		Agent: core.Agent{
			IPAddress: req.Agent.GetIPAddress(),
			UserID:    int(req.Agent.UserID),
		},
		Promo:   req.Promo,
		Message: req.Message,
	}
	return rm, nil
}

func (s *PHUBVal) CreateRemitValidate(ctx context.Context, req *tpb.CreateRemitRequest, remType *core.SendRemitType) (*core.Remittance, error) {
	return nil, fmt.Errorf("service not available for PerahubRemit")
}
