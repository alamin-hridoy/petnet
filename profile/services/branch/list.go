package branch

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	bpb "brank.as/petnet/gunk/dsa/v2/branch"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
)

func (s *Svc) ListBranches(ctx context.Context, req *bpb.ListBranchesRequest) (*bpb.ListBranchesResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Limit, validation.Min(0)),
		validation.Field(&req.Offset, validation.Min(0)),
		validation.Field(&req.OrgID, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	br, err := s.core.ListBranches(ctx, req.GetOrgID(), int(req.GetLimit()), int(req.GetOffset()), req.GetTitle())
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to retrieve branch records")
	}

	ts := func(t time.Time) *tspb.Timestamp {
		if t.IsZero() {
			return nil
		}
		return tspb.New(t)
	}
	brs := make([]*bpb.Branch, len(br))
	for i, b := range br {
		brs[i] = &bpb.Branch{
			ID:    b.ID,
			OrgID: b.OrgID,
			Title: b.Title,
			Address: &ppb.Address{
				Address1:   b.Address1,
				City:       b.City,
				State:      b.State,
				PostalCode: b.PostalCode,
			},
			PhoneNumber:   b.PhoneNumber,
			FaxNumber:     b.FaxNumber,
			ContactPerson: b.ContactPerson,
			Created:       ts(b.Created),
			Updated:       ts(b.Updated),
			Deleted:       ts(b.Deleted.Time),
		}
	}
	var tot int32
	if len(br) > 0 {
		tot = int32(br[0].Count)
	}
	return &bpb.ListBranchesResponse{
		Branches: brs,
		Total:    tot,
	}, nil
}
