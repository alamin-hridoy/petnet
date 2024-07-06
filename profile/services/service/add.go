package service

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/cms/mw"
	ptnr "brank.as/petnet/gunk/dsa/v2/partner"
	spb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func (s *Svc) AddServiceRequest(ctx context.Context, req *spb.AddServiceRequestRequest) (*empty.Empty, error) {
	required := validation.Required
	remit := s.allPartners(ctx, spb.ServiceType_REMITTANCE)
	cico := s.allPartners(ctx, spb.ServiceType_CASHINCASHOUT)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.Type, required),
		validation.Field(&req.Partners,
			validation.When(req.Type == spb.ServiceType_REMITTANCE || req.Type == spb.ServiceType_CASHINCASHOUT && !req.AllPartners),
			required, validation.By(func(interface{}) error {
				for _, p := range req.Partners {
					if _, ok := remit[p]; !ok && req.Type == spb.ServiceType_REMITTANCE {
						return fmt.Errorf("partner: %v, is not a valid partner for service %v", p, spb.ServiceType_REMITTANCE)
					}
					if _, ok := cico[p]; !ok && req.Type == spb.ServiceType_CASHINCASHOUT {
						return fmt.Errorf("cico partner: %v, is not a valid cico partner for service %v", p, spb.ServiceType_CASHINCASHOUT)
					}
				}
				return nil
			})),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return s.doSubmitSvcRequest(ctx, req)
}

func (s *Svc) doSubmitSvcRequest(ctx context.Context, req *spb.AddServiceRequestRequest) (*empty.Empty, error) {
	log := logging.FromContext(ctx)
	res, err := s.st.GetOrgProfile(ctx, req.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("getting org profile")
		return nil, err
	}
	_, err = s.getListAndRemoveSvcRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	for _, p := range req.Partners {
		if err := s.createSvcRequest(ctx, storage.ServiceRequest{
			OrgID:       req.OrgID,
			Partner:     p,
			SvcName:     req.Type.String(),
			CompanyName: res.CompanyName,
		}); err != nil {
			return nil, err
		}
	}
	return new(empty.Empty), nil
}

func (s *Svc) getListAndRemoveSvcRequest(ctx context.Context, req *spb.AddServiceRequestRequest) ([]storage.ServiceRequest, error) {
	log := logging.FromContext(ctx)
	lRes, err := s.st.ListSvcRequest(ctx, storage.SvcRequestFilter{
		OrgID:   []string{req.OrgID},
		SvcName: []string{req.Type.String()},
	})
	if err != nil {
		return nil, err
	}
	for _, lP := range lRes {
		if lP.Status == spb.ServiceRequestStatus_ACCEPTED.String() || lP.Status == spb.ServiceRequestStatus_PENDING.String() {
			continue
		}
		av, _ := mw.InArray(lP.Partner, req.Partners)
		if !av {
			err := s.st.RemoveSvcRequest(ctx, storage.ServiceRequest{
				OrgID:   req.OrgID,
				Partner: lP.Partner,
				SvcName: req.Type.String(),
			})
			if err != nil {
				logging.WithError(err, log).Error("Remove Service Request Failed")
			}
		}
	}
	return lRes, nil
}

func (s *Svc) createSvcRequest(ctx context.Context, req storage.ServiceRequest) error {
	log := logging.FromContext(ctx)
	if _, err := s.st.CreateSvcRequest(ctx, storage.ServiceRequest{
		OrgID:       req.OrgID,
		Partner:     req.Partner,
		SvcName:     req.SvcName,
		CompanyName: req.CompanyName,
	}); err != nil {
		if err != storage.Conflict {
			logging.WithError(err, log).Error("create Svc Request failed")
			return err
		}
	}
	return nil
}

func acceptedAndPendingPartners(sReq []storage.ServiceRequest) (res []string) {
	if len(sReq) > 0 {
		for _, v := range sReq {
			if v.Status == spb.ServiceRequestStatus_ACCEPTED.String() || v.Status == spb.ServiceRequestStatus_PENDING.String() {
				res = append(res, v.Partner)
			}
		}
	}
	return
}

func (s *Svc) allPartners(ctx context.Context, t spb.ServiceType) map[string]string {
	ps := map[string]string{}
	res, err := s.st.GetPartnerList(ctx, &storage.PartnerList{
		ServiceName: t.String(),
		Status:      ptnr.PartnerStatusType_ENABLED.String(),
	})
	if err != nil || res == nil {
		return ps
	}

	for _, v := range res {
		ps[v.Stype] = v.Name
	}

	if len(ps) > 0 {
		return ps
	}

	return allStaticPartners(t)
}

func allStaticPartners(t spb.ServiceType) map[string]string {
	ps := map[string]string{}
	switch t {
	case spb.ServiceType_REMITTANCE:
		for k := range spb.RemittancePartner_value {
			ps[k] = k
		}
	case spb.ServiceType_CASHINCASHOUT:
		for k := range spb.CICOPartner_value {
			ps[k] = k
		}
	}
	return ps
}
