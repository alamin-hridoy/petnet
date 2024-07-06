package microinsurance

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"brank.as/petnet/api/util"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
)

// GetProduct ...
func (s *Svc) GetProduct(ctx context.Context, req *migunk.GetProductRequest) (*migunk.ProductResult, error) {
	if req == nil || req.ProductCode == "" {
		return nil, status.Error(codes.InvalidArgument, "product code is required")
	}

	res, err := s.store.GetProduct(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

// GetOfferProduct ...
func (s *Svc) GetOfferProduct(ctx context.Context, req *migunk.GetOfferProductRequest) (*migunk.OfferProduct, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	err := validation.ValidateStruct(req,
		validation.Field(&req.FirstName, validation.Required),
		validation.Field(&req.LastName, validation.Required),
		validation.Field(&req.Birthdate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&req.Gender, validation.Required),
		validation.Field(&req.TrxType, validation.Required),
		validation.Field(&req.Amount, validation.Required),
	)
	if err != nil {
		return nil, err
	}

	res, err := s.store.GetOfferProduct(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

// CheckActiveProduct ...
func (s *Svc) CheckActiveProduct(ctx context.Context, req *migunk.CheckActiveProductRequest) (*migunk.ActiveProduct, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	err := validation.ValidateStruct(req,
		validation.Field(&req.FirstName, validation.Required),
		validation.Field(&req.LastName, validation.Required),
		validation.Field(&req.Birthdate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&req.Gender, validation.Required),
		validation.Field(&req.ProductCode, validation.Required),
	)
	if err != nil {
		return nil, err
	}

	res, err := s.store.CheckActiveProduct(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

// GetProductList ...
func (s *Svc) GetProductList(ctx context.Context, r *emptypb.Empty) (*migunk.GetProductListResult, error) {
	res, err := s.store.GetProductList(ctx)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}
