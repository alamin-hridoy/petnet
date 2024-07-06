package branch

import (
	"context"
	"database/sql"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bpb "brank.as/petnet/gunk/dsa/v2/branch"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) UpsertBranch(ctx context.Context, req *bpb.UpsertBranchRequest) (*bpb.UpsertBranchResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Branch, validation.Required, validation.By(func(interface{}) error {
			r := req.GetBranch()
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.ID, is.UUIDv4),
			)
		})),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := req.GetBranch()

	br, err := s.core.UpsertBranch(ctx, storage.Branch{
		ID:           r.GetID(),
		OrgID:        r.GetOrgID(),
		OrgProfileID: "10000000-0000-0000-0000-000000000000", // OrgProfileID will be deleted once we refactored profile, setting to avoid error for now
		Title:        r.GetTitle(),
		BranchAddress: storage.BranchAddress{
			Address1:   r.GetAddress().GetAddress1(),
			City:       r.GetAddress().GetCity(),
			State:      r.GetAddress().GetState(),
			PostalCode: r.GetAddress().GetPostalCode(),
		},
		PhoneNumber:   r.GetPhoneNumber(),
		FaxNumber:     r.GetFaxNumber(),
		ContactPerson: r.GetContactPerson(),
		Updated:       time.Now(),
		Deleted: sql.NullTime{
			Time:  r.GetDeleted().AsTime(),
			Valid: r.GetDeleted().IsValid(),
		},
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to store branch record")
	}

	return &bpb.UpsertBranchResponse{
		ID: br.ID,
	}, nil
}
