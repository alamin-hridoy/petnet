package util

import (
	"context"
	"testing"
	"time"

	"brank.as/petnet/api/storage"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRecordRTA(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	orgID := uuid.NewString()
	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitToAccountHistory{}, "ID", "OrgID", "TrxDate", "Info", "Details", "CreatedBy", "UpdatedBy", "SortByColumn", "SortOrder", "Limit", "Offset", "Total", "Created", "Updated",
		),
	}
	tests := []struct {
		name string
		in   *rta.RTAPaymentRequest
		res  *rta.RTAPaymentResponse
		want *storage.RemitToAccountHistory
	}{
		{
			name: "Success Stage",
			in: &rta.RTAPaymentRequest{
				Partner:                     "Test Partner",
				ReferenceNumber:             "ABCDEFGHIJ",
				TrxDate:                     "022-06-10T16:59:51.515",
				AccountNumber:               "0663718000107",
				Currency:                    "1",
				ServiceCharge:               "100",
				Remarks:                     "Meow Meow",
				Particulars:                 "Transfer Particulars",
				MerchantName:                "Perahub",
				BankID:                      "30",
				LocationID:                  371,
				UserID:                      5188,
				CurrencyID:                  "1",
				CustomerID:                  "6925597",
				FormType:                    "OAR",
				FormNumber:                  "HOA0021942",
				TrxType:                     "",
				RemoteLocationID:            371,
				RemoteUserID:                5188,
				BillerName:                  "BPI",
				TrxTime:                     "16:59:51",
				TotalAmount:                 "",
				AccountName:                 "Naparate Jerica Reas",
				BeneficiaryAddress:          "1953 PH 3B BLOCK 6 LOT 9 CAMARIN, 175",
				BeneficiaryBirthdate:        "1996-08-10",
				BeneficiaryCity:             "UNIVERSITY OF THE PH",
				BeneficiaryCivil:            "S",
				BeneficiaryCountry:          "Philippines",
				BeneficiaryCustomertype:     "I",
				BeneficiaryFirstname:        "JERICA",
				BeneficiaryLastname:         "NAPARATE",
				BeneficiaryMiddlename:       "JERICA",
				BeneficiaryTin:              "000000000000000",
				BeneficiarySex:              "F",
				BeneficiaryState:            "QUEZON CITY",
				CurrencyCodePrincipalAmount: "PHP",
				PrincipalAmount:             "10",
				RecordType:                  "01",
				RemitterAddress:             "1953 PH 3B BLOCK 6 LOT 9 CAMARIN, 175",
				RemitterBirthdate:           "1996-08-10",
				RemitterCity:                "UNIVERSITY OF THE PH",
				RemitterCivil:               "S",
				RemitterCountry:             "PH",
				RemitterCustomerType:        "I",
				RemitterFirstname:           "Jerica",
				RemitterGender:              "F",
				RemitterID:                  6925597,
				RemitterLastname:            "Naparate",
				RemitterMiddlename:          "Reas",
				RemitterState:               "QUEZON CITY",
				SettlementMode:              "03",
				Notification:                "false",
				BeneZipCode:                 "1000",
			},
			res: &rta.RTAPaymentResponse{
				Message: "POSITIVE: BENEFICIARY ACCOUNT CREDITED",
				Result: &rta.RTAPaymentResult{
					BeneAmount:                     "205.58",
					BeneficiaryAddress:             "1953PH3B",
					BeneficiaryBankAccountno:       "0030555797",
					BeneficiaryCity:                "MAKATICITY",
					BeneficiaryCountry:             "PHILIPPINES",
					BeneficiaryFirstName:           "JERICA",
					BeneficiaryLastName:            "NAPARATE",
					BeneficiaryMiddleName:          "JERICA",
					BeneficiaryStateOrProvince:     "METROMANILA",
					BpiBranchCode:                  "1003",
					CurrencyCodeOfFundingAmount:    "PHP",
					CurrencyCodeOfSettlementAmount: "PHP",
					TxnDistributionDate:            "2021-05-18 00:00:00.0",
					FundingAmount:                  "205.58",
					Reason:                         "POSITIVE: BENEFICIARY ACCOUNT CREDITED",
					RemitterCity:                   "CALOOCANCITY",
					RemitterCountry:                "PHILIPPINES",
					RemitterFirstName:              "DEIB LOHR",
					RemitterLastName:               "ENRILE",
					RemitterMessageToBeneficiary:   "RemitterMessageToBeneficiary",
					RemitterMiddleName:             "DEL VALLE",
					RemitterStateOrProvince:        "METROMANILA",
					SettlementMode:                 "1",
					StatusCode:                     "8",
					TransactionDate:                "2021-05-25 00:00:00.0",
					TransactionReferenceNo:         "PHRBBPI10350",
					Message:                        "POSITIVE: BENEFICIARY ACCOUNT CREDITED",
					Code:                           "TS",
					SenderRefID:                    "SenderRefID",
					State:                          "Credited Beneficiary Account",
					UUID:                           "123",
					Description:                    "Successful transaction",
					Type:                           "CASH_IN",
					Amount:                         "1.00",
					UbpTranID:                      "56756756",
					TranRequestDate:                "2019-04-11",
					TranFinacleDate:                "2019-04-11",
					Created:                        &timestamppb.Timestamp{},
					Updated:                        &timestamppb.Timestamp{},
				},
				BankCode: "30",
			},
			want: &storage.RemitToAccountHistory{
				ID:                          orgID,
				OrgID:                       orgID,
				Partner:                     "Test Partner",
				ReferenceNumber:             "ABCDEFGHIJ",
				TrxDate:                     time.Now(),
				AccountNumber:               "0663718000107",
				Currency:                    "1",
				ServiceCharge:               "100",
				Remarks:                     "Meow Meow",
				Particulars:                 "Transfer Particulars",
				MerchantName:                "Perahub",
				BankID:                      30,
				LocationID:                  371,
				UserID:                      5188,
				CurrencyID:                  "1",
				CustomerID:                  "6925597",
				FormType:                    "OAR",
				FormNumber:                  "HOA0021942",
				TrxType:                     "",
				RemoteLocationID:            371,
				RemoteUserID:                5188,
				BillerName:                  "BPI",
				TrxTime:                     "16:59:51",
				TotalAmount:                 "",
				AccountName:                 "Naparate Jerica Reas",
				BeneficiaryAddress:          "1953 PH 3B BLOCK 6 LOT 9 CAMARIN, 175",
				BeneficiaryBirthDate:        "1996-08-10",
				BeneficiaryCity:             "UNIVERSITY OF THE PH",
				BeneficiaryCivil:            "S",
				BeneficiaryCountry:          "Philippines",
				BeneficiaryCustomerType:     "I",
				BeneficiaryFirstName:        "JERICA",
				BeneficiaryLastName:         "NAPARATE",
				BeneficiaryMiddleName:       "JERICA",
				BeneficiaryTin:              "000000000000000",
				BeneficiarySex:              "F",
				BeneficiaryState:            "QUEZON CITY",
				CurrencyCodePrincipalAmount: "PHP",
				PrincipalAmount:             "10",
				RecordType:                  "01",
				RemitterAddress:             "1953 PH 3B BLOCK 6 LOT 9 CAMARIN, 175",
				RemitterBirthDate:           "1996-08-10",
				RemitterCity:                "UNIVERSITY OF THE PH",
				RemitterCivil:               "S",
				RemitterCountry:             "PH",
				RemitterCustomerType:        "I",
				RemitterFirstName:           "Jerica",
				RemitterGender:              "F",
				RemitterID:                  6925597,
				RemitterLastName:            "Naparate",
				RemitterMiddleName:          "Reas",
				RemitterState:               "QUEZON CITY",
				SettlementMode:              "03",
				Notification:                false,
				BeneZipCode:                 "1000",
				Info:                        []byte{},
				Details:                     []byte{},
				TxnStatus:                   "SUCCESS",
				ErrorCode:                   "",
				ErrorMessage:                "",
				ErrorTime:                   "",
				ErrorType:                   "",
				CreatedBy:                   "",
				UpdatedBy:                   "",
				Created:                     time.Time{},
				Updated:                     time.Time{},
				SortByColumn:                "",
				SortOrder:                   "",
				Limit:                       0,
				Offset:                      0,
				Total:                       0,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			var err error
			aRes, aErr := RecordRTA(ctx, orgID, st, test.in, test.res, err)
			if aErr != nil {
				t.Fatal(aErr)
			}
			if !cmp.Equal(test.want, aRes, o) {
				t.Error(cmp.Diff(test.want, aRes, o))
			}
			uRes, uErr := UpdateRTA(ctx, orgID, st, &rta.RTARetryRequest{
				Partner:         test.in.Partner,
				ReferenceNumber: test.in.ReferenceNumber,
				LocationID:      test.in.LocationID,
				PrincipalAmount: test.in.PrincipalAmount,
				FormNumber:      test.in.FormNumber,
			}, &rta.RTARetryResponse{
				Message: test.res.Message,
				Result: &rta.RTAPaymentResult{
					BeneAmount:                     test.res.Result.BeneAmount,
					BeneficiaryAddress:             test.res.Result.BeneficiaryAddress,
					BeneficiaryBankAccountno:       test.res.Result.BeneficiaryBankAccountno,
					BeneficiaryCity:                test.res.Result.BeneficiaryCity,
					BeneficiaryCountry:             test.res.Result.BeneficiaryCountry,
					BeneficiaryFirstName:           test.res.Result.BeneficiaryFirstName,
					BeneficiaryLastName:            test.res.Result.BeneficiaryLastName,
					BeneficiaryMiddleName:          test.res.Result.BeneficiaryMiddleName,
					BeneficiaryStateOrProvince:     test.res.Result.BeneficiaryStateOrProvince,
					BpiBranchCode:                  test.res.Result.BpiBranchCode,
					CurrencyCodeOfFundingAmount:    test.res.Result.CurrencyCodeOfFundingAmount,
					CurrencyCodeOfSettlementAmount: test.res.Result.CurrencyCodeOfSettlementAmount,
					TxnDistributionDate:            test.res.Result.TxnDistributionDate,
					FundingAmount:                  test.res.Result.FundingAmount,
					Reason:                         test.res.Result.Reason,
					RemitterCity:                   test.res.Result.RemitterCity,
					RemitterCountry:                test.res.Result.RemitterCountry,
					RemitterFirstName:              test.res.Result.RemitterFirstName,
					RemitterLastName:               test.res.Result.RemitterLastName,
					RemitterMessageToBeneficiary:   test.res.Result.RemitterMessageToBeneficiary,
					RemitterMiddleName:             test.res.Result.RemitterMiddleName,
					RemitterStateOrProvince:        test.res.Result.RemitterStateOrProvince,
					SettlementMode:                 test.res.Result.SettlementMode,
					StatusCode:                     test.res.Result.StatusCode,
					TransactionDate:                test.res.Result.TransactionDate,
					TransactionReferenceNo:         test.res.Result.TransactionReferenceNo,
					Message:                        test.res.Result.Message,
					Code:                           test.res.Result.Code,
					SenderRefID:                    test.res.Result.SenderRefID,
					State:                          test.res.Result.State,
					UUID:                           test.res.Result.UUID,
					Description:                    test.res.Result.Description,
					Type:                           test.res.Result.Type,
					Amount:                         test.res.Result.Amount,
					UbpTranID:                      test.res.Result.UbpTranID,
					TranRequestDate:                test.res.Result.TranRequestDate,
					TranFinacleDate:                test.res.Result.TranFinacleDate,
					Created:                        &timestamppb.Timestamp{},
					Updated:                        &timestamppb.Timestamp{},
				},
				BankCode: "",
			}, err)
			if uErr != nil {
				t.Fatal(uErr)
			}
			if !cmp.Equal(test.want, uRes, o) {
				t.Error(cmp.Diff(test.want, uRes, o))
			}
		})
	}
}
