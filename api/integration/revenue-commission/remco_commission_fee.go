package revenue_commission

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/serviceutil/logging"
)

// RemcoCommissionFee is commission fee details for remco
type RemcoCommissionFee struct {
	FeeID               json.Number    `json:"id"`
	RemcoID             json.Number    `json:"remco_id"`
	MinAmount           json.Number    `json:"min_amount"`
	MaxAmount           json.Number    `json:"max_amount"`
	ServiceFee          json.Number    `json:"service_fee"`
	CommissionAmount    json.Number    `json:"commission_amount"`
	CommissionAmountOTC json.Number    `json:"commission_amount_otc"`
	CommissionType      CommissionType `json:"commission_type"`
	TrxType             TrxType        `json:"trx_type"`
	UpdatedBy           string         `json:"updated_by"`
	CreatedAt           *time.Time     `json:"created_at"`
	UpdatedAt           *time.Time     `json:"updated_at"`
}

// SaveRemcoCommissionFeeRequest is request for save new remco commission fee
type SaveRemcoCommissionFeeRequest struct {
	// Ignore fee id in request body
	FeeID               uint32         `json:"-"`
	RemcoID             json.Number    `json:"remco_id"`
	MinAmount           json.Number    `json:"min_amount"`
	MaxAmount           json.Number    `json:"max_amount"`
	ServiceFee          json.Number    `json:"service_fee"`
	CommissionType      CommissionType `json:"commission_type"`
	CommissionAmount    json.Number    `json:"commission_amount"`
	CommissionAmountOTC json.Number    `json:"commission_amount_otc"`
	TrxType             TrxType        `json:"trx_type"`
	UpdatedBy           string         `json:"updated_by"`
}

// RemcoCommissionFeeResult result object returned by remco commission fee apis
type RemcoCommissionFeeResult struct {
	Code    json.Number         `json:"code"`
	Message string              `json:"message"`
	Result  *RemcoCommissionFee `json:"result"`
}

// ListRemcoCommissionFeeResult result object returned by get all remco commission fee api
type ListRemcoCommissionFeeResult struct {
	Code    json.Number          `json:"code"`
	Message string               `json:"message"`
	Result  []RemcoCommissionFee `json:"result"`
}

// ListRemcoCommissionFee returns all remco commission fee list
func (c *Client) ListRemcoCommissionFee(ctx context.Context) ([]RemcoCommissionFee, error) {
	log := logging.FromContext(ctx)
	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl("remco-sf"))
	if err != nil {
		log.WithError(err).Error("list remco commission fee perahub error")
		return nil, coreerror.ToCoreError(err)
	}

	var resp ListRemcoCommissionFeeResult
	err = json.Unmarshal(rawRes, &resp)
	if err != nil {
		log.WithError(err).Error("list remco commission fee unmarshal error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return resp.Result, nil
}

// GetRemcoCommissionFeeByID returns remco commission fee by fee ID
func (c *Client) GetRemcoCommissionFeeByID(ctx context.Context, feeID uint32) (*RemcoCommissionFee, error) {
	log := logging.FromContext(ctx)
	if feeID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "remco commission fee id required")
	}

	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl(fmt.Sprintf("remco-sf/%d", feeID)))

	return toRemcoCommissionFeeResponse(rawRes, err, log)
}

// CreateRemcoCommissionFee create new remco commission fee
func (c *Client) CreateRemcoCommissionFee(ctx context.Context, request *SaveRemcoCommissionFeeRequest) (*RemcoCommissionFee, error) {
	log := logging.FromContext(ctx)
	if request == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	rawRes, err := c.phService.PostRevComm(ctx, c.getUrl("remco-sf"), *request)

	return toRemcoCommissionFeeResponse(rawRes, err, log)
}

// UpdateRemcoCommissionFee update remco commission fee by fee id and request
func (c *Client) UpdateRemcoCommissionFee(ctx context.Context, request *SaveRemcoCommissionFeeRequest) (*RemcoCommissionFee, error) {
	log := logging.FromContext(ctx)
	if request == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	if request.FeeID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "remco commission fee id required")
	}

	rawRes, err := c.phService.PutRevComm(ctx, c.getUrl(fmt.Sprintf("remco-sf/%d", request.FeeID)), *request)

	return toRemcoCommissionFeeResponse(rawRes, err, log)
}

// DeleteRemcoCommissionFee delete remco commission fee by fee id
func (c *Client) DeleteRemcoCommissionFee(ctx context.Context, feeID uint32) (*RemcoCommissionFee, error) {
	log := logging.FromContext(ctx)
	if feeID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "remco commission fee id required")
	}

	rawRes, err := c.phService.DeleteRevComm(ctx, c.getUrl(fmt.Sprintf("remco-sf/%d", feeID)))

	return toRemcoCommissionFeeResponse(rawRes, err, log)
}

func toRemcoCommissionFeeResponse(jsonRaw json.RawMessage, err error, log *logrus.Entry) (*RemcoCommissionFee, error) {
	if err != nil {
		log.WithError(err).Error("remco commission fee perahub error")
		return nil, coreerror.ToCoreError(err)
	}

	var resp RemcoCommissionFeeResult
	err = json.Unmarshal(jsonRaw, &resp)
	if err != nil {
		log.WithError(err).Error("remco commission fee unmarshal error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return resp.Result, nil
}
