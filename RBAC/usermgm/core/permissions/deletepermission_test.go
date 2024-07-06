package permissions

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus/hooks/test"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/core/challenge"
	"brank.as/rbac/usermgm/integrations/keto"
	"brank.as/rbac/usermgm/storage"
	"brank.as/rbac/usermgm/storage/postgres"
)

func TestDelete(t *testing.T) {
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
	log := logr.WithField("test", "core permission delete")
	ctx := logging.WithLogger(context.Background(), log)

	orgID, err := db.CreateOrg(ctx, storage.Organization{
		OrgName: "Core Delete",
		Active:  true,
	})
	if err != nil {
		t.Error(err)
	}
	testU, err := db.CreateUser(ctx, storage.User{
		OrgID:      orgID,
		Username:   "testPermissionDelete",
		Email:      "testPermissionDelete@example.com",
		InviteCode: random.InvitationCode(20),
	}, storage.Credential{
		Password: "somepassword",
	})
	if err != nil {
		t.Error(err)
	}
	uid := uuid.New().String()
	ctx = metautils.ExtractIncoming(ctx).Set(hydra.ClientIDKey, uid).ToIncoming(ctx)
	res := core.ServiceResource{
		Name:        "test-delete-core",
		Description: "Delete test in core permissions.",
		Resource:    "permissions",
		Actions:     []string{"get"},
	}
	svc, err := s.CreatePermission(ctx, core.ServicePermission{
		Service: core.Service{
			Name:        "TestDeletePM",
			Description: "Testing delete permission",
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

	v := core.Validation{
		Environment: p.Environment,
		Action:      "get",
		Resource:    "permissions",
		ID:          testU.ID,
		OrgID:       p.OrgID,
	}
	// if _, err := val.Validate(ctx, v); err != nil {
	// 	t.Error(err)
	// }

	grs, err := s.store.GetAssignedServices(ctx, orgID)
	if err != nil {
		t.Fatal(err)
	}
	var gr storage.ServiceAssignment
	for _, g := range grs {
		if g.GrantID == grID {
			gr = g
			break
		}
	}

	if err := s.RevokeService(ctx, core.Grant{
		RoleID:      orgID,
		GrantID:     svc.Service.ID,
		Environment: "delete",
	}); err != nil {
		t.Fatal(err)
	}

	want := gr
	want.RevokeUserID = sql.NullString{Valid: true, String: uid}
	want.Revoked.Valid = true

	if _, err := val.Validate(ctx, v); err == nil || status.Code(err) != codes.PermissionDenied {
		t.Error("expected permission removed")
	}

	pms, err = s.ListPermission(ctx, core.ListPermissionFilter{OrgID: p.OrgID})
	if err != nil {
		t.Error(err)
	}
	if len(pms) != 0 {
		t.Errorf("permission not deleted %+v", pms)
	}
}
