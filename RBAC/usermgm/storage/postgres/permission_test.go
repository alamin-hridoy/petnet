package postgres

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"brank.as/rbac/usermgm/storage"
)

func TestCreatePermission(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing DATABASE_CONNECTION env")
	}
	db, clean := NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)
	t.Parallel()
	opt := cmp.FilterValues(func(a, b time.Time) bool {
		return true
	}, cmp.Comparer(func(a, b time.Time) bool {
		return a.Sub(b) < 10*time.Second && b.Sub(a) < 10*time.Second
	}))
	testOrg, testCreator := uuid.New().String(), uuid.New().String()

	pid := uuid.New().String()
	tests := []struct {
		in   storage.Permission
		want *storage.Permission
		err  error
	}{
		{
			in: storage.Permission{
				ID:          "some-keto-id1",
				SvcPermID:   pid,
				Name:        "test1",
				OrgID:       testOrg,
				Description: "First test permission",
				CreateUID:   testCreator,
			},
			want: &storage.Permission{
				ID:          "some-keto-id1",
				SvcPermID:   pid,
				Name:        "test1",
				OrgID:       testOrg,
				Description: "First test permission",
				CreateUID:   testCreator,
				Created:     time.Now(),
				Updated:     time.Now(),
			},
		},
	}
	for _, tst := range tests {
		t.Run(tst.in.Name, func(t *testing.T) {
			ctx := context.Background()
			got, err := db.CreatePermission(ctx, tst.in)
			if !cmp.Equal(tst.err, err) {
				t.Error(cmp.Diff(tst.err, err))
			}
			if !cmp.Equal(tst.want, got, opt) {
				t.Error(cmp.Diff(tst.want, got, opt))
			}
		})
	}
}

func TestDeletePermission(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing DATABASE_CONNECTION env")
	}
	opt := cmp.FilterValues(func(a, b time.Time) bool {
		return true
	}, cmp.Comparer(func(a, b time.Time) bool {
		return a.Sub(b) < 10*time.Second && b.Sub(a) < 10*time.Second
	}))
	testOrg, testCreator := uuid.New().String(), uuid.New().String()
	pid := uuid.New().String()
	tests := []struct {
		in   storage.Permission
		want *storage.Permission
		err  error
	}{
		{
			in: storage.Permission{
				ID:          "some-keto-id1",
				SvcPermID:   pid,
				Name:        "test1",
				OrgID:       testOrg,
				Description: "First test permission",
				CreateUID:   testCreator,
				DeleteUID:   sql.NullString{Valid: true, String: testCreator},
			},
			want: &storage.Permission{
				ID:          "some-keto-id1",
				SvcPermID:   pid,
				Name:        "test1",
				OrgID:       testOrg,
				Description: "First test permission",
				CreateUID:   testCreator,
				DeleteUID:   sql.NullString{Valid: true, String: testCreator},
				Created:     time.Now(),
				Updated:     time.Now(),
			},
		},
	}
	db, clean := NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)
	for _, tst := range tests {
		t.Run(tst.in.Name, func(t *testing.T) {
			ctx := context.Background()
			got, err := db.CreatePermission(ctx, tst.in)
			if !cmp.Equal(tst.err, err) {
				t.Error(cmp.Diff(tst.err, err))
			}
			if !cmp.Equal(tst.want, got, opt) {
				t.Error(cmp.Diff(tst.want, got, opt))
			}
			got, err = db.GetPermission(ctx, got.ID)
			if err != nil {
				t.Error(err)
			}

			rem, err := db.DeletePermission(ctx, tst.in)
			if err != nil {
				t.Error(err)
			}
			if !rem.Delete.Valid || rem.Delete.Time.IsZero() {
				t.Error("unexpected delete", rem.Delete)
			}
		})
	}
}

func TestListPermission(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("missing DATABASE_CONNECTION env")
	}
	opt := cmp.FilterValues(func(a, b time.Time) bool {
		return true
	}, cmp.Comparer(func(a, b time.Time) bool {
		return a.Sub(b) < 10*time.Second && b.Sub(a) < 10*time.Second
	}))
	testOrg, testCreator := uuid.New().String(), uuid.New().String()
	tests := []struct {
		name string
		in   []storage.Permission
		err  error
	}{
		{
			name: "Success",
			in: []storage.Permission{
				{
					ID:          "some-keto-id1",
					SvcPermID:   uuid.New().String(),
					Name:        "test1",
					OrgID:       testOrg,
					Description: "First test permission",
					CreateUID:   testCreator,
					DeleteUID:   sql.NullString{Valid: true, String: testCreator},
				},
				{
					ID:          "some-keto-id2",
					SvcPermID:   uuid.New().String(),
					Name:        "test2",
					OrgID:       testOrg,
					Description: "Second test permission",
					CreateUID:   testCreator,
					DeleteUID:   sql.NullString{Valid: true, String: testCreator},
				},
			},
		},
	}
	db, clean := NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)
	for _, tst := range tests {
		tst := tst
		t.Run(tst.name, func(t *testing.T) {
			ctx := context.Background()
			var wantPerms []*storage.Permission
			for _, p := range tst.in {
				perm, err := db.CreatePermission(ctx, p)
				if !cmp.Equal(tst.err, err) {
					t.Error(cmp.Diff(tst.err, err))
				}
				wantPerms = append(wantPerms, perm)
				t.Cleanup(func() { db.DeletePermission(ctx, *perm) })
			}

			gotPerms, err := db.ListPermission(ctx, testOrg)
			if err != nil {
				t.Error(err)
			}
			if len(gotPerms) != len(wantPerms) {
				t.Error(cmp.Diff(len(gotPerms), len(wantPerms), opt))
			}
		})
	}
}
