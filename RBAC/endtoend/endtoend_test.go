//go:build ignore
// +build ignore

package endtoend

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	ipb "brank.as/rbac/gunk/v1/invite"
	ppb "brank.as/rbac/gunk/v1/permissions"
	sapb "brank.as/rbac/gunk/v1/serviceaccount"
	upb "brank.as/rbac/gunk/v1/user"
)

func TestEndtoend(t *testing.T) {
	env := envDefault()
	if env != "staging" {
		t.Skip("Skipping e2e test")
	}
	if os.Getenv("RUN_E2E_TEST") == "" {
		t.Skip("Skipping e2e test, $RUN_E2E_TEST not set")
	}
	cnf := newConfig(t)
	ut := util{t: t, cnf: cnf}

	t.Log("Connecting to usermgm")
	u := cnf.GetString("usermgm.api")
	conn, ctx := ut.newUsermgmConn(u)
	defer conn.Close()

	t.Log("Signing up with new user")
	email := "e2e-" + uuid.New().String() + "@mail.com"
	pass := "secret0123456789"
	uid, oid := ut.signup(ctx, conn, email, pass)

	clid := cnf.GetString("endtoend.clientID")
	clsec := cnf.GetString("endtoend.clientSecret")

	t.Log("Getting token of endtoend user")
	tok := ut.getToken(clid, clsec)

	t.Log("Connecting to usermgm with endtoend user")
	conn, ctx = ut.newUsermgmConn(u, withToken(tok))
	defer conn.Close()

	t.Log("Approving new user")
	ut.approve(ctx, conn, uid)

	t.Log("Creating new permissions for new user")
	_ = ut.createPermission(ctx, conn, &ppb.CreatePermissionRequest{
		ServiceName: "",
		Description: "",
		Permissions: []*ppb.ServicePermission{
			{},
			{},
			{},
			{},
		},
	})

	t.Log("Logging in with new user")
	redirectURL := cnf.GetString("auth.url") + cnf.GetString("auth.redirectPath")
	cdpCtx, canc := ut.newCDP()
	code := ut.login(cdpCtx, redirectURL, email, pass, true)
	defer canc()

	t.Log("Exchanging code for token")
	tok = ut.exchangeCode(ctx, code)

	t.Log("Connecting to usermgm with new user")
	conn, ctx = ut.newUsermgmConn(u, withToken(tok))
	defer conn.Close()

	t.Log("Creating role and permission")
	rl := "e2e-invite-user-role"
	rid := ut.createRole(ctx, conn, rl)
	pname := "E2E Role Assign test"
	_ = ut.createPermission(ctx, conn, &ppb.CreatePermissionRequest{
		Name:        pname,
		Description: "E2E permission get user",
		Actions:     []string{"view", "create", "assign"},
		Resources:   []string{"ACCOUNT:user", "ACCOUNT:service", "RBAC:role"},
		Groups:      []string{rl},
	})

	// TODO(robin): when creating permissions with user it still returns the svcaccount permissions id
	t.Log("Getting created permissions id")
	var pid string
	ps := ut.listPermissions(ctx, conn, oid)
	for _, p := range ps {
		if p.Name == pname {
			pid = p.ID
		}
	}

	t.Log("Assigning role to permission")
	ut.assignRolePermission(ctx, conn, rid, pid)

	t.Log("Creating webhook token")
	token := ut.createWHToken(t)
	invEmail := token + "@email.webhook.site"

	t.Log("Inviting user")
	uid = ut.inviteUser(ctx, conn, rid, oid, invEmail)

	t.Log("Adding invited user to role")
	ut.addUserToRole(ctx, conn, rid, uid)

	t.Log("Getting invite URL")
	iu, err := ut.getInvURL(t, token, cnf.GetString("idp.url"))
	if err != nil {
		t.Fatal("didn't get invite url")
	}

	t.Log("Signing up with invited user")
	ut.inviteSignup(cdpCtx, t, iu, pass, true)

	t.Log("Logging in with invited user")
	code = ut.login(cdpCtx, redirectURL, invEmail, pass, true)

	t.Log("Exchanging code for token")
	tok = ut.exchangeCode(ctx, code)

	t.Log("Connecting to usermgm with invited user")
	conn, ctx = ut.newUsermgmConn(u, withToken(tok))
	defer conn.Close()

	t.Log("Creating service account")
	clid, clsec = ut.createSvcAccount(ctx, conn, rl)

	t.Log("Adding svc account to role")
	ut.addUserToRole(ctx, conn, rid, clid)

	t.Log("Getting token of service account")
	tok = ut.getToken(clid, clsec)

	t.Log("Connecting to usermgm with service account")
	conn, ctx = ut.newUsermgmConn(u, withToken(tok))
	defer conn.Close()

	t.Log("Getting user")
	ut.getUser(ctx, conn, uid)

	t.Log("Creating role and permission")
	rl2 := "e2e-delete-role-permission"
	rid2 := ut.createRole(ctx, conn, rl2)
	_ = ut.createPermission(ctx, conn, &ppb.CreatePermissionRequest{
		Name:        "E2E delete role permission test",
		Description: "delete role permission",
		Actions:     []string{"assign", "delete"},
		Resources:   []string{"RBAC:role", "RBAC:permission"},
		Groups:      []string{rl2},
	})

	// TODO(robin): when creating permissions with user it still returns the svcaccount permissions id
	t.Log("Getting created permissions id")
	var pid2 string
	ps = ut.listPermissions(ctx, conn, oid)
	for _, p := range ps {
		if p.Name == pname {
			pid2 = p.ID
		}
	}

	t.Log("Assigning role to permission")
	ut.assignRolePermission(ctx, conn, rid2, pid2)

	t.Log("Revoke permission from role")
	ut.revokeRolePermission(ctx, conn, rid2, pid2)

	t.Log("Add user to role")
	ut.addUserToRole(ctx, conn, rid2, uid)

	t.Log("Revoke user from role")
	ut.removeUserFromRole(ctx, conn, rid2, uid)

	t.Log("Delete role")
	ut.deleteRole(ctx, conn, rid2)

	t.Log("Delete permission")
	ut.removePermission(ctx, conn, pid2)
}

