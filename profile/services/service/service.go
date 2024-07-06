package service

import (
	"context"
	"strings"

	spb "brank.as/petnet/gunk/dsa/v2/service"
	eml "brank.as/petnet/profile/integrations/email"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type iEmailStore interface {
	SendDsaServiceRequestNotification(ctx context.Context, req eml.DsaServiceRequestNotificationForm) error
}

type Svc struct {
	spb.UnimplementedServiceServiceServer
	st         *postgres.Storage
	emailStore iEmailStore
}

func New(st *postgres.Storage, emailStore iEmailStore) *Svc {
	return &Svc{
		st:         st,
		emailStore: emailStore,
	}
}

// RegisterPartner with grpc server.
func (s *Svc) Register(srv *grpc.Server) { spb.RegisterServiceServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return spb.RegisterServiceServiceHandlerFromEndpoint(ctx, mux, address, options)
}

func (s *Svc) sendServiceNotification(ctx context.Context, orgID string, svc string) error {
	log := logging.FromContext(ctx).WithField("method", "service.sendServiceNotification")
	svcReq, err := s.st.GetAllServiceRequest(ctx, storage.SvcRequestFilter{
		OrgID:   []string{orgID},
		SvcName: []string{svc},
	})
	if err != nil {
		logging.WithError(err, log).Error("unable to get all service request")
		return err
	}

	if svcReq == nil || len(svcReq) == 0 {
		err = status.Error(codes.NotFound, "service request not found.")
		logging.WithError(err, log).Error("service request not found.")
		return err
	}

	totalReq := len(svcReq)
	acceptedPtnrs, rejectedPtnrs := []string{}, []string{}
	acceptedRemark, rejectedRemark, email := "", "", ""
	for _, svcR := range svcReq {
		if svcR.Status == spb.ServiceRequestStatus_ACCEPTED.String() {
			acceptedPtnrs = append(acceptedPtnrs, svcR.Partner)
			acceptedRemark = svcR.Remarks
		}
		if svcR.Status == spb.ServiceRequestStatus_REJECTED.String() {
			rejectedPtnrs = append(rejectedPtnrs, svcR.Partner)
			rejectedRemark = svcR.Remarks
		}
	}

	if acceptedRemark == "" {
		acceptedRemark = "Accepted"
	}

	if rejectedRemark == "" {
		rejectedRemark = "Rejected"
	}

	email = s.getEmailByOrgid(ctx, orgID)
	if email == "" {
		err = status.Error(codes.NotFound, "org email not found.")
		logging.WithError(err, log).Error("org email not found.")
		return err
	}

	AcptRjtCnt := len(acceptedPtnrs) + len(rejectedPtnrs)

	if totalReq == AcptRjtCnt && len(acceptedPtnrs) > 0 {
		ptnrs := s.getPartnersByStypes(ctx, acceptedPtnrs)
		if len(ptnrs) == 0 {
			err = status.Error(codes.NotFound, "partner not found.")
			logging.WithError(err, log).Error("partner not found.")
		}
		if len(ptnrs) > 0 {
			s.emailStore.SendDsaServiceRequestNotification(ctx, eml.DsaServiceRequestNotificationForm{
				Email:        email,
				Status:       spb.ServiceRequestStatus_ACCEPTED.String(),
				ServiceName:  svc,
				Remark:       acceptedRemark,
				PartnerNames: strings.Join(ptnrs, ", "),
			})
		}
	}

	if totalReq == AcptRjtCnt && len(rejectedPtnrs) > 0 {
		ptnrs := s.getPartnersByStypes(ctx, rejectedPtnrs)
		if len(ptnrs) == 0 {
			err = status.Error(codes.NotFound, "partner not found.")
			logging.WithError(err, log).Error("partner not found.")
		}
		if len(ptnrs) > 0 {
			s.emailStore.SendDsaServiceRequestNotification(ctx, eml.DsaServiceRequestNotificationForm{
				Email:        email,
				Status:       spb.ServiceRequestStatus_REJECTED.String(),
				ServiceName:  svc,
				Remark:       rejectedRemark,
				PartnerNames: strings.Join(ptnrs, ", "),
			})
		}
	}

	return nil
}

func (s *Svc) getPartnersByStypes(ctx context.Context, stype []string) []string {
	partners := []string{}
	if mw.InSlice("RuralNet", stype) {
		stype = append(stype, "RLN")
	}
	ptnrs, err := s.st.GetPartnerList(ctx, &storage.PartnerList{
		Stype: strings.Join(stype, ","),
	})
	if err != nil {
		return partners
	}

	if ptnrs == nil || len(ptnrs) == 0 {
		return partners
	}

	for _, p := range ptnrs {
		partners = append(partners, p.Name)
	}

	return partners
}

func (s *Svc) getEmailByOrgid(ctx context.Context, orgID string) string {
	email := ""
	usrs, err := s.st.GetUserProfiles(ctx, orgID)
	if err != nil {
		return email
	}

	if usrs == nil || len(usrs) == 0 {
		return email
	}

	for _, u := range usrs {
		email = u.Email
	}

	return email
}
