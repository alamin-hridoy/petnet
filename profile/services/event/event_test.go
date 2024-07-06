package event

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"

	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	ec "brank.as/petnet/profile/core/event"
)

func TestService(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	test.NewNullLogger()
	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)

	oid := uuid.New().String()
	_, err := st.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uuid.NewString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	eid := uuid.New().String()
	want := &tpb.EventData{
		EventID:  eid,
		Resource: "branch",
		Action:   tpb.ActionType_Create,
		Data:     `{"data": "branch"}`,
	}

	s := New(ec.New(st))
	if _, err := s.CreateEventData(ctx, &tpb.CreateEventDataRequest{EventData: want}); err != nil {
		t.Fatal("create event data: ", err)
	}

	o := cmp.Options{
		cmpopts.IgnoreUnexported(
			tpb.EventData{},
		),
	}
	got, err := s.GetEventData(ctx, &tpb.GetEventDataRequest{EventID: eid})
	if err != nil {
		t.Fatal("get event data: ", err)
	}
	if !cmp.Equal(want, got.EventData, o) {
		t.Error("(-want +got): ", cmp.Diff(want, got.EventData, o))
	}
}