func (ut util) signup(ctx context.Context, conn *grpc.ClientConn, email, pass string) (string, string) {
	cl := upb.NewSignupClient(conn)
	res, err := cl.Signup(ctx, &upb.SignupRequest{
		Username:  email,
		FirstName: "first",
		LastName:  "last",
		Email:     email,
		Password:  pass,
	})
	if err != nil {
		ut.t.Fatal("signing up: ", err)
	}
	return res.UserID, res.OrgID
}

func (ut util) approve(ctx context.Context, conn *grpc.ClientConn, uid string) {
	cl := ipb.NewInviteServiceClient(conn)
	_, err := cl.Approve(ctx, &ipb.ApproveRequest{
		ID: uid,
	})
	if err != nil {
		ut.t.Fatal("approving user: ", err)
	}
}

func (ut util) createPermission(ctx context.Context, conn *grpc.ClientConn, p *ppb.CreatePermissionRequest) string {
	cl := ppb.NewPermissionServiceClient(conn)
	res, err := cl.CreatePermission(ctx, p)
	if err != nil {
		ut.t.Fatal("creating permission: ", err)
	}
	return res.ID
}

func (ut util) listPermissions(ctx context.Context, conn *grpc.ClientConn, oid string) []*ppb.Permission {
	cl := ppb.NewPermissionServiceClient(conn)
	res, err := cl.ListPermission(ctx, &ppb.ListPermissionRequest{
		OrgID: oid,
	})
	if err != nil {
		ut.t.Fatal("list permissions: ", err)
	}
	return res.Permissions
}

func (ut util) removePermission(ctx context.Context, conn *grpc.ClientConn, pid string) {
	cl := ppb.NewPermissionServiceClient(conn)
	_, err := cl.DeletePermission(ctx, &ppb.DeletePermissionRequest{
		ID: pid,
	})
	if err != nil {
		ut.t.Fatal("delete permission: ", err)
	}
}

