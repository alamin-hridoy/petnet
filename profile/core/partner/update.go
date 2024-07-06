package partner

import (
	"context"
	"encoding/json"

	spb "brank.as/petnet/gunk/dsa/v2/partner"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdatePartners(ctx context.Context, pnr *spb.Partners) error {
	log := logging.FromContext(ctx)

	if sv := pnr.GetWesternUnionPartner(); sv != nil {
		b, err := json.Marshal(storage.WesternUnionPartner{
			Coy:        sv.Coy,
			TerminalID: sv.TerminalID,
			StartDate:  sv.StartDate.AsTime(),
			EndDate:    sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.WesternUnionPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if sv := pnr.GetIRemitPartner(); sv != nil {
		b, err := json.Marshal(storage.IRemitPartner{
			// todo: once we have parameters
			Param1:    pnr.IRemitPartner.Param1,
			Param2:    pnr.IRemitPartner.Param2,
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.IRemitPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetTransfastPartner() != nil {
		b, err := json.Marshal(storage.TransfastPartner{
			// todo: once we have parameters
			Param1:    pnr.TransfastPartner.Param1,
			Param2:    pnr.TransfastPartner.Param2,
			StartDate: pnr.TransfastPartner.StartDate.AsTime(),
			EndDate:   pnr.TransfastPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.TransfastPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetRemitlyPartner() != nil {
		b, err := json.Marshal(storage.RemitlyPartner{
			// todo: once we have parameters
			Param1:    pnr.RemitlyPartner.Param1,
			Param2:    pnr.RemitlyPartner.Param2,
			StartDate: pnr.RemitlyPartner.StartDate.AsTime(),
			EndDate:   pnr.RemitlyPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.RemitlyPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetRiaPartner() != nil {
		b, err := json.Marshal(storage.RiaPartner{
			// todo: once we have parameters
			Param1:    pnr.RiaPartner.Param1,
			Param2:    pnr.RiaPartner.Param2,
			StartDate: pnr.RiaPartner.StartDate.AsTime(),
			EndDate:   pnr.RiaPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.RiaPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetMetroBankPartner() != nil {
		b, err := json.Marshal(storage.MetroBankPartner{
			// todo: once we have parameters
			Param1:    pnr.MetroBankPartner.Param1,
			Param2:    pnr.MetroBankPartner.Param2,
			StartDate: pnr.MetroBankPartner.StartDate.AsTime(),
			EndDate:   pnr.MetroBankPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.MetroBankPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetBPIPartner() != nil {
		b, err := json.Marshal(storage.BPIPartner{
			// todo: once we have parameters
			Param1:    pnr.BPIPartner.Param1,
			Param2:    pnr.BPIPartner.Param2,
			StartDate: pnr.BPIPartner.StartDate.AsTime(),
			EndDate:   pnr.BPIPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.BPIPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetUSSCPartner() != nil {
		b, err := json.Marshal(storage.USSCPartner{
			// todo: once we have parameters
			Param1:    pnr.USSCPartner.Param1,
			Param2:    pnr.USSCPartner.Param2,
			StartDate: pnr.USSCPartner.StartDate.AsTime(),
			EndDate:   pnr.USSCPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.USSCPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetJapanRemitPartner() != nil {
		b, err := json.Marshal(storage.JapanRemitPartner{
			// todo: once we have parameters
			Param1:    pnr.JapanRemitPartner.Param1,
			Param2:    pnr.JapanRemitPartner.Param2,
			StartDate: pnr.JapanRemitPartner.StartDate.AsTime(),
			EndDate:   pnr.JapanRemitPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.JapanRemitPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetInstantCashPartner() != nil {
		b, err := json.Marshal(storage.InstantCashPartner{
			// todo: once we have parameters
			Param1:    pnr.InstantCashPartner.Param1,
			Param2:    pnr.InstantCashPartner.Param2,
			StartDate: pnr.InstantCashPartner.StartDate.AsTime(),
			EndDate:   pnr.InstantCashPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.InstantCashPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetUnitellerPartner() != nil {
		b, err := json.Marshal(storage.UnitellerPartner{
			// todo: once we have parameters
			Param1:    pnr.UnitellerPartner.Param1,
			Param2:    pnr.UnitellerPartner.Param2,
			StartDate: pnr.UnitellerPartner.StartDate.AsTime(),
			EndDate:   pnr.UnitellerPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.UnitellerPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetCebuanaPartner() != nil {
		b, err := json.Marshal(storage.CebuanaPartner{
			// todo: once we have parameters
			Param1:    pnr.CebuanaPartner.Param1,
			Param2:    pnr.CebuanaPartner.Param2,
			StartDate: pnr.CebuanaPartner.StartDate.AsTime(),
			EndDate:   pnr.CebuanaPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.CebuanaPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetTransferWisePartner() != nil {
		b, err := json.Marshal(storage.TransferWisePartner{
			// todo: once we have parameters
			Param1:    pnr.TransferWisePartner.Param1,
			Param2:    pnr.TransferWisePartner.Param2,
			StartDate: pnr.TransferWisePartner.StartDate.AsTime(),
			EndDate:   pnr.TransferWisePartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.TransferWisePartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetCebuanaIntlPartner() != nil {
		b, err := json.Marshal(storage.CebuanaIntlPartner{
			// todo: once we have parameters
			Param1:    pnr.CebuanaIntlPartner.Param1,
			Param2:    pnr.CebuanaIntlPartner.Param2,
			StartDate: pnr.CebuanaIntlPartner.StartDate.AsTime(),
			EndDate:   pnr.CebuanaIntlPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.CebuanaIntlPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetAyannahPartner() != nil {
		b, err := json.Marshal(storage.AyannahPartner{
			// todo: once we have parameters
			Param1:    pnr.AyannahPartner.Param1,
			Param2:    pnr.AyannahPartner.Param2,
			StartDate: pnr.AyannahPartner.StartDate.AsTime(),
			EndDate:   pnr.AyannahPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.AyannahPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	if pnr.GetIntelExpressPartner() != nil {
		b, err := json.Marshal(storage.IntelExpressPartner{
			// todo: once we have parameters
			Param1:    pnr.IntelExpressPartner.Param1,
			Param2:    pnr.IntelExpressPartner.Param2,
			StartDate: pnr.IntelExpressPartner.StartDate.AsTime(),
			EndDate:   pnr.IntelExpressPartner.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}
		if _, err := s.st.UpdatePartner(ctx, &storage.Partner{
			ID:        pnr.IntelExpressPartner.ID,
			OrgID:     pnr.OrgID,
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
		}); err != nil {
			if err == storage.NotFound {
				logging.WithError(err, log).Error("partner doesn't exists")
				return status.Error(codes.NotFound, "failed to update partner record")
			}
			logging.WithError(err, log).Error("update partner")
			return status.Error(codes.Internal, "failed to update partner record")
		}
	}
	return nil
}
