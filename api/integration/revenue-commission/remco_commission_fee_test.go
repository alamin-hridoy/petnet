package revenue_commission

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
)

func TestClient_ListRemcoCommissionFee(t *testing.T) {
	basePath := "http://some-url.test"
	baseUrl, err := url.Parse(basePath)
	require.Nil(t, err)
	remcoFee := getMockRemcoCommissionFee("1")

	for _, tc := range []struct {
		name          string
		serviceResp   *ListRemcoCommissionFeeResult
		serviceErr    error
		expectedURL   string
		expectedResp  []RemcoCommissionFee
		expectedError error
	}{
		{
			name: "happy path",
			serviceResp: &ListRemcoCommissionFeeResult{
				Code:    "200",
				Message: "good",
				Result:  []RemcoCommissionFee{*remcoFee},
			},
			serviceErr:    nil,
			expectedURL:   basePath + "/remco-sf",
			expectedResp:  []RemcoCommissionFee{*remcoFee},
			expectedError: nil,
		},
		{
			name:        "perahub error",
			serviceResp: nil,
			serviceErr: &perahub.Error{
				GRPCCode: codes.InvalidArgument,
				Msg:      "something is invalid",
			},
			expectedURL:   basePath + "/remco-sf",
			expectedResp:  nil,
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "something is invalid"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			mockPerahubService := new(MockPerahubService)

			if tc.serviceResp != nil || tc.serviceErr != nil {
				var expectedServiceRes json.RawMessage = nil
				if tc.serviceResp != nil {
					expectedServiceRes, err = json.Marshal(tc.serviceResp)
					require.Nil(t, err)
				}

				mockPerahubService.On("GetRevComm", ctx, tc.expectedURL).
					Return(expectedServiceRes, tc.serviceErr)
			}

			client := NewRevCommClient(mockPerahubService, baseUrl)

			res, err := client.ListRemcoCommissionFee(ctx)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResp, res)
		})
	}

	t.Run("unmarshal error", func(t *testing.T) {
		ctx := context.Background()
		mockPerahubService := new(MockPerahubService)
		mockPerahubService.On("GetRevComm", ctx, basePath+"/remco-sf").
			Return(json.RawMessage("invalid"), nil)

		client := NewRevCommClient(mockPerahubService, baseUrl)

		res, err := client.ListRemcoCommissionFee(ctx)

		assert.Nil(t, res)
		assert.Equal(t, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError), err)
	})
}

func TestClient_GetRemcoCommissionFeeByID(t *testing.T) {
	basePath := "http://some-url.test"
	baseUrl, err := url.Parse(basePath)
	require.Nil(t, err)

	for _, tc := range []struct {
		name          string
		inputFeeID    uint32
		serviceResp   *RemcoCommissionFeeResult
		serviceErr    error
		expectedURL   string
		expectedResp  *RemcoCommissionFee
		expectedError error
	}{
		{
			name:       "happy path",
			inputFeeID: 1,
			serviceResp: &RemcoCommissionFeeResult{
				Code:    "200",
				Message: "Good",
				Result:  getMockRemcoCommissionFee("1"),
			},
			serviceErr:    nil,
			expectedURL:   fmt.Sprintf("%s/remco-sf/%d", basePath, 1),
			expectedResp:  getMockRemcoCommissionFee("1"),
			expectedError: nil,
		},
		{
			name:          "empty remco commission fee id",
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "remco commission fee id required"),
		},
		{
			name:        "perahub error",
			inputFeeID:  1,
			serviceResp: nil,
			serviceErr: &perahub.Error{
				GRPCCode: codes.InvalidArgument,
				Msg:      "something is invalid",
			},
			expectedURL:   fmt.Sprintf("%s/remco-sf/%d", basePath, 1),
			expectedResp:  nil,
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "something is invalid"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			mockPerahubService := new(MockPerahubService)

			if tc.serviceResp != nil || tc.serviceErr != nil {
				var expectedServiceRes json.RawMessage = nil
				if tc.serviceResp != nil {
					expectedServiceRes, err = json.Marshal(tc.serviceResp)
					require.Nil(t, err)
				}

				mockPerahubService.On("GetRevComm", ctx, tc.expectedURL).
					Return(expectedServiceRes, tc.serviceErr)
			}

			client := NewRevCommClient(mockPerahubService, baseUrl)

			res, err := client.GetRemcoCommissionFeeByID(ctx, tc.inputFeeID)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResp, res)
		})
	}
}

func TestClient_CreateRemcoCommissionFee(t *testing.T) {
	basePath := "http://some-url.test"
	baseUrl, err := url.Parse(basePath)
	require.Nil(t, err)

	for _, tc := range []struct {
		name          string
		inputReq      *SaveRemcoCommissionFeeRequest
		serviceResp   *RemcoCommissionFeeResult
		serviceErr    error
		expectedURL   string
		expectedResp  *RemcoCommissionFee
		expectedError error
	}{
		{
			name: "happy path",
			inputReq: &SaveRemcoCommissionFeeRequest{
				RemcoID: "222",
			},
			serviceResp: &RemcoCommissionFeeResult{
				Code:    "200",
				Message: "Good",
				Result:  getMockRemcoCommissionFee("1"),
			},
			serviceErr:    nil,
			expectedURL:   basePath + "/remco-sf",
			expectedResp:  getMockRemcoCommissionFee("1"),
			expectedError: nil,
		},
		{
			name:          "empty request",
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "empty request"),
		},
		{
			name: "perahub error",
			inputReq: &SaveRemcoCommissionFeeRequest{
				RemcoID: "222",
			},
			serviceResp: nil,
			serviceErr: &perahub.Error{
				GRPCCode: codes.InvalidArgument,
				Msg:      "something is invalid",
			},
			expectedURL:   basePath + "/remco-sf",
			expectedResp:  nil,
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "something is invalid"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			mockPerahubService := new(MockPerahubService)

			if tc.inputReq != nil && (tc.serviceResp != nil || tc.serviceErr != nil) {
				var expectedServiceRes json.RawMessage = nil
				if tc.serviceResp != nil {
					expectedServiceRes, err = json.Marshal(tc.serviceResp)
					require.Nil(t, err)
				}

				mockPerahubService.On("PostRevComm", ctx, tc.expectedURL, *tc.inputReq).
					Return(expectedServiceRes, tc.serviceErr)
			}

			client := NewRevCommClient(mockPerahubService, baseUrl)

			res, err := client.CreateRemcoCommissionFee(ctx, tc.inputReq)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResp, res)
		})
	}
}

