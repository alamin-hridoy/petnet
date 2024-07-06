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

func TestAssign(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	baseURL := os.Getenv("KETO_URL")
	switch "" {
	case baseURL:
		t.Skip("missing env 'KETO_URL'")
	case conn:
		t.Error("missing env 'DATABASE_CONNECTION'")
	}

	if conn == "" {
		t.Skip("missing DATABASE_CONNECTION environment")
	}
	db, clean := postgres.NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)

	// opt := cmp.FilterPath(func(p cmp.Path) bool { return p.Last().String() == ".ID" }, cmp.Ignore())
	// tmOpt := cmp.FilterValues(func(a, b time.Time) bool { return true }, cmp.Ignore())

	k := keto.New(baseURL)
	val := challenge.New(db, k)

	s := New(db, k)
	logr, _ := test.NewNullLogger()
	log := logr.WithField("test", "core permission delete")
	ctx := logging.WithLogger(context.Background(), log)
	orgID, err := db.CreateOrg(ctx, storage.Organization{
		OrgName: "Core Assign",
		Active:  true,
	})
	if err != nil {
		t.Error(err)
	}
	ctx = metautils.ExtractIncoming(ctx).Add(hydra.OrgIDKey, orgID).ToIncoming(ctx)
	u, err := db.CreateUser(ctx, storage.User{
		OrgID:      orgID,
		Username:   "testPermissionAssignGrant",
		Email:      "testPermissionAssignGrant@example.com",
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
	var pm core.OrgPermission
	for _, p := range pms {
		if p.GrantID == grID {
			pm = p
			break
		}
	}
	if !cmp.Equal(res.Resource, pm.Resource) {
		t.Error(cmp.Diff(res.Resource, pm.Resource))
	}
	if !cmp.Equal(res.Name, pm.Name) {
		t.Error(cmp.Diff(res.Name, pm.Name))
	}
	t.Cleanup(func() { s.DeleteOrgPermission(ctx, pm) })

	// New user has not been granted
	v := core.Validation{
		Environment: pm.Environment,
		Action:      "get",
		Resource:    "permissions",
		ID:          u.ID,
		OrgID:       pm.OrgID,
	}
	if _, err := val.Validate(ctx, v); err == nil {
		t.Error("permission not granted but authorized")
	}

	// Empty role
	r := core.Role{
		OrgID:     pm.OrgID,
		Name:      "test-assign-role",
		CreateUID: pm.CreateUID,
	}
	rID, err := s.CreateRole(ctx, r)
	if err != nil {
		t.Error(err)
	}
	r.ID = rID

	// Grant Permission to role
	if _, err := s.RoleGrant(ctx, core.Grant{
		RoleID:  r.ID,
		GrantID: pm.ID,
	}); err != nil {
		t.Error(err)
	}

	if _, err := val.Validate(ctx, v); err == nil {
		t.Error("user not assigned but still authorized")
	}

	// Assign user to role
	if _, err := s.AssignRole(ctx, core.Grant{
		RoleID:  r.ID,
		GrantID: u.ID,
	}); err != nil {
		t.Error(err)
	}

	// Check permission was granted
	if _, err := val.Validate(ctx, v); err != nil {
		t.Error(err)
	}

	// Assign user to role
	if _, err := s.RemoveRole(ctx, core.Grant{
		RoleID:  r.ID,
		GrantID: u.ID,
	}); err != nil {
		t.Error(err)
	}

	if _, err := val.Validate(ctx, v); err == nil {
		t.Error("user not assigned but still authorized")
	}
}
