package postgres_test

import (
	"context"
	"encoding/json"
	"testing"

	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	"brank.as/petnet/profile/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestEventData(t *testing.T) {
	ts := newTestStorage(t)

	ctx := context.Background()
	oid := uuid.NewString()
	_, err := ts.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uuid.NewString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	eid := uuid.New().String()
	d := &storage.Branch{
		OrgID:        oid,
		OrgProfileID: uuid.New().String(),
		Title:        "title",
		BranchAddress: storage.BranchAddress{
			Address1:   "addr1",
			City:       "city",
			State:      "state",
			PostalCode: "12345",
		},
	}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	want := &storage.EventData{
		EventID:  eid,
		Resource: "branch",
		Action:   tpb.ActionType_Create.String(),
		Data:     string(b),
	}

	o := []cmp.Option{
		cmpopts.IgnoreFields(storage.EventData{}, "Created", "Updated", "Data"),
	}
	if err := ts.CreateEventData(ctx, want); err != nil {
		t.Fatal(err)
	}
	got, err := ts.GetEventData(ctx, eid)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got, o...) {
		t.Error(cmp.Diff(want, got, o...))
	}
	if got.Created.IsZero() {
		t.Error("created is empty")
	}
	if got.Updated.IsZero() {
		t.Error("updated is empty")
	}

	got2 := &storage.Branch{}
	if err := json.Unmarshal([]byte(got.Data), got2); err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(d, got2) {
		t.Error(cmp.Diff(d, got2))
	}

	if err = ts.DeleteEventData(ctx, eid); err != nil {
		t.Fatal(err)
	}
	if _, err := ts.GetEventData(ctx, eid); err != storage.NotFound {
		t.Fatal("event data not deleted")
	}
}
