package service

import (
	"context"
	"strconv"

	spb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListServiceRequest(ctx context.Context, req *spb.ListServiceRequestRequest) (*spb.ListServiceRequestResponse, error) {
	log := logging.FromContext(ctx)

	// todo clean this up
	var ss []string
	for _, st := range req.Statuses {
		ss = append(ss, st.String())
	}
	var ts []string
	for _, t := range req.Types {
		ts = append(ts, t.String())
	}
	var ps []string
	for _, p := range req.Partners {
		ps = append(ps, p)
	}

	rs, err := s.st.ListSvcRequest(ctx, storage.SvcRequestFilter{
		OrgID:        req.OrgIDs,
		Status:       ss,
		SvcName:      ts,
		Partner:      ps,
		SortByColumn: req.SortByColumn.String(),
		SortOrder:    req.SortBy.String(),
		Limit:        int(req.Limit),
		Offset:       int(req.Offset),
	})
	if err != nil {
		logging.WithError(err, log).Error("validate request")
		return nil, err
	}
	res := &spb.ListServiceRequestResponse{}
	for _, r := range rs {
		rr := &spb.ServiceRequest{
			ID:          r.ID,
			OrgID:       r.OrgID,
			CompanyName: r.CompanyName,
			Partner:     r.Partner,
			Type:        spb.ServiceType(spb.ServiceType_value[r.SvcName]),
			Status:      spb.ServiceRequestStatus(spb.ServiceRequestStatus_value[r.Status]),
			Enabled:     r.Enabled,
			Remarks:     r.Remarks,
			Created:     tspb.New(r.Created),
			Updated:     tspb.New(r.Updated),
			UpdatedBy:   r.UpdatedBy,
		}
		if r.Applied.Valid {
			rr.Applied = tspb.New(r.Applied.Time)
		}
		res.ServiceRequst = append(res.ServiceRequst, rr)
	}
	if len(rs) > 0 {
		total, err := strconv.Atoi(rs[0].Total)
		if err != nil {
			logging.WithError(err, log).Error("Cant convert total")
			return nil, err
		}
		res.Total = int32(total)
	}
	return res, nil
}