func TestClient_UpdateRemcoCommissionFee(t *testing.T) {
	basePath := "http://some-url.test"
	baseUrl, err := url.Parse(basePath)
	require.Nil(t, err)

	for _, tc := range []struct {
		name          string
		inputReq      *SaveRemcoCommissionFeeRequest
		serviceResp   *RemcoCommissionFeeResult
		serviceErr    error
		expectedURL   string
		expectedResp  *RemcoCommissionFee
		expectedError error
	}{
		{
			name: "happy path",
			inputReq: &SaveRemcoCommissionFeeRequest{
				FeeID:   1,
				RemcoID: "222",
			},
			serviceResp: &RemcoCommissionFeeResult{
				Code:    "200",
				Message: "Good",
				Result:  getMockRemcoCommissionFee("1"),
			},
			serviceErr:    nil,
			expectedURL:   basePath + "/remco-sf/1",
			expectedResp:  getMockRemcoCommissionFee("1"),
			expectedError: nil,
		},
		{
			name:          "empty request",
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "empty request"),
		},
		{
			name: "remco commission fee id required",
			inputReq: &SaveRemcoCommissionFeeRequest{
				RemcoID: "222",
			},
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "remco commission fee id required"),
		},
		{
			name: "perahub error",
			inputReq: &SaveRemcoCommissionFeeRequest{
				FeeID:   1,
				RemcoID: "222",
			},
			serviceResp: nil,
			serviceErr: &perahub.Error{
				GRPCCode: codes.InvalidArgument,
				Msg:      "something is invalid",
			},
			expectedURL:   basePath + "/remco-sf/1",
			expectedResp:  nil,
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "something is invalid"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			mockPerahubService := new(MockPerahubService)

			if tc.inputReq != nil && (tc.serviceResp != nil || tc.serviceErr != nil) {
				var expectedServiceRes json.RawMessage = nil
				if tc.serviceResp != nil {
					expectedServiceRes, err = json.Marshal(tc.serviceResp)
					require.Nil(t, err)
				}

				mockPerahubService.On("PutRevComm", ctx, tc.expectedURL, *tc.inputReq).
					Return(expectedServiceRes, tc.serviceErr)
			}

			client := NewRevCommClient(mockPerahubService, baseUrl)

			res, err := client.UpdateRemcoCommissionFee(ctx, tc.inputReq)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResp, res)
		})
	}
}

func TestClient_DeleteRemcoCommissionFee(t *testing.T) {
	basePath := "http://some-url.test"
	baseUrl, err := url.Parse(basePath)
	require.Nil(t, err)

	for _, tc := range []struct {
		name          string
		inputFeeID    uint32
		serviceResp   *RemcoCommissionFeeResult
		serviceErr    error
		expectedURL   string
		expectedResp  *RemcoCommissionFee
		expectedError error
	}{
		{
			name:       "happy path",
			inputFeeID: 1,
			serviceResp: &RemcoCommissionFeeResult{
				Code:    "200",
				Message: "Good",
				Result:  getMockRemcoCommissionFee("1"),
			},
			serviceErr:    nil,
			expectedURL:   fmt.Sprintf("%s/remco-sf/%d", basePath, 1),
			expectedResp:  getMockRemcoCommissionFee("1"),
			expectedError: nil,
		},
		{
			name:          "empty remco commission fee id",
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "remco commission fee id required"),
		},
		{
			name:        "perahub error",
			inputFeeID:  1,
			serviceResp: nil,
			serviceErr: &perahub.Error{
				GRPCCode: codes.InvalidArgument,
				Msg:      "something is invalid",
			},
			expectedURL:   fmt.Sprintf("%s/remco-sf/%d", basePath, 1),
			expectedResp:  nil,
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "something is invalid"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			mockPerahubService := new(MockPerahubService)

			if tc.serviceResp != nil || tc.serviceErr != nil {
				var expectedServiceRes json.RawMessage = nil
				if tc.serviceResp != nil {
					expectedServiceRes, err = json.Marshal(tc.serviceResp)
					require.Nil(t, err)
				}

				mockPerahubService.On("DeleteRevComm", ctx, tc.expectedURL).
					Return(expectedServiceRes, tc.serviceErr)
			}

			client := NewRevCommClient(mockPerahubService, baseUrl)

			res, err := client.DeleteRemcoCommissionFee(ctx, tc.inputFeeID)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResp, res)
		})
	}
}

func getMockRemcoCommissionFee(id json.Number) *RemcoCommissionFee {
	return &RemcoCommissionFee{
		FeeID:               id,
		RemcoID:             "22",
		MinAmount:           "1000",
		MaxAmount:           "2000",
		ServiceFee:          "12",
		CommissionAmount:    "0",
		CommissionAmountOTC: "0",
		CommissionType:      CommissionTypeAbsolute,
		TrxType:             TrxTypeInbound,
		UpdatedBy:           "DRP",
	}
}
