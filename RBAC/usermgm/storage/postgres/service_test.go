package postgres

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"brank.as/rbac/usermgm/storage"
)

func TestUpsertService(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing DATABASE_CONNECTION env")
	}
	db, clean := NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)

	tests := []struct {
		name string
		in   storage.Service
		err  error
	}{
		{
			name: "Success",
			in: storage.Service{
				Name: "TestUpsertService",
			},
		},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := db.UpsertService(ctx, tst.in)
			if !cmp.Equal(tst.err, err) {
				t.Error(cmp.Diff(tst.err, err))
			}
			if got.ID == "" {
				t.Error("want ID, got:", got)
			}
		})
	}
}

func TestGetService(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing DATABASE_CONNECTION env")
	}
	db, clean := NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)

	tests := []struct {
		name string
		in   storage.Service
		err  error
	}{
		{
			name: "Success",
			in: storage.Service{
				Name: "TestGetService",
			},
		},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := db.UpsertService(ctx, tst.in)
			if !cmp.Equal(tst.err, err) {
				t.Error(cmp.Diff(tst.err, err))
			}
			if got.ID == "" {
				t.Error("want ID, got:", got)
			}
			t.Cleanup(func() {
				db.DeleteService(ctx, got.ID)
			})

			sp, err := db.GetService(ctx, got.ID)
			if !cmp.Equal(tst.err, err) {
				t.Error(cmp.Diff(tst.err, err))
			}

			if !cmp.Equal(sp, got) {
				t.Error(cmp.Diff(sp, got))
			}
		})
	}
}

func TestDeleteService(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing DATABASE_CONNECTION env")
	}
	db, clean := NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)

	tests := []struct {
		name string
		in   storage.Service
		err  error
	}{
		{
			name: "Success",
			in: storage.Service{
				Name:    "TestDeleteService",
				Default: false,
			},
		},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := db.UpsertService(ctx, tst.in)
			if !cmp.Equal(tst.err, err) {
				t.Error(cmp.Diff(tst.err, err))
			}
			if got.ID == "" {
				t.Error("want ID, got:", got)
			}
			t.Cleanup(func() {
				db.DeleteService(ctx, got.ID)
			})

			err = db.DeleteService(ctx, got.ID)
			if !cmp.Equal(tst.err, err) {
				t.Error(cmp.Diff(tst.err, err))
			}
		})
	}
}
