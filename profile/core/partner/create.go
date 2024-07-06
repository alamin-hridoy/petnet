package partner

import (
	"context"
	"encoding/json"

	ppb "brank.as/petnet/gunk/dsa/v2/partner"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CreatePartners(ctx context.Context, pnr *ppb.Partners) error {
	log := logging.FromContext(ctx)
	if sv := pnr.GetWesternUnionPartner(); sv != nil {
		b, err := json.Marshal(storage.WesternUnionPartner{
			Coy:        sv.Coy,
			TerminalID: sv.TerminalID,
			UpdatedBy:  pnr.UpdatedBy,
			Status:     sv.GetStatus().String(),
			Created:    sv.Created.AsTime(),
			Updated:    sv.Updated.AsTime(),
			StartDate:  sv.StartDate.AsTime(),
			EndDate:    sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_WU.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.WesternUnionPartner.ID = id
	}
	if sv := pnr.GetIRemitPartner(); sv != nil {
		b, err := json.Marshal(storage.IRemitPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_IR.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.IRemitPartner.ID = id
	}
	if sv := pnr.GetTransfastPartner(); sv != nil {
		b, err := json.Marshal(storage.TransfastPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_TF.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.TransfastPartner.ID = id
	}
	if sv := pnr.GetRemitlyPartner(); sv != nil {
		b, err := json.Marshal(storage.RemitlyPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_RM.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.RemitlyPartner.ID = id
	}
	if sv := pnr.GetRiaPartner(); sv != nil {
		b, err := json.Marshal(storage.RiaPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_RIA.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.RiaPartner.ID = id
	}
	if sv := pnr.GetMetroBankPartner(); sv != nil {
		b, err := json.Marshal(storage.MetroBankPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_MB.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.MetroBankPartner.ID = id
	}
	if sv := pnr.GetBPIPartner(); sv != nil {
		b, err := json.Marshal(storage.BPIPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_BPI.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.BPIPartner.ID = id
	}
	if sv := pnr.GetUSSCPartner(); sv != nil {
		b, err := json.Marshal(storage.USSCPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_USSC.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.USSCPartner.ID = id
	}
	if sv := pnr.GetJapanRemitPartner(); sv != nil {
		b, err := json.Marshal(storage.JapanRemitPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_JPR.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.JapanRemitPartner.ID = id
	}
	if sv := pnr.GetInstantCashPartner(); sv != nil {
		b, err := json.Marshal(storage.InstantCashPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_IC.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.InstantCashPartner.ID = id
	}
	if sv := pnr.GetUnitellerPartner(); sv != nil {
		b, err := json.Marshal(storage.UnitellerPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_UNT.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.UnitellerPartner.ID = id
	}
	if sv := pnr.GetCebuanaPartner(); sv != nil {
		b, err := json.Marshal(storage.CebuanaPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_CEB.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.CebuanaPartner.ID = id
	}
	if sv := pnr.GetTransferWisePartner(); sv != nil {
		b, err := json.Marshal(storage.TransferWisePartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_WISE.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.TransferWisePartner.ID = id
	}
	if sv := pnr.GetCebuanaIntlPartner(); sv != nil {
		b, err := json.Marshal(storage.CebuanaIntlPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_CEBI.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.CebuanaIntlPartner.ID = id
	}
	if sv := pnr.GetAyannahPartner(); sv != nil {
		b, err := json.Marshal(storage.AyannahPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_AYA.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.AyannahPartner.ID = id
	}
	if sv := pnr.GetIntelExpressPartner(); sv != nil {
		b, err := json.Marshal(storage.IntelExpressPartner{
			Param1:    sv.Param1,
			Param2:    sv.Param2,
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
			Created:   sv.Created.AsTime(),
			Updated:   sv.Updated.AsTime(),
			StartDate: sv.StartDate.AsTime(),
			EndDate:   sv.EndDate.AsTime(),
		})
		if err != nil {
			logging.WithError(err, log).Error("marshaling json")
			return status.Error(codes.Internal, "failed to update partner record")
		}

		id, err := s.st.CreatePartner(ctx, &storage.Partner{
			OrgID:     pnr.OrgID,
			Type:      ppb.PartnerType_IE.String(),
			Partner:   string(b),
			UpdatedBy: pnr.UpdatedBy,
			Status:    sv.GetStatus().String(),
		})
		if err != nil {
			if err == storage.Conflict {
				logging.WithError(err, log).Error("partner already exists")
				return status.Error(codes.AlreadyExists, "failed to record partner record")
			}
			logging.WithError(err, log).Error("store partner")
			return status.Error(codes.Internal, "failed to record partner record")
		}
		pnr.IntelExpressPartner.ID = id
	}
	return nil
}
