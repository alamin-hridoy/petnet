package util

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestRecordRemittanceHistory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	st := newTestStorage(t)
	dsaid := uuid.NewString()
	uid := uuid.NewString()
	var err error
	var bj []byte
	bj, err = json.Marshal(&storage.PerahubRemittanceHistoryDetails{
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
		SenderOfFund:           "SALARY",
		SenderOccupation:       "DOCTOR",
		SenderEmploymentNature: "IT",
		SendPartnerCode:        "USP",
		PayoutPartnerCode:      "USP",
		PartnerCode:            "USP",
	})
	if err != nil {
		bj = []byte{}
	}

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.PerahubRemittanceHistory{}, "Details", "Total", "RemittanceHistoryID", "TxnCreatedTime", "ErrorTime", "PayHisErr", "ErrorCode", "ErrorMessage", "ErrorType", "TxnUpdatedTime", "TxnConfirmTime",
		),
	}
	tests := []struct {
		name string
		want *storage.PerahubRemittanceHistory
		in   *storage.PerahubRemittanceHistory
	}{
		{
			name: "success validate send",
			in: &storage.PerahubRemittanceHistory{
				DsaID:                       dsaid,
				UserID:                      uid,
				SendValidateReferenceNumber: "1111111111111",
				Phrn:                        "p1p1p1p1p1p1p1",
				TxnStatus:                   storage.VALIDATE_SEND,
				Details:                     bj,
				TxnCreatedTime:              time.Now(),
			},
			want: &storage.PerahubRemittanceHistory{
				DsaID:                       dsaid,
				UserID:                      uid,
				SendValidateReferenceNumber: "1111111111111",
				Phrn:                        "p1p1p1p1p1p1p1",
				TxnStatus:                   storage.VALIDATE_SEND,
				Details:                     bj,
				TxnCreatedTime:              time.Now(),
			},
		},
		{
			name: "error validate send",
			in: &storage.PerahubRemittanceHistory{
				DsaID:                       dsaid,
				UserID:                      uid,
				SendValidateReferenceNumber: "1111111111111",
				Phrn:                        "p1p1p1p1p1p1p1",
				TxnStatus:                   storage.VALIDATE_SEND,
				Details:                     []byte{},
				TxnCreatedTime:              time.Now(),
				PayHisErr: &perahub.Error{
					Code:       "503",
					GRPCCode:   503,
					Msg:        "failed to create validate send",
					UnknownErr: "failed to create validate send",
					Type:       perahub.RemitanceError,
					Errors: map[string][]string{
						"dsss": {"failed to create validate send"},
					},
				},
			},
			want: &storage.PerahubRemittanceHistory{
				DsaID:                       dsaid,
				UserID:                      uid,
				SendValidateReferenceNumber: "1111111111111",
				Phrn:                        "p1p1p1p1p1p1p1",
				TxnStatus:                   storage.TRANSACTION_FAIL,
				Details:                     []byte{},
				TxnCreatedTime:              time.Now(),
			},
		},
		{
			name: "success confirm send",
			in: &storage.PerahubRemittanceHistory{
				DsaID:                       dsaid,
				UserID:                      uid,
				Phrn:                        "p1p1p1p1p1p1p1",
				SendValidateReferenceNumber: "1111111111111",
				TxnStatus:                   storage.CONFIRM_SEND,
			},
			want: &storage.PerahubRemittanceHistory{
				DsaID:                       dsaid,
				UserID:                      uid,
				Phrn:                        "p1p1p1p1p1p1p1",
				SendValidateReferenceNumber: "1111111111111",
				TxnStatus:                   storage.CONFIRM_SEND,
			},
		},
		{
			name: "success validate receive",
			in: &storage.PerahubRemittanceHistory{
				DsaID:                         dsaid,
				UserID:                        uid,
				Phrn:                          "p1p1p1p1p1p1p1",
				SendValidateReferenceNumber:   "1111111111111",
				PayoutValidateReferenceNumber: "PayoutValidateReferenceNumber",
				TxnStatus:                     storage.VALIDATE_RECEIVE,
			},
			want: &storage.PerahubRemittanceHistory{
				DsaID:                         dsaid,
				UserID:                        uid,
				Phrn:                          "p1p1p1p1p1p1p1",
				SendValidateReferenceNumber:   "1111111111111",
				PayoutValidateReferenceNumber: "PayoutValidateReferenceNumber",
				TxnStatus:                     storage.VALIDATE_RECEIVE,
			},
		},
		{
			name: "success confirm receive",
			in: &storage.PerahubRemittanceHistory{
				DsaID:                         dsaid,
				UserID:                        uid,
				Phrn:                          "p1p1p1p1p1p1p1",
				SendValidateReferenceNumber:   "1111111111111",
				PayoutValidateReferenceNumber: "PayoutValidateReferenceNumber",
				TxnStatus:                     storage.CONFIRM_RECEIVE,
			},
			want: &storage.PerahubRemittanceHistory{
				DsaID:                         dsaid,
				UserID:                        uid,
				Phrn:                          "p1p1p1p1p1p1p1",
				SendValidateReferenceNumber:   "1111111111111",
				PayoutValidateReferenceNumber: "PayoutValidateReferenceNumber",
				TxnStatus:                     storage.CONFIRM_RECEIVE,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, err := RecordRemittanceHistory(ctx, st, *test.in)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.want, got, o) {
				t.Error(cmp.Diff(test.want, got, o))
			}
		})
	}
}
