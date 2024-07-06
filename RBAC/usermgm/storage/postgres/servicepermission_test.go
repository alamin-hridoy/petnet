package postgres

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"brank.as/rbac/usermgm/storage"
)

func TestUpsertServicePermission(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing DATABASE_CONNECTION env")
	}
	db, clean := NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)

	svc, err := db.UpsertService(context.Background(), storage.Service{
		Name:    "TestUpsertServicePermission",
		Default: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		in   storage.ServicePermission
		err  error
	}{
		{
			name: "Success",
			in: storage.ServicePermission{
				ServiceID: svc.ID,
				Resource:  "RBCA:role",
				Action:    "create",
			},
		},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := db.UpsertServicePermission(ctx, tst.in)
			if !cmp.Equal(tst.err, err) {
				t.Error(cmp.Diff(tst.err, err))
			}
			if got.ID == "" {
				t.Error("want ID, got:", got)
			}
		})
	}
}

func TestGetServicePermission(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing DATABASE_CONNECTION env")
	}
	db, clean := NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)

	id, err := db.UpsertService(context.Background(), storage.Service{
		Name:    "TestGetServicePermission",
		Default: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		in   storage.ServicePermission
		err  error
	}{
		{
			name: "Success",
			in: storage.ServicePermission{
				ServiceID: id.ID,
				Resource:  "RBCA:role",
				Action:    "create",
			},
		},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := db.UpsertServicePermission(ctx, tst.in)
			if !cmp.Equal(tst.err, err) {
				t.Error(cmp.Diff(tst.err, err))
			}
			if got.ID == "" {
				t.Error("want ID, got:", got)
			}
			t.Cleanup(func() {
				db.DeleteServicePermission(ctx, got.ID)
			})

			sp, err := db.GetServicePermission(ctx, got.ID)
			if !cmp.Equal(tst.err, err) {
				t.Error(cmp.Diff(tst.err, err))
			}

			if !cmp.Equal(sp, got) {
				t.Error(cmp.Diff(sp, got))
			}
		})
	}
}

func TestDeleteServicePermission(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing DATABASE_CONNECTION env")
	}
	db, clean := NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)

	svc, err := db.UpsertService(context.Background(), storage.Service{
		Name:    "TestDeleteServicePermission",
		Default: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		in   storage.ServicePermission
		err  error
	}{
		{
			name: "Success",
			in: storage.ServicePermission{
				ServiceID: svc.ID,
				Resource:  "RBCA:role",
				Action:    "create",
			},
		},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := db.UpsertServicePermission(ctx, tst.in)
			if !cmp.Equal(tst.err, err) {
				t.Error(cmp.Diff(tst.err, err))
			}
			if got.ID == "" {
				t.Error("want ID, got:", got)
			}
			t.Cleanup(func() {
				db.DeleteServicePermission(ctx, got.ID)
			})

			err = db.DeleteServicePermission(ctx, got.ID)
			if !cmp.Equal(tst.err, err) {
				t.Error(cmp.Diff(tst.err, err))
			}
		})
	}
}
