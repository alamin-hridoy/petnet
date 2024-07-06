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

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func TestListRole(t *testing.T) {
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
	rl := []storage.Role{
		{
			ID:          "some-keto-id1",
			OrgID:       testOrg,
			Name:        "test1",
			Description: "First test permission",
			CreateUID:   testCreator,
			UpdatedUID:  testCreator,
			DeleteUID:   sql.NullString{Valid: true, String: testCreator},
		},
		{
			ID:          "some-keto-id2",
			Name:        "test2",
			OrgID:       testOrg,
			Description: "Second test permission",
			CreateUID:   testCreator,
			UpdatedUID:  testCreator,
			DeleteUID:   sql.NullString{Valid: true, String: testCreator},
		},
		{
			ID:          "some-keto-id3",
			Name:        "test3",
			OrgID:       testOrg,
			Description: "third test permission",
			CreateUID:   testCreator,
			UpdatedUID:  testCreator,
			DeleteUID:   sql.NullString{Valid: true, String: testCreator},
		},
	}
	tests := []struct {
		name string
		in   []storage.Role
		err  error
		f    core.ListRoleFilter
		want []storage.Role
	}{
		{
			name: "Success",
			in:   rl,
			f: core.ListRoleFilter{
				OrgID: testOrg,
			},
			want: []storage.Role{rl[0], rl[1], rl[2]},
		},
		{
			name: "Limit",
			f: core.ListRoleFilter{
				OrgID:  testOrg,
				Limit:  2,
				Offset: 0,
			},
			want: []storage.Role{rl[2], rl[1]},
		},
		{
			name: "Offset",
			f: core.ListRoleFilter{
				OrgID:  testOrg,
				Limit:  2,
				Offset: 1,
			},
			want: []storage.Role{rl[1], rl[0]},
		},
		{
			name: "Limit+Offset",
			f: core.ListRoleFilter{
				OrgID:  testOrg,
				Limit:  1,
				Offset: 1,
			},
			want: []storage.Role{rl[1]},
		},
	}
	db, clean := NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)
	var wantPerms []*storage.Role
	for _, tst := range tests {
		tst := tst
		t.Run(tst.name, func(t *testing.T) {
			ctx := context.Background()
			for _, p := range tst.in {
				perm, err := db.CreateRole(ctx, p)
				if !cmp.Equal(tst.err, err) {
					t.Error(cmp.Diff(tst.err, err))
				}
				wantPerms = append(wantPerms, perm)
			}
		})
	}
	for _, tst := range tests {
		tst := tst
		t.Run(tst.name, func(t *testing.T) {
			ctx := context.Background()
			gotPerms, err := db.ListRole(ctx, tst.f)
			if err != nil {
				t.Error(err)
			}
			if len(gotPerms) != len(tst.want) {
				t.Log(tst.name)
				t.Error(cmp.Diff(gotPerms, tst.want, opt))
			}
		})
	}
	for _, wp := range wantPerms {
		wp := wp
		t.Cleanup(func() {
			ctx := context.Background()
			_, err := db.DeleteRole(ctx, *wp)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
