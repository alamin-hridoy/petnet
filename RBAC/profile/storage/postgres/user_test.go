package postgres

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"brank.as/rbac/profile/storage"
)

func TestCRDUser(t *testing.T) {
	ts := newTestStorage(t)
	want := &storage.User{
		OrgID: uuid.New().String(),
	}
	uid, err := ts.CreateUser(context.TODO(), *want)
	if err != nil {
		t.Fatal("CreateUser: ", err)
	}
	if uid == "" {
		t.Fatal("CreateUser, userID should not be empty")
	}

	got, err := ts.GetUserByID(context.TODO(), uid)
	if err != nil {
		t.Fatal("GetUserByID: ", err)
	}
	tOps := cmpopts.IgnoreFields(storage.User{}, "ID", "Created", "Updated", "Deleted")
	if !cmp.Equal(want, got, tOps) {
		t.Fatal("GetUserByID: ", cmp.Diff(want, got, tOps))
	}
	if got.ID == "" {
		t.Fatal("GetUserByID, userID should not be empty")
	}
	if got.Created.IsZero() || got.Updated.IsZero() {
		t.Fatal("GetUserByID, created and updated shouldn't be empty")
	}
	if got.Deleted.Valid {
		t.Fatal("GetUserByID, deleted should be null")
	}

	want.ID = got.ID
	tOps = cmpopts.IgnoreFields(storage.User{}, "Created", "Updated", "Deleted")
	got2, err := ts.GetUsersByOrg(context.TODO(), got.OrgID)
	if err != nil {
		t.Fatal("GetUsersByOrg: ", err)
	}
	if !cmp.Equal(want, &got2[0], tOps) {
		t.Fatal("GetUsersByOrg: ", cmp.Diff(want, got2[0], tOps))
	}
	if got.Created.IsZero() || got.Updated.IsZero() {
		t.Fatal("GetUsersByOrg, created and updated shouldn't be empty")
	}
	if got.Deleted.Valid {
		t.Fatal("GetUserByOrg, deleted should be null")
	}

	got3, err := ts.DeleteUserByID(context.TODO(), got2[0].ID)
	if err != nil {
		t.Fatal("DeleteUserByID: ", err)
	}
	if !cmp.Equal(want, got3, tOps) {
		t.Fatal("DeleteUserByID: ", cmp.Diff(want, got3, tOps))
	}
	if !got3.Deleted.Valid {
		t.Fatal("DeleteUserByID, Deleted time should not be null")
	}
}
