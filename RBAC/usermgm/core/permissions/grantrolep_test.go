package permissions

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus/hooks/test"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/core/challenge"
	"brank.as/rbac/usermgm/integrations/keto"
	"brank.as/rbac/usermgm/storage"
	"brank.as/rbac/usermgm/storage/postgres"
)

func TestGrant(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	baseURL := os.Getenv("KETO_URL")
	switch "" {
	case baseURL:
		t.Skip("missing env 'KETO_URL'")
	case conn:
		t.Error("missing env 'DATABASE_CONNECTION'")
	}
	db, clean := postgres.NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)

	// opt := cmp.FilterPath(func(p cmp.Path) bool { return p.Last().String() == ".ID" }, cmp.Ignore())
	// tmOpt := cmp.FilterValues(func(a, b time.Time) bool { return true }, cmp.Ignore())

	k := keto.New(baseURL)

	s := New(db, k)
	val := challenge.New(db, k)
	logr, _ := test.NewNullLogger()
	log := logr.WithField("test", "core permission delete")
	ctx := logging.WithLogger(context.Background(), log)
	orgID, err := db.CreateOrg(ctx, storage.Organization{
		OrgName: "Core Grant",
		Active:  true,
	})
	ctx = metautils.ExtractIncoming(ctx).Add(hydra.OrgIDKey, orgID).ToIncoming(ctx)
	u, err := db.CreateUser(ctx, storage.User{
		OrgID:      orgID,
		Username:   "testPermissionRoleGrant",
		Email:      "testPermissionRoleGrant@example.com",
		InviteCode: random.InvitationCode(20),
	}, storage.Credential{
		Password: "somepassword",
	})
	if err != nil {
		t.Error(err)
	}

	ctx = metautils.ExtractIncoming(ctx).Add(hydra.ClientIDKey, u.ID).ToIncoming(ctx)
	res := core.ServiceResource{
		Name:        "test-assign-core",
		Description: "Assign test in core permissions.",
		Resource:    "permissions",
		Actions:     []string{"get"},
	}
	svc, err := s.CreatePermission(ctx, core.ServicePermission{
		Service: core.Service{
			Name:        "TestAssignPM",
			Description: "Testing assign permission",
		},
		Res: []core.ServiceResource{res},
	})
	if err != nil {
		t.Fatal(err)
	}

	grID, err := s.GrantService(ctx, core.Grant{
		RoleID:      orgID,
		GrantID:     svc.Service.ID,
		Environment: "delete",
	})
	if err != nil {
		t.Fatal(err)
	}

	pms, err := s.ListPermission(ctx, core.ListPermissionFilter{OrgID: orgID})
	if err != nil {
		t.Fatal(err)
	}
	var p core.OrgPermission
	for _, pp := range pms {
		if pp.GrantID == grID {
			p = pp
			break
		}
	}
	if !cmp.Equal(res.Resource, p.Resource) {
		t.Error(cmp.Diff(res.Resource, p.Resource))
	}
	if !cmp.Equal(res.Name, p.Name) {
		t.Error(cmp.Diff(res.Name, p.Name))
	}
	t.Cleanup(func() { s.DeleteOrgPermission(ctx, p) })

	// New user
	v := core.Validation{
		Environment: p.Environment,
		Action:      "get",
		Resource:    "permissions",
		ID:          u.ID,
		OrgID:       p.OrgID,
	}
	if _, err := val.Validate(ctx, v); err == nil {
		t.Error("permission not granted but authorized")
	}

	// Empty role
	r := core.Role{
		OrgID:     p.OrgID,
		Name:      "test-assign-role",
		CreateUID: p.CreateUID,
	}
	rID, err := s.CreateRole(ctx, r)
	if err != nil {
		t.Error(err)
	}
	r.ID = rID

	if _, err := s.AssignRole(ctx, core.Grant{RoleID: r.ID, GrantID: u.ID}); err != nil {
		t.Error(err)
	}

	g := core.Grant{
		RoleID:  r.ID,
		GrantID: p.ID,
	}
	// Grant Permission to role
	if _, err := s.RoleGrant(ctx, g); err != nil {
		t.Error(err)
	}

	if _, err := val.Validate(ctx, v); err != nil {
		t.Error(err)
	}

	// Revoke Permission from role
	if _, err := s.RoleRevoke(ctx, g); err != nil {
		t.Error(err)
	}

	if _, err := val.Validate(ctx, v); err == nil {
		t.Error("permission revoked but authorized")
	}
}
