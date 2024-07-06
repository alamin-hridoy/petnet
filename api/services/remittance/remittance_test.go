package remittance

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	bpc "brank.as/petnet/api/core/remittance"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/services"
	"brank.as/petnet/api/storage/postgres"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/auth/hydra"
)

var _testStorage *postgres.Storage

func TestMain(m *testing.M) {
	const dbConnEnv = "DATABASE_CONNECTION"
	ddlConnStr := os.Getenv(dbConnEnv)
	if ddlConnStr == "" {
		log.Printf("%s is not set, skipping", dbConnEnv)
		return
	}

	var teardown func()
	_testStorage, teardown = postgres.NewTestStorage(ddlConnStr, filepath.Join("..", "..", "migrations", "sql"))

	exitCode := m.Run()

	if teardown != nil {
		teardown()
	}

	os.Exit(exitCode)
}

func newTestStorage(tb testing.TB) *postgres.Storage {
	if testing.Short() {
		tb.Skip("skipping tests that use postgres on -short")
	}

	return _testStorage
}

func newTestSvc(t *testing.T, st *postgres.Storage) (*Svc, *services.Mock) {
	cl := perahub.NewTestHTTPMock(st, perahub.MockConfig{})
	log := logrus.New().WithField("stage", "testing")
	ph, err := perahub.New(cl,
		"dev",
		"https://newkycgateway.dev.perahub.com.ph/gateway/",
		"https://privatedrp.dev.perahub.com.ph/v1/remit/nonex/",
		"https://privatedrp.dev.perahub.com.ph/v1/billspay/wrapper/api/",
		"https://privatedrp.dev.perahub.com.ph/v1/transactions/api/",
		"https://privatedrp.dev.perahub.com.ph/v1/billspay/",
		"partner-id",
		"client-key",
		"api-key",
		"",
		"",
		nil,
		perahub.WithLogger(log),
		perahub.WithPerahubDefaultAPIKey("api-key"),
		perahub.WithCiCoURL("https://privatedrp.dev.perahub.com.ph/v1/cico/wrapper/"),
		perahub.WithPHRemittanceURL("https://privatedrp.dev.perahub.com.ph/v1/remit/dmt/"),
	)
	if err != nil {
		t.Fatal("setting up perahub integration: ", err)
	}
	m := &services.Mock{}
	return New(bpc.New(st, ph)), m
}

