package microinsurance

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
)

// GetProduct ...
func (c *Client) GetProduct(ctx context.Context, req *GetProductRequest) (*ProductResult, error) {
	if req == nil || req.ProductCode == "" {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "product code is required")
	}

	rawRes, err := c.phService.PostMicroInsurance(ctx, c.getUrl("get-product"), *req)
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var res GetProductResult
	err = json.Unmarshal(rawRes, &res)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, "Invalid perahub response")
	}

	return res.Result, nil
}

// GetOfferProduct ...
func (c *Client) GetOfferProduct(ctx context.Context, req *GetOfferProductRequest) (*OfferProduct, error) {
	if req == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	rawRes, err := c.phService.PostMicroInsurance(ctx, c.getUrl("get-offer-product"), *req)
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var res GetOfferProductResult
	err = json.Unmarshal(rawRes, &res)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, "Invalid perahub response")
	}

	return res.Result, nil
}

// CheckActiveProduct ...
func (c *Client) CheckActiveProduct(ctx context.Context, req *CheckActiveProductRequest) (*ActiveProduct, error) {
	if req == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	rawRes, err := c.phService.PostMicroInsurance(ctx, c.getUrl("check-active-product"), *req)
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var res CheckActiveProductResult
	err = json.Unmarshal(rawRes, &res)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, "Invalid perahub response")
	}

	return res.Result, nil
}

// GetProductList ...
func (c *Client) GetProductList(ctx context.Context) ([]ActiveProduct, error) {
	rawRes, err := c.phService.GetMicroInsurance(ctx, c.getUrl("product-code-list"))
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var res GetProductListResult
	err = json.Unmarshal(rawRes, &res)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, "Invalid perahub response")
	}

	return res.Result, nil
}
