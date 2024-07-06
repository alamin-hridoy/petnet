package util

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	"github.com/bojanz/currency"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

func TestErrors(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)

	zamt, _ := currency.NewMinor("0", "PHP")
	rmt := core.Remittance{
		DsaID:          uuid.New().String(),
		UserID:         uuid.New().String(),
		GeneratedRefNo: "genref",
		ControlNo:      "refno",
		MyWUNumber:     "wuno",
		RemitPartner:   "RIA",
		GrossTotal:     zamt,
		SourceAmount:   zamt,
		DestAmount:     zamt,
		Tax:            zamt,
		Charge:         zamt,
	}

	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.RemitHistory{}, "TxnID", "DsaID", "DsaOrderID", "UserID", "RemcoID", "RemType", "SenderID", "ReceiverID", "RemcoControlNo", "Remittance", "TxnStagedTime", "TxnCompletedTime", "Updated", "Total", "ErrorTime",
		),
	}
	tests := []struct {
		name     string
		want     *storage.RemitHistory
		stgOpts  *StageTxnOpts
		confOpts *ConfirmTxnOpts
	}{
		{
			name: "Success Stage",
			want: &storage.RemitHistory{
				TxnStatus: string(storage.SuccessStatus),
				TxnStep:   string(storage.StageStep),
				TransactionType: sql.NullString{
					String: PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			stgOpts: &StageTxnOpts{
				TxnType:     storage.DisburseType,
				PtnrRemType: "PO",
			},
		},
		{
			name: "Nonex Error Stage",
			want: &storage.RemitHistory{
				TxnStatus: string(storage.FailStatus),
				TxnStep:   string(storage.StageStep),
				ErrorCode: "200",
				ErrorMsg:  "<html>server-error<html>",
				ErrorType: string(perahub.NonexError),
				TransactionType: sql.NullString{
					String: PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			stgOpts: &StageTxnOpts{
				TxnType:     storage.DisburseType,
				PtnrRemType: "PO",
				TxnErr: &perahub.Error{
					Code:       "200",
					GRPCCode:   codes.Internal,
					Msg:        codes.Internal.String(),
					UnknownErr: "<html>server-error<html>",
					Type:       perahub.NonexError,
				},
			},
		},
		{
			name: "DRP Error Stage",
			want: &storage.RemitHistory{
				TxnStatus: string(storage.FailStatus),
				TxnStep:   string(storage.StageStep),
				ErrorCode: "400",
				ErrorMsg:  `{"non":{"standard":"error"}}`,
				ErrorType: string(perahub.DRPError),
				TransactionType: sql.NullString{
					String: PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			stgOpts: &StageTxnOpts{
				TxnType:     storage.DisburseType,
				PtnrRemType: "PO",
				TxnErr: &perahub.Error{
					Code:       "400",
					GRPCCode:   codes.Internal,
					Msg:        codes.Internal.String(),
					UnknownErr: `{"non":{"standard":"error"}}`,
					Type:       perahub.DRPError,
				},
			},
		},
		{
			name: "Partner Error Stage",
			want: &storage.RemitHistory{
				TxnStatus: string(storage.FailStatus),
				TxnStep:   string(storage.StageStep),
				ErrorCode: "400",
				ErrorMsg:  codes.Internal.String(),
				ErrorType: string(perahub.PartnerError),
				TransactionType: sql.NullString{
					String: PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			stgOpts: &StageTxnOpts{
				TxnType:     storage.DisburseType,
				PtnrRemType: "PO",
				TxnErr: &perahub.Error{
					Code:     "400",
					GRPCCode: codes.Internal,
					Msg:      codes.Internal.String(),
					Type:     perahub.PartnerError,
				},
			},
		},
		{
			name: "Internal Error Stage",
			want: &storage.RemitHistory{
				TxnStatus: string(storage.FailStatus),
				TxnStep:   string(storage.StageStep),
				ErrorCode: "",
				ErrorMsg:  "some error",
				ErrorType: string(perahub.DRPError),
				TransactionType: sql.NullString{
					String: PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			stgOpts: &StageTxnOpts{
				TxnType:     storage.DisburseType,
				PtnrRemType: "PO",
				TxnErr:      errors.New("some error"),
			},
		},
		{
			name: "Success Confirm",
			want: &storage.RemitHistory{
				TxnStatus: string(storage.SuccessStatus),
				TxnStep:   string(storage.ConfirmStep),
				TransactionType: sql.NullString{
					String: PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			confOpts: &ConfirmTxnOpts{
				TxnType:     storage.DisburseType,
				PtnrRemType: "PO",
			},
		},
		{
			name: "Nonex Error Confirm",
			want: &storage.RemitHistory{
				TxnStatus: string(storage.FailStatus),
				TxnStep:   string(storage.ConfirmStep),
				ErrorCode: "200",
				ErrorMsg:  "<html>server-error<html>",
				ErrorType: string(perahub.NonexError),
				TransactionType: sql.NullString{
					String: PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			confOpts: &ConfirmTxnOpts{
				TxnType:     storage.DisburseType,
				PtnrRemType: "PO",
				TxnErr: &perahub.Error{
					Code:       "200",
					GRPCCode:   codes.Internal,
					Msg:        codes.Internal.String(),
					UnknownErr: "<html>server-error<html>",
					Type:       perahub.NonexError,
				},
			},
		},
		{
			name: "DRP Error Confirm",
			want: &storage.RemitHistory{
				TxnStatus: string(storage.FailStatus),
				TxnStep:   string(storage.ConfirmStep),
				ErrorCode: "400",
				ErrorMsg:  `{"non":{"standard":"error"}}`,
				ErrorType: string(perahub.DRPError),
				TransactionType: sql.NullString{
					String: PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			confOpts: &ConfirmTxnOpts{
				TxnType:     storage.DisburseType,
				PtnrRemType: "PO",
				TxnErr: &perahub.Error{
					Code:       "400",
					GRPCCode:   codes.Internal,
					Msg:        codes.Internal.String(),
					UnknownErr: `{"non":{"standard":"error"}}`,
					Type:       perahub.DRPError,
				},
			},
		},
		{
			name: "Partner Error Confirm",
			want: &storage.RemitHistory{
				TxnStatus: string(storage.FailStatus),
				TxnStep:   string(storage.ConfirmStep),
				ErrorCode: "400",
				ErrorMsg:  codes.Internal.String(),
				ErrorType: string(perahub.PartnerError),
				TransactionType: sql.NullString{
					String: PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			confOpts: &ConfirmTxnOpts{
				TxnType:     storage.DisburseType,
				PtnrRemType: "PO",
				TxnErr: &perahub.Error{
					Code:     "400",
					GRPCCode: codes.Internal,
					Msg:      codes.Internal.String(),
					Type:     perahub.PartnerError,
				},
			},
		},
		{
			name: "Internal Error Confirm",
			want: &storage.RemitHistory{
				TxnStatus: string(storage.FailStatus),
				TxnStep:   string(storage.ConfirmStep),
				ErrorCode: "",
				ErrorMsg:  "some error",
				ErrorType: string(perahub.DRPError),
				TransactionType: sql.NullString{
					String: PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			confOpts: &ConfirmTxnOpts{
				TxnType:     storage.DisburseType,
				PtnrRemType: "PO",
				TxnErr:      errors.New("some error"),
			},
		},
		{
			name: "Stage MultiError",
			want: &storage.RemitHistory{
				TxnStatus: string(storage.FailStatus),
				TxnStep:   string(storage.StageStep),
				ErrorCode: "400",
				ErrorMsg:  `Internal, {"errors":["error1","error2"]}`,
				ErrorType: string(perahub.PartnerError),
				TransactionType: sql.NullString{
					String: PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			stgOpts: &StageTxnOpts{
				TxnType:     storage.DisburseType,
				PtnrRemType: "PO",
				TxnErr: &perahub.Error{
					Code:     "400",
					GRPCCode: codes.Internal,
					Msg:      codes.Internal.String(),
					Type:     perahub.PartnerError,
					Errors: map[string][]string{
						"errors": {
							"error1",
							"error2",
						},
					},
				},
			},
		},
		{
			name: "Confirm MultiError",
			want: &storage.RemitHistory{
				TxnStatus: string(storage.FailStatus),
				TxnStep:   string(storage.ConfirmStep),
				ErrorCode: "400",
				ErrorMsg:  `Internal, {"errors":["error1","error2"]}`,
				ErrorType: string(perahub.PartnerError),
				TransactionType: sql.NullString{
					String: PerahubTrxTypeOTC,
					Valid:  true,
				},
			},
			confOpts: &ConfirmTxnOpts{
				TxnType:     storage.DisburseType,
				PtnrRemType: "PO",
				TxnErr: &perahub.Error{
					Code:     "400",
					GRPCCode: codes.Internal,
					Msg:      codes.Internal.String(),
					Type:     perahub.PartnerError,
					Errors: map[string][]string{
						"errors": {
							"error1",
							"error2",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rmt.DsaOrderID = uuid.New().String()
			ctx := context.Background()

			var res *storage.RemitHistory
			if test.stgOpts == nil {
				test.stgOpts = &StageTxnOpts{
					TxnType:     storage.DisburseType,
					PtnrRemType: "PO",
				}
			}
			res, err := RecordStageTxn(ctx, st, rmt, *test.stgOpts)
			if err != nil {
				t.Fatal(err)
			}

			if test.confOpts == nil {
				got, err := st.GetRemitHistory(ctx, res.TxnID)
				if err != nil {
					t.Fatal(err)
				}

				if test.stgOpts.TxnErr != nil && got.ErrorTime.Time.IsZero() {
					t.Error("ErrorTime should be set")
				}
				if got.TxnStagedTime.Time.IsZero() {
					t.Error("TxnStagedTime should be set")
				}
				if !cmp.Equal(test.want, got, o) {
					t.Error(cmp.Diff(test.want, got, o))
				}
				return
			}

			test.confOpts.TxnID = res.TxnID
			res, err = RecordConfirmTxn(ctx, st, rmt, *test.confOpts)
			if err != nil {
				t.Fatal(err)
			}

			got, err := st.GetRemitHistory(ctx, res.TxnID)
			if err != nil {
				t.Fatal(err)
			}

			if test.confOpts.TxnErr != nil && got.ErrorTime.Time.IsZero() {
				t.Error("ErrorTime should be set")
			}
			if got.TxnStagedTime.Time.IsZero() {
				t.Error("TxnStagedTime should be set")
			}
			if got.TxnCompletedTime.Time.IsZero() {
				t.Error("TxnCompletedTime should be set")
			}
			if !cmp.Equal(test.want, got, o) {
				t.Error(cmp.Diff(test.want, got, o))
			}
		})
	}
}
