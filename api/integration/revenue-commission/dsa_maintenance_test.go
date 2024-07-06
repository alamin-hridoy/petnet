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

func TestClient_GetDSAByID(t *testing.T) {
	basePath := "http://some-url.test"
	baseUrl, err := url.Parse(basePath)
	require.Nil(t, err)

	for _, tc := range []struct {
		name          string
		inputDsaID    uint32
		serviceResp   *DSAResult
		serviceErr    error
		expectedURL   string
		expectedResp  *DSA
		expectedError error
	}{
		{
			name:       "happy path",
			inputDsaID: 1,
			serviceResp: &DSAResult{
				Code:    "200",
				Message: "Good",
				Result:  getMockDSA("1"),
			},
			serviceErr:    nil,
			expectedURL:   fmt.Sprintf("%s/dsa/%d", basePath, 1),
			expectedResp:  getMockDSA("1"),
			expectedError: nil,
		},
		{
			name:          "empty dsa id",
			inputDsaID:    0,
			serviceResp:   nil,
			serviceErr:    nil,
			expectedURL:   fmt.Sprintf("%s/dsa/%d", basePath, 1),
			expectedResp:  nil,
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "dsa id required"),
		},
		{
			name:        "perahub error",
			inputDsaID:  1,
			serviceResp: nil,
			serviceErr: &perahub.Error{
				GRPCCode: codes.InvalidArgument,
				Msg:      "something is invalid",
			},
			expectedURL:   fmt.Sprintf("%s/dsa/%d", basePath, 1),
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

			res, err := client.GetDSAByID(ctx, tc.inputDsaID)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResp, res)
		})
	}

	t.Run("unmarshal error", func(t *testing.T) {
		ctx := context.Background()
		mockPerahubService := new(MockPerahubService)
		mockPerahubService.On("GetRevComm", ctx, basePath+"/dsa/1").
			Return(json.RawMessage("invalid"), nil)

		client := NewRevCommClient(mockPerahubService, baseUrl)

		res, err := client.GetDSAByID(ctx, 1)

		assert.Nil(t, res)
		assert.Equal(t, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError), err)
	})
}

func TestClient_ListDSA(t *testing.T) {
	basePath := "http://some-url.test"
	baseUrl, err := url.Parse(basePath)
	require.Nil(t, err)

	dsa := getMockDSA("1")

	for _, tc := range []struct {
		name          string
		serviceResp   *ListDSAResult
		serviceErr    error
		expectedURL   string
		expectedResp  []DSA
		expectedError error
	}{
		{
			name: "happy path",
			serviceResp: &ListDSAResult{
				Code:    "200",
				Message: "good",
				Result:  []DSA{*dsa},
			},
			serviceErr:    nil,
			expectedURL:   basePath + "/dsa",
			expectedResp:  []DSA{*dsa},
			expectedError: nil,
		},
		{
			name:        "perahub error",
			serviceResp: nil,
			serviceErr: &perahub.Error{
				GRPCCode: codes.InvalidArgument,
				Msg:      "something is invalid",
			},
			expectedURL:   basePath + "/dsa",
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

			res, err := client.ListDSA(ctx)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResp, res)
		})
	}

	t.Run("unmarshal error", func(t *testing.T) {
		ctx := context.Background()
		mockPerahubService := new(MockPerahubService)
		mockPerahubService.On("GetRevComm", ctx, basePath+"/dsa").
			Return(json.RawMessage("invalid"), nil)

		client := NewRevCommClient(mockPerahubService, baseUrl)

		res, err := client.ListDSA(ctx)

		assert.Nil(t, res)
		assert.Equal(t, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError), err)
	})
}

func TestClient_CreateDSA(t *testing.T) {
	basePath := "http://some-url.test"
	baseUrl, err := url.Parse(basePath)
	require.Nil(t, err)

	for _, tc := range []struct {
		name          string
		inputReq      *SaveDSARequest
		serviceResp   *DSAResult
		serviceErr    error
		expectedURL   string
		expectedResp  *DSA
		expectedError error
	}{
		{
			name: "happy path",
			inputReq: &SaveDSARequest{
				DsaCode: "DAS",
			},
			serviceResp: &DSAResult{
				Code:    "200",
				Message: "Good",
				Result:  getMockDSA("1"),
			},
			serviceErr:    nil,
			expectedURL:   basePath + "/dsa",
			expectedResp:  getMockDSA("1"),
			expectedError: nil,
		},
		{
			name:          "empty request",
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "empty request"),
		},
		{
			name: "perahub error",
			inputReq: &SaveDSARequest{
				DsaCode: "DAS",
			},
			serviceResp: nil,
			serviceErr: &perahub.Error{
				GRPCCode: codes.InvalidArgument,
				Msg:      "something is invalid",
			},
			expectedURL:   basePath + "/dsa",
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

			res, err := client.CreateDSA(ctx, tc.inputReq)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResp, res)
		})
	}
}

func TestClient_UpdateDSA(t *testing.T) {
	basePath := "http://some-url.test"
	baseUrl, err := url.Parse(basePath)
	require.Nil(t, err)

	for _, tc := range []struct {
		name          string
		inputReq      *SaveDSARequest
		serviceResp   *DSAResult
		serviceErr    error
		expectedURL   string
		expectedResp  *DSA
		expectedError error
	}{
		{
			name: "happy path",
			inputReq: &SaveDSARequest{
				DsaID:   1,
				DsaCode: "DAS",
			},
			serviceResp: &DSAResult{
				Code:    "200",
				Message: "Good",
				Result:  getMockDSA("1"),
			},
			serviceErr:    nil,
			expectedURL:   basePath + "/dsa/1",
			expectedResp:  getMockDSA("1"),
			expectedError: nil,
		},
		{
			name:          "empty request",
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "empty request"),
		},
		{
			name: "empty dsa id",
			inputReq: &SaveDSARequest{
				DsaCode: "aaa",
			},
			expectedError: coreerror.NewCoreError(codes.InvalidArgument, "dsa id required"),
		},
		{
			name: "perahub error",
			inputReq: &SaveDSARequest{
				DsaID:   1,
				DsaCode: "DAS",
			},
			serviceResp: nil,
			serviceErr: &perahub.Error{
				GRPCCode: codes.InvalidArgument,
				Msg:      "something is invalid",
			},
			expectedURL:   basePath + "/dsa/1",
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

			res, err := client.UpdateDSA(ctx, tc.inputReq)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResp, res)
		})
	}
}

func getMockDSA(dsaID json.Number) *DSA {
	return &DSA{
		DsaID:        dsaID,
		DsaCode:      "UB",
		DsaName:      "UNIONBANK",
		EmailAddress: "ub@petnet.com",
		Status:       "1",
		Vatable:      "2",
	}
}
