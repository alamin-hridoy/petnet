package product

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

func (s *Svc) ListServices(ctx context.Context, req *ppb.ListServicesRequest) (*ppb.ListServicesResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.product.listservices")

	svcs, err := s.st.ListService(ctx)
	if err != nil {
		logging.WithError(err, log).Error("fetch services")
		return nil, status.Error(codes.Internal, "failed to list services")
	}

	pub, err := s.st.ListServicePublic(ctx)
	if err != nil {
		logging.WithError(err, log).Error("fetch public services")
		return nil, status.Error(codes.Internal, "failed to list services")
	}
	log.WithField("pub", pub).Info("public services")
	m := map[string][]string{}
	for _, p := range pub {
		if p.Retracted.Valid {
			continue
		}
		e := p.Environment
		if e == "" {
			e = "All"
		}
		if len(m[p.ServiceID]) == 0 {
			m[p.ServiceID] = []string{e}
			continue
		}
		m[p.ServiceID] = append(m[p.ServiceID], e)
	}

	svc := make([]*ppb.Service, len(svcs))
	pbs := make([]*ppb.Service, 0, len(pub))
	for i, sc := range svcs {
		svc[i] = &ppb.Service{
			ServiceID:   sc.ID,
			ServiceName: sc.Name,
			Description: sc.Description,
		}
		for _, e := range m[sc.ID] {
			pbs = append(pbs, &ppb.Service{
				ServiceID:   sc.ID,
				ServiceName: sc.Name,
				Description: sc.Description,
				Environment: e,
			})
		}

	}

	return &ppb.ListServicesResponse{
		Services: svc,
		Public:   pbs,
	}, nil
}
