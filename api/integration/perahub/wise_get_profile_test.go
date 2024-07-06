package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWISEGetProfile(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)
	wantReq := WISEGetProfileReq{
		Email: "user@email.com",
	}
	wantResp := &WISEGetProfileResp{
		Type: "personal",
		Details: WISEGetPFDetails{
			FirstName:   "Brankas",
			LastName:    "Sender",
			BirthDate:   "1990-01-10",
			PhoneNumber: "+639999999999",
			Occupations: []WISEOccupation{
				{
					Code:   "Software Engineer",
					Format: "FREE_FORM",
				},
			},
			PrimaryAddress: json.Number("37325852"),
		},
		ProfileID: json.Number("16325688"),
		Address: WISEPFAddress{
			Country:   "ph",
			FirstLine: "East Offices Bldg., 114 Aguirre St.,Legaspi Village,",
			PostCode:  "1229",
			City:      "Makati",
		},
	}
	gotResp, err := s.WISEGetProfile(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}

	var gotReq WISEGetProfileReq
	if err := json.Unmarshal(m.GetMockRequest(), &gotReq); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(wantReq, gotReq) {
		t.Error(cmp.Diff(wantReq, gotReq))
	}
	if !cmp.Equal(wantResp, gotResp) {
		t.Error(cmp.Diff(wantResp, gotResp))
	}
}
