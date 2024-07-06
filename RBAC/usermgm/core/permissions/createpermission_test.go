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

func TestCreatePermission(t *testing.T) {
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

	k := keto.New(baseURL)

	s := New(db, k)
	val := challenge.New(db, k)
	logr, _ := test.NewNullLogger()
	log := logr.WithField("test", "core permission create")
	ctx := logging.WithLogger(context.Background(), log)

	orgs := []string{}
	orgNames := []string{"Test Org 1", "Test Org 2", "Test Org 3"}

	for _, name := range orgNames {
		orgID, err := db.CreateOrg(ctx, storage.Organization{
			OrgName: name,
			Active:  true,
		})
		if err != nil {
			t.Fatal(err)
		}
		orgs = append(orgs, orgID)
	}

	ctx = metautils.ExtractIncoming(ctx).Add(hydra.OrgIDKey, orgs[0]).ToIncoming(ctx)
	testU, err := db.CreateUser(ctx, storage.User{
		OrgID:      orgs[0],
		Username:   "testPermissionRole",
		Email:      "testPermissionRole@example.com",
		InviteCode: random.InvitationCode(20),
	}, storage.Credential{
		Password: "somepassword",
	})
	if err != nil {
		t.Error(err)
	}

	u, err := db.CreateUser(ctx, storage.User{
		OrgID:      orgs[0],
		Username:   "testPermissionRoleGrant",
		Email:      "testPermissionRoleGrant@example.com",
		InviteCode: random.InvitationCode(20),
	}, storage.Credential{
		Password: "somepassword",
	})
	if err != nil {
		t.Error(err)
	}
	ctx = metautils.ExtractIncoming(ctx).Set(hydra.ClientIDKey, u.ID).ToIncoming(ctx)
	res := core.ServiceResource{
		Name:        "test-assign-core",
		Description: "Assign test in core permissions.",
		Resource:    "RBCA:role",
		Actions:     []string{"get"},
	}
	svc, err := s.CreatePermission(ctx, core.ServicePermission{
		Service: core.Service{
			Name:        "TestCreatePM",
			Description: "Testing create permission",
		},
		Res: []core.ServiceResource{res},
	})
	if err != nil {
		t.Fatal(err)
	}

	grants := map[string]core.OrgPermission{}
	for _, org := range orgs {
		ctx = metautils.ExtractIncoming(ctx).Set(hydra.OrgIDKey, org).ToIncoming(ctx)
		grID, err := s.GrantService(ctx, core.Grant{
			RoleID:      org,
			GrantID:     svc.Service.ID,
			Environment: "create",
		})
		if err != nil {
			t.Fatal(err)
		}
		pms, err := s.ListPermission(ctx, core.ListPermissionFilter{OrgID: org})
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
		grants[org] = pm
	}

	for _, orgID := range orgs {
		ctx = metautils.ExtractIncoming(ctx).Set(hydra.OrgIDKey, orgID).ToIncoming(ctx)
		perms, err := s.ListPermission(ctx, core.ListPermissionFilter{
			OrgID: orgID,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(perms) == 0 {
			t.Fatal("wants at least 1 org, but got ", len(perms))
		}
		gr := grants[orgID]

		// Direct assignment of permission
		v := core.Validation{
			Environment: gr.Environment,
			Action:      "get",
			Resource:    "RBCA:role",
			ID:          testU.ID,
			OrgID:       orgID,
		}
		// if _, err := val.Validate(ctx, v); err != nil {
		// 	t.Error(err)
		// }

		// New user
		v.ID = u.ID
		if _, err := val.Validate(ctx, v); err == nil {
			t.Error("permission not granted but authorized")
		}

		// Empty role
		r := core.Role{
			OrgID:     orgID,
			Name:      "test-assign-role" + orgID,
			CreateUID: gr.CreateUID,
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
			GrantID: perms[0].ID,
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

		// Direct assignment of permission should remain
		// v.ID = testID
		// if _, err := val.Validate(ctx, v); err != nil {
		// 	t.Error(err)
		// }
	}
}