func (ut util) createRole(ctx context.Context, conn *grpc.ClientConn, rl string) string {
	cl := ppb.NewRoleServiceClient(conn)
	res, err := cl.CreateRole(ctx, &ppb.CreateRoleRequest{
		Name:        rl,
		Description: "e2e-role",
	})
	if err != nil {
		ut.t.Fatal("creating role: ", err)
	}
	return res.ID
}

func (ut util) assignRolePermission(ctx context.Context, conn *grpc.ClientConn, rid, pid string) {
	cl := ppb.NewRoleServiceClient(conn)
	_, err := cl.AssignRolePermission(ctx, &ppb.AssignRolePermissionRequest{
		RoleID:     rid,
		Permission: pid,
	})
	if err != nil {
		ut.t.Fatal("assigning role permissions: ", err)
	}
}

func (ut util) revokeRolePermission(ctx context.Context, conn *grpc.ClientConn, rid, pid string) {
	cl := ppb.NewRoleServiceClient(conn)
	_, err := cl.RevokeRolePermission(ctx, &ppb.RevokeRolePermissionRequest{
		RoleID:     rid,
		Permission: pid,
	})
	if err != nil {
		ut.t.Fatal("revoke role permissions: ", err)
	}
}

func (ut util) addUserToRole(ctx context.Context, conn *grpc.ClientConn, rid, uid string) {
	cl := ppb.NewRoleServiceClient(conn)
	_, err := cl.AddUser(ctx, &ppb.AddUserRequest{
		RoleID: rid,
		UserID: uid,
	})
	if err != nil {
		ut.t.Fatal("adding user to role: ", err)
	}
}

func (ut util) removeUserFromRole(ctx context.Context, conn *grpc.ClientConn, rid, uid string) {
	cl := ppb.NewRoleServiceClient(conn)
	_, err := cl.RemoveUser(ctx, &ppb.RemoveUserRequest{
		RoleID: rid,
		UserID: uid,
	})
	if err != nil {
		ut.t.Fatal("remove user from role: ", err)
	}
}

func (ut util) deleteRole(ctx context.Context, conn *grpc.ClientConn, id string) {
	cl := ppb.NewRoleServiceClient(conn)
	res, err := cl.DeleteRole(ctx, &ppb.DeleteRoleRequest{
		ID: id,
	})
	if err != nil {
		ut.t.Fatal("delete role: ", err)
	}
	if res.Deleted == nil {
		ut.t.Fatal("delete role success with empty delete time")
	}
}

func (ut util) inviteUser(ctx context.Context, conn *grpc.ClientConn, rid, oid, email string) string {
	cl := ipb.NewInviteServiceClient(conn)
	res, err := cl.InviteUser(ctx, &ipb.InviteUserRequest{
		OrgID:     oid,
		OrgName:   "e2e-org",
		FirstName: "first",
		LastName:  "last",
		Email:     email,
		Phone:     "123456789",
		Role:      rid,
	})
	if err != nil {
		ut.t.Fatal("inviting user: ", err)
	}

	return res.ID
}

func (ut util) createSvcAccount(ctx context.Context, conn *grpc.ClientConn, rl string) (string, string) {
	cl := sapb.NewSvcAccountServiceClient(conn)
	res, err := cl.CreateAccount(ctx, &sapb.CreateAccountRequest{
		Name:     "E2E SvcAccount",
		Env:      "Sandbox",
		Role:     rl,
		AuthType: sapb.AuthType_OAuth2,
	})
	if err != nil {
		ut.t.Fatal("creating svc account: ", err)
	}
	return res.ClientID, res.Secret
}

func (ut util) getUser(ctx context.Context, conn *grpc.ClientConn, id string) {
	cl := upb.NewUserServiceClient(conn)
	if res, err := cl.GetUser(ctx, &upb.GetUserRequest{
		ID: id,
	}); err != nil || res == nil {
		ut.t.Fatal("getting user: ", err)
	}
}
