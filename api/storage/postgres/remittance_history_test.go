package postgres

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"brank.as/petnet/api/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestRemittanceHistory(t *testing.T) {
	ts := newTestStorage(t)
	var err error
	var bj []byte
	dsaid := uuid.NewString()
	uid := uuid.NewString()
	bj, err = json.Marshal(storage.PerahubRemittanceHistoryDetails{
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
	})
	if err != nil {
		bj = []byte{}
	}
	t.Run("RemittanceHistory", func(t *testing.T) {
		t.Parallel()
		in := &storage.PerahubRemittanceHistory{
			RemittanceHistoryID:         uuid.NewString(),
			DsaID:                       dsaid,
			UserID:                      uid,
			SendValidateReferenceNumber: "1111111111111",
			TxnStatus:                   "VALIDATE_SEND",
			Details:                     bj,
			TxnCreatedTime:              time.Now(),
		}
		r, err := ts.CreateRemittanceHistory(context.TODO(), *in)
		if err != nil {
			t.Fatalf("Create Remittance History = got error %v, want nil", err)
		}
		if r.RemittanceHistoryID == "" {
			t.Fatal("Create Remittance History = returned empty ID")
		}
		r.Phrn = "p2p2p2p2p2p2p"
		updateResult, err := ts.UpdateRemittanceHistory(context.TODO(), *r)
		if err != nil {
			t.Fatalf("Update Remittance History = got error %v, want nil", err)
		}
		if updateResult.Phrn == "" {
			t.Fatal("Update Remittance History = returned empty Phrn")
		}
		if updateResult.Phrn != r.Phrn {
			t.Fatal("Update Remittance History not match")
		}
		gr, err := ts.GetRemittanceHistory(context.TODO(), r.RemittanceHistoryID)
		if err != nil {
			t.Fatalf("Get Remittance History = got error %v, want nil", err)
		}
		opt := cmpopts.IgnoreFields(storage.PerahubRemittanceHistory{}, "Details", "Total")
		if !cmp.Equal(*r, *gr, opt) {
			t.Fatal(cmp.Diff(*r, *gr, opt))
		}
		lph, err := ts.ListRemittanceHistory(context.TODO(), storage.RemittanceHistoryFilter{})
		if err != nil {
			t.Fatalf("List Remittance History = got error %v, want nil", err)
		}
		if !cmp.Equal([]storage.PerahubRemittanceHistory{*r}, lph, opt) {
			t.Fatal(cmp.Diff([]storage.PerahubRemittanceHistory{*r}, lph, opt))
		}
		cph, err := ts.ConfirmRemittanceHistory(context.TODO(), *r)
		if err != nil {
			t.Fatalf("Confirm Remittance History = got error %v, want nil", err)
		}

		if cph.TxnStatus == "" {
			t.Fatal("Confirm Remittance History = returned empty TxnStatus")
		}
		if cph.TxnStatus != "CONFIRM_SEND" {
			t.Fatal("Confirm Remittance History not match the CONFIRM_SEND")
		}
		r.CancelSendReferenceNumber = "2222222334235345556654"
		cphd, err := ts.CancelRemittanceHistory(context.TODO(), *r)
		if err != nil {
			t.Fatalf("Cancel Remittance History = got error %v, want nil", err)
		}
		if cphd.TxnStatus == "" {
			t.Fatal("Cancel Remittance History = returned empty TxnStatus")
		}
		if cphd.TxnStatus != "CANCEL_SEND" {
			t.Fatal("Cancel Remittance History not match the CANCEL_SEND")
		}
		r.PayoutValidateReferenceNumber = "ghesds4523523654356242366"
		vrph, err := ts.ValidateReceiveRemittanceHistory(context.TODO(), *r)
		if err != nil {
			t.Fatalf("Validate Receive Remittance History = got error %v, want nil", err)
		}
		if vrph.TxnStatus == "" {
			t.Fatal("Cancel Remittance History = returned empty TxnStatus")
		}
		if vrph.TxnStatus != "VALIDATE_RECEIVE" {
			t.Fatal("Cancel Remittance History not match the VALIDATE_RECEIVE")
		}
		crph, err := ts.ConfirmReceiveRemittanceHistory(context.TODO(), *r)
		if err != nil {
			t.Fatalf("Confirm Receive Remittance History = got error %v, want nil", err)
		}
		if crph.TxnStatus == "" {
			t.Fatal("Cancel Remittance History = returned empty TxnStatus")
		}

		if crph.TxnStatus != "CONFIRM_RECEIVE" {
			t.Fatal("Cancel Remittance History not match the CONFIRM_RECEIVE")
		}
	})
}
