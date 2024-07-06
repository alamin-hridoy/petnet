package partner

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/petnet/gunk/dsa/v2/partner"
)

func (s *Svc) UpdatePartners(ctx context.Context, req *ppb.UpdatePartnersRequest) (*ppb.UpdatePartnersResponse, error) {
	required := validation.Required
	pnr := req.GetPartners()
	if err := validation.ValidateStruct(pnr,
		validation.Field(&pnr.OrgID, validation.Required, is.UUIDv4),
		validation.Field(&pnr.WesternUnionPartner, validation.By(func(interface{}) error {
			r := pnr.GetWesternUnionPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Coy, required),
				validation.Field(&r.TerminalID, required),
			)
		})),
		validation.Field(&pnr.IRemitPartner, validation.By(func(interface{}) error {
			r := pnr.GetIRemitPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.TransfastPartner, validation.By(func(interface{}) error {
			r := pnr.GetTransfastPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.RemitlyPartner, validation.By(func(interface{}) error {
			r := pnr.GetRemitlyPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.RiaPartner, validation.By(func(interface{}) error {
			r := pnr.GetRiaPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.MetroBankPartner, validation.By(func(interface{}) error {
			r := pnr.GetMetroBankPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.BPIPartner, validation.By(func(interface{}) error {
			r := pnr.GetBPIPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.USSCPartner, validation.By(func(interface{}) error {
			r := pnr.GetUSSCPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.JapanRemitPartner, validation.By(func(interface{}) error {
			r := pnr.GetJapanRemitPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.InstantCashPartner, validation.By(func(interface{}) error {
			r := pnr.GetInstantCashPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.UnitellerPartner, validation.By(func(interface{}) error {
			r := pnr.GetUnitellerPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.CebuanaPartner, validation.By(func(interface{}) error {
			r := pnr.GetCebuanaPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.TransferWisePartner, validation.By(func(interface{}) error {
			r := pnr.GetTransferWisePartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.CebuanaIntlPartner, validation.By(func(interface{}) error {
			r := pnr.GetCebuanaIntlPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.AyannahPartner, validation.By(func(interface{}) error {
			r := pnr.GetAyannahPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
		validation.Field(&pnr.IntelExpressPartner, validation.By(func(interface{}) error {
			r := pnr.GetIntelExpressPartner()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, required),
				validation.Field(&r.Param1, required),
				validation.Field(&r.Param2, required),
				// todo: validation once we know parameters
			)
		})),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.core.UpdatePartners(ctx, pnr); err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update service record")
	}
	return &ppb.UpdatePartnersResponse{Partners: pnr}, nil
}
