package partner

import (
	"context"
	"encoding/json"

	spb "brank.as/petnet/gunk/dsa/v2/partner"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetPartner(ctx context.Context, oid string, tp string) (*storage.Partner, error) {
	return s.st.GetPartner(ctx, oid, tp)
}

func (s *Svc) GetPartners(ctx context.Context, oid string) (*spb.Partners, error) {
	log := logging.FromContext(ctx)

	ss, err := s.st.GetPartners(ctx, oid)
	if err != nil {
		if err == storage.NotFound {
			logging.WithError(err, log).Error("partner doesn't exists")
			return nil, status.Error(codes.NotFound, "failed to get partner record")
		}
		logging.WithError(err, log).Error("get partner")
		return nil, status.Error(codes.Internal, "failed to get partner record")
	}

	ps := &spb.Partners{}
	ps.PartnerStatuses = make(map[string]string)
	for _, pnr := range ss {
		switch spb.PartnerType(spb.PartnerType_value[pnr.Type]) {
		case spb.PartnerType_WU:
			sv := &storage.WesternUnionPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.WesternUnionPartner = &spb.WesternUnionPartner{
				ID:         pnr.ID,
				Coy:        sv.Coy,
				TerminalID: sv.TerminalID,
				Status:     spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:    tspb.New(pnr.Created),
				Updated:    tspb.New(pnr.Updated),
				StartDate:  tspb.New(sv.StartDate),
				EndDate:    tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_WU.String()] = pnr.Status
		case spb.PartnerType_IR:
			sv := &storage.IRemitPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.IRemitPartner = &spb.IRemitPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_IR.String()] = pnr.Status
		case spb.PartnerType_TF:
			sv := &storage.TransfastPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.TransfastPartner = &spb.TransfastPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_TF.String()] = pnr.Status
		case spb.PartnerType_RM:
			sv := &storage.RemitlyPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.RemitlyPartner = &spb.RemitlyPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_RM.String()] = pnr.Status
		case spb.PartnerType_RIA:
			sv := &storage.RiaPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.RiaPartner = &spb.RiaPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_RIA.String()] = pnr.Status
		case spb.PartnerType_MB:
			sv := &storage.MetroBankPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.MetroBankPartner = &spb.MetroBankPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_MB.String()] = pnr.Status
		case spb.PartnerType_BPI:
			sv := &storage.BPIPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.BPIPartner = &spb.BPIPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_BPI.String()] = pnr.Status
		case spb.PartnerType_USSC:
			sv := &storage.USSCPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.USSCPartner = &spb.USSCPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_USSC.String()] = pnr.Status
		case spb.PartnerType_JPR:
			sv := &storage.JapanRemitPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.JapanRemitPartner = &spb.JapanRemitPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_JPR.String()] = pnr.Status
		case spb.PartnerType_IC:
			sv := &storage.InstantCashPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.InstantCashPartner = &spb.InstantCashPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_IC.String()] = pnr.Status
		case spb.PartnerType_UNT:
			sv := &storage.UnitellerPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.UnitellerPartner = &spb.UnitellerPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_UNT.String()] = pnr.Status
		case spb.PartnerType_CEB:
			sv := &storage.CebuanaPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.CebuanaPartner = &spb.CebuanaPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_CEB.String()] = pnr.Status
		case spb.PartnerType_WISE:
			sv := &storage.TransferWisePartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.TransferWisePartner = &spb.TransferWisePartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_WISE.String()] = pnr.Status
		case spb.PartnerType_CEBI:
			sv := &storage.CebuanaIntlPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.CebuanaIntlPartner = &spb.CebuanaIntlPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_CEBI.String()] = pnr.Status
		case spb.PartnerType_AYA:
			sv := &storage.AyannahPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.AyannahPartner = &spb.AyannahPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_AYA.String()] = pnr.Status
		case spb.PartnerType_IE:
			sv := &storage.IntelExpressPartner{}
			ps.UpdatedBy = pnr.UpdatedBy
			if err := json.Unmarshal([]byte(pnr.Partner), sv); err != nil {
				logging.WithError(err, log).Error("unmarshaling partner")
				return nil, status.Error(codes.Internal, "failed to get partner record")
			}
			ps.IntelExpressPartner = &spb.IntelExpressPartner{
				ID:        pnr.ID,
				Param1:    sv.Param1,
				Param2:    sv.Param2,
				Status:    spb.PartnerStatusType(spb.PartnerStatusType_value[pnr.Status]),
				Created:   tspb.New(pnr.Created),
				Updated:   tspb.New(pnr.Updated),
				StartDate: tspb.New(sv.StartDate),
				EndDate:   tspb.New(sv.EndDate),
			}
			ps.PartnerStatuses[spb.PartnerType_IE.String()] = pnr.Status
		}
	}
	if len(ps.PartnerStatuses) == 0 {
		log.Error("no partners found")
		return nil, status.Error(codes.NotFound, "service not found")
	}
	ps.OrgID = oid
	return ps, nil
}
