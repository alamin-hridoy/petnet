package partner

import (
	"context"

	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/petnet/gunk/drp/v1/partner"
	pfSvc "brank.as/petnet/gunk/dsa/v2/service"

	spb "brank.as/petnet/gunk/dsa/v2/partner"
	pl "brank.as/petnet/gunk/dsa/v2/partnerlist"
)

func (s *Svc) RemitPartners(ctx context.Context, req *ppb.RemitPartnersRequest) (*ppb.RemitPartnersResponse, error) {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Country, is.CountryCode2),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	p, err := s.remit.ListPartners(ctx, req.GetCountry())
	if err != nil {
		return nil, status.Error(codes.Internal, "listing partners failed")
	}

	pt := map[string]*ppb.RemitPartner{}
	for _, r := range p {
		ptr := &ppb.RemitPartner{
			PartnerCode:            r.Code,
			PartnerName:            r.Name,
			SupportedSendTypes:     make(map[string]*ppb.RemitType, len(r.SendTypes)),
			SupportedDisburseTypes: make(map[string]*ppb.RemitType, len(r.DisburseTypes)),
		}
		for k, v := range r.SendTypes {
			ptr.SupportedSendTypes[k] = &ppb.RemitType{Code: k, Description: v.Description}
		}
		for k, v := range r.DisburseTypes {
			ptr.SupportedDisburseTypes[k] = &ppb.RemitType{Code: k, Description: v.Description}
		}
		pt[r.Code] = ptr
	}

	ep, err := s.enabledPartners(ctx, pt)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, "listing partners failed")
	}

	epg, err := s.enabledPartnersGlobally(ctx, ep)
	if err != nil {
		logging.WithError(err, log).Error("unable to get enabled Partners Globally")
		return nil, status.Error(codes.Internal, "listing partners failed")
	}

	return &ppb.RemitPartnersResponse{Partners: epg}, nil
}

func (s *Svc) enabledPartners(ctx context.Context, pt map[string]*ppb.RemitPartner) (map[string]*ppb.RemitPartner, error) {
	if phmw.GetEnv(ctx) == "sandbox" {
		return pt, nil
	}
	res, err := s.scl.ListServiceRequest(ctx, &pfSvc.ListServiceRequestRequest{
		OrgIDs:   []string{hydra.OrgID(ctx)},
		Statuses: []pfSvc.ServiceRequestStatus{pfSvc.ServiceRequestStatus_ACCEPTED},
	})

	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}

	ptrsts := make(map[string]bool)

	for _, v := range res.GetServiceRequst() {
		ptrsts[v.GetPartner()] = v.GetEnabled()
	}

	for k := range pt {
		v, ok := ptrsts[k]
		if !ok || (ok && !v) {
			delete(pt, k)
		}
	}

	return pt, nil
}

func (s *Svc) enabledPartnersGlobally(ctx context.Context, ep map[string]*ppb.RemitPartner) (map[string]*ppb.RemitPartner, error) {
	gpl, err := s.plcl.GetPartnerList(ctx, &pl.GetPartnerListRequest{
		Status: spb.PartnerStatusType_ENABLED.String(),
	})

	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}

	gplmap := make(map[string]bool)
	for _, z := range gpl.GetPartnerList() {
		gplmap[z.GetStype()] = true
	}

	for k := range ep {
		v, ok := gplmap[k]
		if !ok || (ok && !v) {
			delete(ep, k)
		}
	}
	return ep, nil
}
