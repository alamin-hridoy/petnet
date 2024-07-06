package microinsurance

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Transact ...
func (c *Client) Transact(ctx context.Context, req *TransactRequest) (*Insurance, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	rawRes, err := c.phService.PostMicroInsurance(ctx, c.getUrl("transact"), *req)

	return toInsuranceResponse(rawRes, err)
}

// GetReprint ...
func (c *Client) GetReprint(ctx context.Context, req *GetReprintRequest) (*Insurance, error) {
	if req == nil || req.TraceNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "trace number is required")
	}

	rawRes, err := c.phService.PostMicroInsurance(ctx, c.getUrl("get-reprint"), *req)

	return toInsuranceResponse(rawRes, err)
}

// RetryTransaction ...
func (c *Client) RetryTransaction(ctx context.Context, req *RetryTransactionRequest) (*Insurance, error) {
	if req == nil || req.ID == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	rawRes, err := c.phService.PostMicroInsurance(ctx, c.getUrl("retry"), *req)

	return toInsuranceResponse(rawRes, err)
}

// GetTransactionList ...
func (c *Client) GetTransactionList(ctx context.Context, req *GetTransactionListRequest) (*TransactionListResult, error) {
	if req == nil || req.DateFrom.IsZero() {
		return nil, status.Error(codes.InvalidArgument, "date_from is required")
	}

	if req.DateTo.IsZero() {
		return nil, status.Error(codes.InvalidArgument, "date_to is required")
	}

	rawRes, err := c.phService.PostMicroInsurance(ctx, c.getUrl("get-transactions-list"), *req)
	if err != nil {
		return nil, err
	}

	var res GetTransactionListResult
	err = json.Unmarshal(rawRes, &res)
	if err != nil || res.Result == nil {
		return nil, status.Error(codes.Internal, "Invalid perahub response")
	}

	return res.Result, nil
}

func toInsuranceResponse(jsonRaw json.RawMessage, err error) (*Insurance, error) {
	if err != nil {
		return nil, err
	}

	var resp TransactResult
	err = json.Unmarshal(jsonRaw, &resp)
	if err != nil || resp.Result == nil {
		return nil, status.Error(codes.Internal, "Invalid perahub response")
	}

	return resp.Result, nil
}
