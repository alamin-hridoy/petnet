package product

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/rbac/serviceutil/logging"

	ppb "brank.as/rbac/gunk/v1/permissions"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) ListServiceAssignments(ctx context.Context, req *ppb.ListServiceAssignmentsRequest) (*ppb.ListServiceAssignmentsResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.product.listservices")

	if err := validation.ValidateStruct(req,
		validation.Field(req.OrgID, is.UUIDv4),
		validation.Field(req.ServiceID, validation.Each(is.UUIDv4)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var assn []storage.ServiceAssignment
	if req.OrgID != "" {
		a, err := s.st.ListServiceAssignOrg(ctx, req.OrgID)
		if err != nil {
			logging.WithError(err, log).Error("getting service asssignments")
			return nil, status.Error(codes.Internal, "failed to list services")
		}
		assn = a
	} else {
		a, err := s.st.ListServiceAssign(ctx)
		if err != nil {
			logging.WithError(err, log).Error("getting service asssignments")
			return nil, status.Error(codes.Internal, "failed to list services")
		}
		assn = a
	}

	svcs, err := s.st.ListService(ctx)
	if err != nil {
		logging.WithError(err, log).Error("getting services")
		return nil, status.Error(codes.Internal, "failed to list services")
	}
	sm := make(map[string]storage.Service, len(svcs))
	for _, sv := range svcs {
		sm[sv.ID] = sv
	}

	if len(req.ServiceID) != 0 {
		m := map[string]bool{}
		for _, s := range req.ServiceID {
			m[s] = true
		}

		asn := make([]storage.ServiceAssignment, 0, len(assn))
		for _, a := range assn {
			if m[a.ServiceID] {
				asn = append(asn, a)
			}
		}
		assn = asn
	}

	svc := make([]*ppb.ServiceAssignment, len(assn))
	for i, a := range assn {
		svc[i] = &ppb.ServiceAssignment{
			Grant:       a.GrantID,
			OrgID:       a.OrgID,
			ServiceID:   a.ServiceID,
			ServiceName: sm[a.ServiceID].Name,
			Environment: a.Environment,
			GrantedBy:   a.AssignUserID,
			Granted:     tspb.New(a.Assigned),
			RevokedBy:   a.RevokeUserID.String,
			Revoked:     tspb.New(a.Revoked.Time),
		}
	}

	return &ppb.ListServiceAssignmentsResponse{
		ServiceAssignments: svc,
	}, nil
}