func TestRemittance(t *testing.T) {
	st := newTestStorage(t)
	tests := []struct {
		apiName string
		desc    string
		in      interface{}
		want    interface{}
		tops    cmp.Options
	}{
		{
			apiName: "send-validate",
			desc:    "send validate",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.ValidateSendMoneyResult{}), cmpopts.IgnoreUnexported(bpa.ValidateSendMoneyResponse{}, bpa.ValidateSendMoneyResult{}, timestamppb.Timestamp{})},
			in: &bpa.ValidateSendMoneyRequest{
				PartnerReferenceNumber: "12JEKDWW213",
				PrincipalAmount:        "10000",
				ServiceFee:             "50",
				IsoCurrency:            "PHP",
				ConversionRate:         "1",
				IsoOriginatingCountry:  "PHP",
				IsoDestinationCountry:  "PHP",
				SenderLastName:         "HERMO",
				SenderFirstName:        "IRENE",
				SenderMiddleName:       "M",
				ReceiverLastName:       "HERMO",
				ReceiverFirstName:      "SONNY",
				ReceiverMiddleName:     "D",
				SenderBirthDate:        "1981-06-12",
				SenderBirthPlace:       "TARLAC",
				SenderBirthCountry:     "PH",
				SenderGender:           "FEMALE",
				SenderRelationship:     "SPOUSE",
				SenderPurpose:          "GIFT",
				SenderOccupation:       "DOCTOR",
				SenderEmploymentNature: "IT",
				SendPartnerCode:        "USP",
				SenderSourceOfFund:     "SALARY",
			},
			want: &bpa.ValidateSendMoneyResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.ValidateSendMoneyResult{
					SendValidateReferenceNumber: "1653296685161",
				},
			},
		},
		{
			apiName: "confirm-send-money",
			desc:    "confirm-send-money",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.ConfirmSendMoneyResult{}), cmpopts.IgnoreUnexported(bpa.ConfirmSendMoneyResponse{}, bpa.ConfirmSendMoneyResult{}, timestamppb.Timestamp{})},
			in: &bpa.ConfirmSendMoneyRequest{
				SendValidateReferenceNumber: "084396e59ed3a4acc8da4bd8885bec01",
			},
			want: &bpa.ConfirmSendMoneyResponse{
				Code:    200,
				Message: "Successful",
				Result: &bpa.ConfirmSendMoneyResult{
					Phrn: "PH1654787564",
				},
			},
		},
		{
			apiName: "cancel-send-money",
			desc:    "cancel-send-money",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.CancelSendMoneyResult{}), cmpopts.IgnoreUnexported(bpa.CancelSendMoneyResponse{}, bpa.CancelSendMoneyResult{}, timestamppb.Timestamp{})},
			in: &bpa.CancelSendMoneyRequest{
				Phrn:        "PH1654789142",
				PartnerCode: "DRP",
				Remarks:     "Sample cancel send",
			},
			want: &bpa.CancelSendMoneyResponse{
				Code:    200,
				Message: "Cancel Send Remittance successful",
				Result: &bpa.CancelSendMoneyResult{
					Phrn:                      "PH1654789142",
					CancelSendDate:            "2022-06-14 19:35:02",
					CancelSendReferenceNumber: "6a6e74400561a300d627aba12107bb6c",
				},
			},
		},
		{
			apiName: "validate-receive-money",
			desc:    "validate-receive-money",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.ValidateReceiveMoneyResult{}), cmpopts.IgnoreUnexported(bpa.ValidateReceiveMoneyResponse{}, bpa.ValidateReceiveMoneyResult{}, timestamppb.Timestamp{})},
			in: &bpa.ValidateReceiveMoneyRequest{
				Phrn:                  "PH1654789142",
				PrincipalAmount:       "10000",
				IsoOriginatingCountry: "PHP",
				IsoDestinationCountry: "PHP",
				SenderLastName:        "HERMO",
				SenderFirstName:       "IRENE",
				SenderMiddleName:      "M",
				ReceiverLastName:      "HERMO",
				ReceiverFirstName:     "SONNY",
				ReceiverMiddleName:    "D",
				PayoutPartnerCode:     "USP",
			},
			want: &bpa.ValidateReceiveMoneyResponse{
				Code:    200,
				Message: "Successful",
				Result: &bpa.ValidateReceiveMoneyResult{
					PayoutValidateReferenceNumber: "4f8a09d3b293807aa50305f66d6cc73c",
				},
			},
		},
		{
			apiName: "inquire",
			desc:    "inquire",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.InquireResult{}), cmpopts.IgnoreUnexported(bpa.InquireResponse{}, bpa.InquireResult{}, timestamppb.Timestamp{})},
			in: &bpa.InquireRequest{
				Phrn: "PH1654789142",
			},
			want: &bpa.InquireResponse{
				Code:    200,
				Message: "PeraHUB Reference Number (PHRN) is available for Payout",
				Result: &bpa.InquireResult{
					Phrn:                  "PH1658296732",
					PrincipalAmount:       10000,
					IsoCurrency:           "PHP",
					ConversionRate:        1,
					IsoOriginatingCountry: "PHP",
					IsoDestinationCountry: "PHP",
					SenderLastName:        "HERMO",
					SenderFirstName:       "IRENE",
					SenderMiddleName:      "M",
					ReceiverLastName:      "HERMO",
					ReceiverFirstName:     "SONNY",
					ReceiverMiddleName:    "D",
				},
			},
		},
		{
			apiName: "Confirm-receive-money",
			desc:    "Confirm-receive-money",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.ConfirmReceiveMoneyResult{}), cmpopts.IgnoreUnexported(bpa.ConfirmReceiveMoneyResponse{}, bpa.ConfirmReceiveMoneyResult{}, timestamppb.Timestamp{})},
			in: &bpa.ConfirmReceiveMoneyRequest{
				PayoutValidateReferenceNumber: "4f8a09d3b293807aa50305f66d6cc73c",
			},
			want: &bpa.ConfirmReceiveMoneyResponse{
				Code:    200,
				Message: "Successful",
				Result: &bpa.ConfirmReceiveMoneyResult{
					Phrn:                  "PH1654789142",
					PrincipalAmount:       10000,
					IsoOriginatingCountry: "PHP",
					IsoDestinationCountry: "PHP",
					SenderLastName:        "HERMO",
					SenderFirstName:       "IRENE",
					SenderMiddleName:      "M",
					ReceiverLastName:      "HERMO",
					ReceiverFirstName:     "SONNY",
					ReceiverMiddleName:    "D",
				},
			},
		},
		{
			apiName: "partner-grid",
			desc:    "partner grid",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.PartnersGridResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.PartnersGridResponse{}, bpa.PartnersGridResult{}, timestamppb.Timestamp{})},
			want: &bpa.PartnersGridResponse{
				Code:    200,
				Message: "Good",
				Result: []*bpa.PartnersGridResult{
					{
						ID:           1,
						PartnerCode:  "DRP",
						PartnerName:  "PERA HUB",
						ClientSecret: "26da230221d9e506b1fd823df1869875",
						Status:       1,
						CreatedAt:    timestamppb.Now(),
						UpdatedAt:    timestamppb.Now(),
						DeletedAt:    timestamppb.Now(),
					},
					{
						ID:           2,
						PartnerCode:  "USP",
						PartnerName:  "PERA HUB",
						ClientSecret: "12358fbef0bb08d7a7bab57df956a335",
						Status:       1,
						CreatedAt:    timestamppb.Now(),
						UpdatedAt:    timestamppb.Now(),
						DeletedAt:    timestamppb.Now(),
					},
				},
			},
		},
		{
			apiName: "purpose-of-remittance-grid",
			desc:    "purpose-of-remittance-grid",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.PurposeOfRemittanceGridResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.PurposeOfRemittanceGridResponse{}, bpa.PurposeOfRemittanceGridResult{}, timestamppb.Timestamp{})},
			want: &bpa.PurposeOfRemittanceGridResponse{
				Code:    200,
				Message: "Good",
				Result: []*bpa.PurposeOfRemittanceGridResult{
					{
						ID:                  "1",
						PurposeOfRemittance: "Gift",
						CreatedAt:           timestamppb.Now(),
						UpdatedAt:           timestamppb.Now(),
						DeletedAt:           timestamppb.Now(),
					},
					{
						ID:                  "2",
						PurposeOfRemittance: "Fund",
						CreatedAt:           timestamppb.Now(),
						UpdatedAt:           timestamppb.Now(),
						DeletedAt:           timestamppb.Now(),
					},
					{
						ID:                  "3",
						PurposeOfRemittance: "Allowance",
						CreatedAt:           timestamppb.Now(),
						UpdatedAt:           timestamppb.Now(),
						DeletedAt:           timestamppb.Now(),
					},
				},
			},
		},
		{
			apiName: "purpose-of-remittance-get",
			desc:    "purpose of remittance get",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.PurposeOfRemittanceGetResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.PurposeOfRemittanceGetResponse{}, bpa.PurposeOfRemittanceGetResult{}, timestamppb.Timestamp{})},
			in: &bpa.PurposeOfRemittanceGetRequest{
				ID: "1",
			},
			want: &bpa.PurposeOfRemittanceGetResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.PurposeOfRemittanceGetResult{
					ID:                  1,
					PurposeOfRemittance: "Donation",
					CreatedAt:           timestamppb.Now(),
					UpdatedAt:           timestamppb.Now(),
					DeletedAt:           timestamppb.Now(),
				},
			},
		},
		{
			apiName: "purpose-of-remittance-update",
			desc:    "purpose of remittance update",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.PurposeOfRemittanceUpdateResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.PurposeOfRemittanceUpdateResponse{}, bpa.PurposeOfRemittanceUpdateResult{}, timestamppb.Timestamp{})},
			in: &bpa.PurposeOfRemittanceUpdateRequest{
				PurposeOfRemittance: "Donation",
				ID:                  "1",
			},
			want: &bpa.PurposeOfRemittanceUpdateResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.PurposeOfRemittanceUpdateResult{
					ID:                  1,
					PurposeOfRemittance: "Donation",
					CreatedAt:           timestamppb.Now(),
					UpdatedAt:           timestamppb.Now(),
					DeletedAt:           timestamppb.Now(),
				},
			},
		},
		{
			apiName: "purpose-of-remittance-create",
			desc:    "purpose-of-remittance-create",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.PurposeOfRemittanceCreateResult{}, "CreatedAt", "UpdatedAt"), cmpopts.IgnoreUnexported(bpa.PurposeOfRemittanceCreateResponse{}, bpa.PurposeOfRemittanceCreateResult{}, timestamppb.Timestamp{})},
			in: &bpa.PurposeOfRemittanceCreateRequest{
				PurposeOfRemittance: "Donation",
			},
			want: &bpa.PurposeOfRemittanceCreateResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.PurposeOfRemittanceCreateResult{
					ID:                  1,
					PurposeOfRemittance: "USP",
					CreatedAt:           &timestamppb.Timestamp{},
					UpdatedAt:           &timestamppb.Timestamp{},
				},
			},
		},
		{
			apiName: "partner-create",
			desc:    "partner create",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.PartnersCreateResult{}), cmpopts.IgnoreUnexported(bpa.PartnersCreateResponse{}, bpa.PartnersCreateResult{}, timestamppb.Timestamp{})},
			in: &bpa.PartnersCreateRequest{
				PartnerCode: "DRP",
				PartnerName: "BRANKAS",
			},
			want: &bpa.PartnersCreateResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.PartnersCreateResult{
					ID:           1,
					PartnerCode:  "USP",
					PartnerName:  "PERA HUB",
					ClientSecret: "adawdawdawd",
					CreatedAt:    &timestamppb.Timestamp{},
					UpdatedAt:    &timestamppb.Timestamp{},
				},
			},
		},
		{
			apiName: "source-of-fund-grid",
			desc:    "source of fund grid",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.SourceOfFundGridResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.SourceOfFundGridResponse{}, bpa.SourceOfFundGridResult{}, timestamppb.Timestamp{})},
			want: &bpa.SourceOfFundGridResponse{
				Code:    200,
				Message: "Good",
				Result: []*bpa.SourceOfFundGridResult{
					{
						ID:           "1",
						SourceOfFund: "SALARY",
						CreatedAt:    timestamppb.Now(),
						UpdatedAt:    timestamppb.Now(),
						DeletedAt:    timestamppb.Now(),
					},
					{
						ID:           "2",
						SourceOfFund: "BUSINESS",
						CreatedAt:    timestamppb.Now(),
						UpdatedAt:    timestamppb.Now(),
						DeletedAt:    timestamppb.Now(),
					},
					{
						ID:           "3",
						SourceOfFund: "REMITTANCE",
						CreatedAt:    timestamppb.Now(),
						UpdatedAt:    timestamppb.Now(),
						DeletedAt:    timestamppb.Now(),
					},
				},
			},
		},
		{
			apiName: "source-of-fund-create",
			desc:    "source of fund create",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.SourceOfFundCreateResult{}, "CreatedAt", "UpdatedAt"), cmpopts.IgnoreUnexported(bpa.SourceOfFundCreateResponse{}, bpa.SourceOfFundCreateResult{}, timestamppb.Timestamp{})},
			in: &bpa.SourceOfFundCreateRequest{
				SourceOfFund: "SALARY",
			},
			want: &bpa.SourceOfFundCreateResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.SourceOfFundCreateResult{
					ID:           1,
					SourceOfFund: "SALARY",
					CreatedAt:    &timestamppb.Timestamp{},
					UpdatedAt:    &timestamppb.Timestamp{},
				},
			},
		},
		{
			apiName: "source-of-fund-get",
			desc:    "source of fund get",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.SourceOfFundGetResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.SourceOfFundGetResponse{}, bpa.SourceOfFundGetResult{}, timestamppb.Timestamp{})},
			in: &bpa.SourceOfFundGetRequest{
				ID: "1",
			},
			want: &bpa.SourceOfFundGetResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.SourceOfFundGetResult{
					ID:           1,
					SourceOfFund: "SALARY",
					CreatedAt:    timestamppb.Now(),
					UpdatedAt:    timestamppb.Now(),
					DeletedAt:    timestamppb.Now(),
				},
			},
		},
		{
			apiName: "employment-grid",
			desc:    "employment grid",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.EmploymentGridResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.EmploymentGridResponse{}, bpa.EmploymentGridResult{}, timestamppb.Timestamp{})},
			want: &bpa.EmploymentGridResponse{
				Code:    200,
				Message: "Good",
				Result: []*bpa.EmploymentGridResult{
					{
						ID:               1,
						EmploymentNature: "REGULAR",
						CreatedAt:        timestamppb.Now(),
						UpdatedAt:        timestamppb.Now(),
						DeletedAt:        timestamppb.Now(),
					},
					{
						ID:               2,
						EmploymentNature: "PROBATIONARY",
						CreatedAt:        timestamppb.Now(),
						UpdatedAt:        timestamppb.Now(),
						DeletedAt:        timestamppb.Now(),
					},
					{
						ID:               3,
						EmploymentNature: "CONTRACTUAL",
						CreatedAt:        timestamppb.Now(),
						UpdatedAt:        timestamppb.Now(),
						DeletedAt:        timestamppb.Now(),
					},
				},
			},
		},
		{
			apiName: "remittance-employment-create",
			desc:    "remittance employment create",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.RemittanceEmploymentCreateResult{}, "CreatedAt", "UpdatedAt"), cmpopts.IgnoreUnexported(bpa.RemittanceEmploymentCreateResponse{}, bpa.RemittanceEmploymentCreateResult{}, timestamppb.Timestamp{})},
			in: &bpa.RemittanceEmploymentCreateRequest{
				Employment:       "REGULAR",
				EmploymentNature: "123",
			},
			want: &bpa.RemittanceEmploymentCreateResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.RemittanceEmploymentCreateResult{
					ID:               1,
					EmploymentNature: "REGULAR",
					CreatedAt:        &timestamppb.Timestamp{},
					UpdatedAt:        &timestamppb.Timestamp{},
				},
			},
		},
		{
			apiName: "employment-update",
			desc:    "remittance employment update",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.RemittanceEmploymentUpdateResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.RemittanceEmploymentUpdateResponse{}, bpa.RemittanceEmploymentUpdateResult{}, timestamppb.Timestamp{})},
			in: &bpa.RemittanceEmploymentUpdateRequest{
				ID:               "1",
				Employment:       "REGULAR",
				EmploymentNature: "REGULAR",
			},
			want: &bpa.RemittanceEmploymentUpdateResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.RemittanceEmploymentUpdateResult{
					ID:               1,
					EmploymentNature: "REGULAR",
					CreatedAt:        timestamppb.Now(),
					UpdatedAt:        timestamppb.Now(),
					DeletedAt:        timestamppb.Now(),
				},
			},
		},
		{
			apiName: "occupation-grid",
			desc:    "occupation grid",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.OccupationGridResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.OccupationGridResponse{}, bpa.OccupationGridResult{}, timestamppb.Timestamp{})},
			want: &bpa.OccupationGridResponse{
				Code:    200,
				Message: "Good",
				Result: []*bpa.OccupationGridResult{
					{
						ID:         1,
						Occupation: "Programmer",
						CreatedAt:  timestamppb.Now(),
						UpdatedAt:  timestamppb.Now(),
						DeletedAt:  timestamppb.Now(),
					},
					{
						ID:         2,
						Occupation: "Engineer",
						CreatedAt:  timestamppb.Now(),
						UpdatedAt:  timestamppb.Now(),
						DeletedAt:  timestamppb.Now(),
					},
					{
						ID:         3,
						Occupation: "Doctor",
						CreatedAt:  timestamppb.Now(),
						UpdatedAt:  timestamppb.Now(),
						DeletedAt:  timestamppb.Now(),
					},
				},
			},
		},
		{
			apiName: "occupation-get",
			desc:    "occupation get",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.OccupationGetResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.OccupationGetResponse{}, bpa.OccupationGetResult{}, timestamppb.Timestamp{})},
			in: &bpa.OccupationGetRequest{
				ID: "1",
			},
			want: &bpa.OccupationGetResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.OccupationGetResult{
					ID:         1,
					Occupation: "Programmer",
					CreatedAt:  timestamppb.Now(),
					UpdatedAt:  timestamppb.Now(),
					DeletedAt:  timestamppb.Now(),
				},
			},
		},
		{
			apiName: "occupation-create",
			desc:    "occupation create",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.OccupationCreateResult{}, "CreatedAt", "UpdatedAt"), cmpopts.IgnoreUnexported(bpa.OccupationCreateResponse{}, bpa.OccupationCreateResult{}, timestamppb.Timestamp{})},
			in: &bpa.OccupationCreateRequest{
				Occupation: "Programmer",
			},
			want: &bpa.OccupationCreateResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.OccupationCreateResult{
					ID:         1,
					Occupation: "Programmer",
					CreatedAt:  &timestamppb.Timestamp{},
					UpdatedAt:  &timestamppb.Timestamp{},
				},
			},
		},
		{
			apiName: "occupation-update",
			desc:    "occupation update",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.OccupationUpdateResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.OccupationUpdateResponse{}, bpa.OccupationUpdateResult{}, timestamppb.Timestamp{})},
			in: &bpa.OccupationUpdateRequest{
				Occupation: "Programmer",
				ID:         "1",
			},
			want: &bpa.OccupationUpdateResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.OccupationUpdateResult{
					ID:         1,
					Occupation: "Programmer",
					CreatedAt:  timestamppb.Now(),
					UpdatedAt:  timestamppb.Now(),
					DeletedAt:  timestamppb.Now(),
				},
			},
		},
		{
			apiName: "occupation-delete",
			desc:    "occupation delete",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.OccupationDeleteResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.OccupationDeleteResponse{}, bpa.OccupationDeleteResult{}, timestamppb.Timestamp{})},
			in: &bpa.OccupationDeleteRequest{
				ID: "1",
			},
			want: &bpa.OccupationDeleteResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.OccupationDeleteResult{
					ID:         1,
					Occupation: "Programmer",
					CreatedAt:  timestamppb.Now(),
					UpdatedAt:  timestamppb.Now(),
					DeletedAt:  timestamppb.Now(),
				},
			},
		},
		{
			apiName: "relationship-get",
			desc:    "relationship get",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.RelationshipGetResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.RelationshipGetResponse{}, bpa.RelationshipGetResult{}, timestamppb.Timestamp{})},
			in: &bpa.RelationshipGetRequest{
				ID: "1",
			},
			want: &bpa.RelationshipGetResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.RelationshipGetResult{
					ID:           1,
					Relationship: "Friend",
					CreatedAt:    timestamppb.Now(),
					UpdatedAt:    timestamppb.Now(),
					DeletedAt:    timestamppb.Now(),
				},
			},
		},
		{
			apiName: "source-of-fund-update",
			desc:    "source of fund update",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.SourceOfFundUpdateResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.SourceOfFundUpdateResponse{}, bpa.SourceOfFundUpdateResult{}, timestamppb.Timestamp{})},
			in: &bpa.SourceOfFundUpdateRequest{
				SourceOfFund: "SALARY",
				ID:           "1",
			},
			want: &bpa.SourceOfFundUpdateResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.SourceOfFundUpdateResult{
					ID:           1,
					SourceOfFund: "SALARY",
					CreatedAt:    timestamppb.Now(),
					UpdatedAt:    timestamppb.Now(),
					DeletedAt:    timestamppb.Now(),
				},
			},
		},
		{
			apiName: "source-of-fund-delete",
			desc:    "source of fund delete",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.SourceOfFundDeleteResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.SourceOfFundDeleteResponse{}, bpa.SourceOfFundDeleteResult{}, timestamppb.Timestamp{})},
			in: &bpa.SourceOfFundDeleteRequest{
				ID: "1",
			},
			want: &bpa.SourceOfFundDeleteResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.SourceOfFundDeleteResult{
					ID:           1,
					SourceOfFund: "SALARY",
					CreatedAt:    timestamppb.Now(),
					UpdatedAt:    timestamppb.Now(),
					DeletedAt:    timestamppb.Now(),
				},
			},
		},
		{
			apiName: "purpose-of-remittance-delete",
			desc:    "purpose of remittance delete",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.PurposeOfRemittanceDeleteResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.PurposeOfRemittanceDeleteResponse{}, bpa.PurposeOfRemittanceDeleteResult{}, timestamppb.Timestamp{})},
			in: &bpa.PurposeOfRemittanceDeleteRequest{
				ID: "1",
			},
			want: &bpa.PurposeOfRemittanceDeleteResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.PurposeOfRemittanceDeleteResult{
					ID:                  1,
					PurposeOfRemittance: "Gift",
					CreatedAt:           timestamppb.Now(),
					UpdatedAt:           timestamppb.Now(),
					DeletedAt:           timestamppb.Now(),
				},
			},
		},
		{
			apiName: "relationship-delete",
			desc:    "relationship delete",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.RelationshipDeleteResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.RelationshipDeleteResponse{}, bpa.RelationshipDeleteResult{}, timestamppb.Timestamp{})},
			in: &bpa.RelationshipDeleteRequest{
				ID: "1",
			},
			want: &bpa.RelationshipDeleteResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.RelationshipDeleteResult{
					ID:           1,
					Relationship: "Friend",
					CreatedAt:    timestamppb.Now(),
					UpdatedAt:    timestamppb.Now(),
					DeletedAt:    timestamppb.Now(),
				},
			},
		},
		{
			apiName: "employment-get",
			desc:    "employment get",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.EmploymentGetResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.EmploymentGetResponse{}, bpa.EmploymentGetResult{}, timestamppb.Timestamp{})},
			in:      &bpa.EmploymentGetRequest{ID: "1"},
			want: &bpa.EmploymentGetResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.EmploymentGetResult{
					ID: 1, EmploymentNature: "REGULAR",
					CreatedAt: timestamppb.Now(),
					UpdatedAt: timestamppb.Now(),
					DeletedAt: timestamppb.Now(),
				},
			},
		},
		{
			apiName: "employment-delete",
			desc:    "remittance employment delete",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.RemittanceEmploymentDeleteResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.RemittanceEmploymentDeleteResponse{}, bpa.RemittanceEmploymentDeleteResult{}, timestamppb.Timestamp{})},
			in: &bpa.RemittanceEmploymentDeleteRequest{
				ID: "1",
			},
			want: &bpa.RemittanceEmploymentDeleteResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.RemittanceEmploymentDeleteResult{
					ID:               1,
					EmploymentNature: "REGULAR",
					CreatedAt:        timestamppb.Now(),
					UpdatedAt:        timestamppb.Now(),
					DeletedAt:        timestamppb.Now(),
				},
			},
		},
		{
			apiName: "relationship-grid",
			desc:    "relationship grid",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.RelationshipGridResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.RelationshipGridResponse{}, bpa.RelationshipGridResult{}, timestamppb.Timestamp{})},
			want: &bpa.RelationshipGridResponse{
				Code:    200,
				Message: "Good",
				Result: []*bpa.RelationshipGridResult{
					{
						ID:           1,
						Relationship: "Friend",
						CreatedAt:    timestamppb.Now(),
						UpdatedAt:    timestamppb.Now(),
						DeletedAt:    timestamppb.Now(),
					},
					{
						ID:           2,
						Relationship: "Father",
						CreatedAt:    timestamppb.Now(),
						UpdatedAt:    timestamppb.Now(),
						DeletedAt:    timestamppb.Now(),
					},
					{
						ID:           3,
						Relationship: "Mother",
						CreatedAt:    timestamppb.Now(),
						UpdatedAt:    timestamppb.Now(),
						DeletedAt:    timestamppb.Now(),
					},
				},
			},
		},
		{
			apiName: "relationship-update",
			desc:    "relationship update",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.RelationshipUpdateResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.RelationshipUpdateResponse{}, bpa.RelationshipUpdateResult{}, timestamppb.Timestamp{})},
			in: &bpa.RelationshipUpdateRequest{
				Relationship: "Friend",
				ID:           "1",
			},
			want: &bpa.RelationshipUpdateResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.RelationshipUpdateResult{
					ID:           1,
					Relationship: "Friend",
					CreatedAt:    timestamppb.Now(),
					UpdatedAt:    timestamppb.Now(),
					DeletedAt:    timestamppb.Now(),
				},
			},
		},
		{
			apiName: "partners-delete",
			desc:    "partners delete",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.PartnersDeleteResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.PartnersDeleteResponse{}, bpa.PartnersDeleteResult{}, timestamppb.Timestamp{})},
			in: &bpa.PartnersDeleteRequest{
				ID:          "1",
				PartnerCode: "USP",
				PartnerName: "PERA HUB 2",
			},
			want: &bpa.PartnersDeleteResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.PartnersDeleteResult{
					ID:           1,
					PartnerCode:  "USP",
					PartnerName:  "PERA HUB",
					ClientSecret: "adawdawdawd",
					CreatedAt:    timestamppb.Now(),
					UpdatedAt:    timestamppb.Now(),
					DeletedAt:    timestamppb.Now(),
				},
			},
		},
		{
			apiName: "partner-get",
			desc:    "partner get",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.PartnersGetResult{}, "CreatedAt", "UpdatedAt"), cmpopts.IgnoreUnexported(bpa.PartnersGetResponse{}, bpa.PartnersGetResult{}, timestamppb.Timestamp{})},
			in: &bpa.PartnersGetRequest{
				ID: "1",
			},
			want: &bpa.PartnersGetResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.PartnersGetResult{
					ID:           1,
					PartnerCode:  "DRP",
					PartnerName:  "BRANKAS",
					ClientSecret: "4fab1de660a6b7faef0168ca4788408a",
					Status:       1,
					CreatedAt:    timestamppb.Now(),
					UpdatedAt:    timestamppb.Now(),
					DeletedAt:    "null",
				},
			},
		},
		{
			apiName: "partner-update",
			desc:    "partner update",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.PartnersUpdateResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"), cmpopts.IgnoreUnexported(bpa.PartnersUpdateResponse{}, bpa.PartnersUpdateResult{}, timestamppb.Timestamp{})},
			in: &bpa.PartnersUpdateRequest{
				ID:          "1",
				PartnerCode: "USP",
				PartnerName: "PERA HUB",
				Service:     "PERA HUB",
			},
			want: &bpa.PartnersUpdateResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.PartnersUpdateResult{
					ID:           1,
					PartnerCode:  "USP",
					PartnerName:  "PERA HUB",
					ClientSecret: "adawdawdawd",
					Status:       1,
					CreatedAt:    timestamppb.Now(),
					UpdatedAt:    timestamppb.Now(),
					DeletedAt:    "null",
				},
			},
		},
		{
			apiName: "relationship-create",
			desc:    "relationship create",
			tops:    cmp.Options{cmpopts.IgnoreFields(bpa.RelationshipCreateResult{}, "CreatedAt", "UpdatedAt"), cmpopts.IgnoreUnexported(bpa.RelationshipCreateResponse{}, bpa.RelationshipCreateResult{}, timestamppb.Timestamp{})},
			in: &bpa.RelationshipCreateRequest{
				Relationship: "Friend",
			},
			want: &bpa.RelationshipCreateResponse{
				Code:    200,
				Message: "Good",
				Result: &bpa.RelationshipCreateResult{
					Relationship: "Friend",
					CreatedAt:    timestamppb.Now(),
					UpdatedAt:    timestamppb.Now(),
					ID:           1,
				},
			},
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())
	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h, _ := newTestSvc(t, st)
			switch test.apiName {
			case "send-validate":
				req, ok := test.in.(*bpa.ValidateSendMoneyRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.ValidateSendMoney(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "confirm-send-money":
				req, ok := test.in.(*bpa.ConfirmSendMoneyRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.ConfirmSendMoney(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "cancel-send-money":
				req, ok := test.in.(*bpa.CancelSendMoneyRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.CancelSendMoney(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "validate-receive-money":
				req, ok := test.in.(*bpa.ValidateReceiveMoneyRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.ValidateReceiveMoney(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "inquire":
				req, ok := test.in.(*bpa.InquireRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.Inquire(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "confirm-receive-money":
				req, ok := test.in.(*bpa.ConfirmReceiveMoneyRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.ConfirmReceiveMoney(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "partner-grid":
				got, err := h.PartnersGrid(ctx, &emptypb.Empty{})
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "purpose-of-remittance-grid":
				got, err := h.PurposeOfRemittanceGrid(ctx, &emptypb.Empty{})
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "purpose-of-remittance-get":
				req, ok := test.in.(*bpa.PurposeOfRemittanceGetRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.PurposeOfRemittanceGet(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "purpose-of-remittance-update":
				req, ok := test.in.(*bpa.PurposeOfRemittanceUpdateRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.PurposeOfRemittanceUpdate(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "purpose-of-remittance-create":
				req, ok := test.in.(*bpa.PurposeOfRemittanceCreateRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.PurposeOfRemittanceCreate(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "source-of-fund-grid":
				got, err := h.SourceOfFundGrid(ctx, &emptypb.Empty{})
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "source-of-fund-create":
				req, ok := test.in.(*bpa.SourceOfFundCreateRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.SourceOfFundCreate(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "source-of-fund-get":
				req, ok := test.in.(*bpa.SourceOfFundGetRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.SourceOfFundGet(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "employment-grid":
				got, err := h.EmploymentGrid(ctx, &emptypb.Empty{})
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "remittance-employment-create":
				req, ok := test.in.(*bpa.RemittanceEmploymentCreateRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.RemittanceEmploymentCreate(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "employment-update":
				req, ok := test.in.(*bpa.RemittanceEmploymentUpdateRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.RemittanceEmploymentUpdate(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "occupation-grid":
				got, err := h.OccupationGrid(ctx, &emptypb.Empty{})
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "occupation-get":
				req, ok := test.in.(*bpa.OccupationGetRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.OccupationGet(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "occupation-update":
				req, ok := test.in.(*bpa.OccupationUpdateRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.OccupationUpdate(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "occupation-create":
				req, ok := test.in.(*bpa.OccupationCreateRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.OccupationCreate(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "occupation-delete":
				req, ok := test.in.(*bpa.OccupationDeleteRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.OccupationDelete(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "relationship-get":
				req, ok := test.in.(*bpa.RelationshipGetRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.RelationshipGet(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "source-of-fund-update":
				req, ok := test.in.(*bpa.SourceOfFundUpdateRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.SourceOfFundUpdate(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "source-of-fund-delete":
				req, ok := test.in.(*bpa.SourceOfFundDeleteRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.SourceOfFundDelete(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "purpose-of-remittance-delete":
				req, ok := test.in.(*bpa.PurposeOfRemittanceDeleteRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.PurposeOfRemittanceDelete(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "relationship-delete":
				req, ok := test.in.(*bpa.RelationshipDeleteRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.RelationshipDelete(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "employment-get":
				req, ok := test.in.(*bpa.EmploymentGetRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.EmploymentGet(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "employment-delete":
				req, ok := test.in.(*bpa.RemittanceEmploymentDeleteRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.RemittanceEmploymentDelete(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "relationship-grid":
				got, err := h.RelationshipGrid(ctx, &emptypb.Empty{})
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "relationship-update":
				req, ok := test.in.(*bpa.RelationshipUpdateRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.RelationshipUpdate(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "partners-delete":
				req, ok := test.in.(*bpa.PartnersDeleteRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.PartnersDelete(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "partner-get":
				req, ok := test.in.(*bpa.PartnersGetRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.PartnersGet(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "partner-update":
				req, ok := test.in.(*bpa.PartnersUpdateRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.PartnersUpdate(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			case "relationship-create":
				req, ok := test.in.(*bpa.RelationshipCreateRequest)
				if !ok {
					t.Error("reqest type conversion error")
				}
				got, err := h.RelationshipCreate(ctx, req)
				if err != nil {
					t.Fatal(err)
				}
				o := test.tops
				if !cmp.Equal(test.want, got, o) {
					t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
				}
			}
		})
	}
}
