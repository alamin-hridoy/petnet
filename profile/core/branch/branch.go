package branch

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw/meta/md"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Svc struct {
	st *postgres.Storage
}

func New(st *postgres.Storage) *Svc {
	return &Svc{st: st}
}

func (s *Svc) UpsertBranch(ctx context.Context, b storage.Branch) (*storage.Branch, error) {
	log := logging.FromContext(ctx)

	if b.OrgProfileID == "" {
		switch "" {
		case md.GetProfileID(ctx):
			b.OrgProfileID = md.GetProfileID(ctx)
		case hydra.OrgID(ctx):
			pr, err := s.st.GetProfileID(ctx, hydra.OrgID(ctx))
			if err != nil {
				logging.WithError(err, log).Error("no profile id found")
				return nil, status.Error(codes.NotFound, "org not found")
			}
			b.OrgProfileID = pr
		default:
			return nil, status.Error(codes.InvalidArgument, "no org identified")
		}
	}
	br, err := s.st.UpsertBranch(ctx, b)
	if err != nil {
		logging.WithError(err, log).Error("store branch")
		return nil, status.Error(codes.Internal, "failed to record branch record")
	}
	return br, nil
}

func (s *Svc) ListBranches(ctx context.Context, org string, lim, off int, title string) ([]storage.Branch, error) {
	log := logging.FromContext(ctx)

	br, err := s.st.ListBranches(ctx, org, storage.LimitOffsetFilter{
		Limit:  int32(lim),
		Offset: int32(off),
	}, title)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "branch records not found")
		}
		logging.WithError(err, log).Error("fetch org branches")
		return nil, status.Error(codes.Internal, "fetch branch records failed")
	}
	return br, nil
}
