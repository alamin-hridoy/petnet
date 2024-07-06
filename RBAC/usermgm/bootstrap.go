package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/core/permissions"
	"brank.as/rbac/usermgm/core/svcacct"
	"brank.as/rbac/usermgm/storage"
	"brank.as/rbac/usermgm/storage/postgres"
)

type bs struct {
	complete chan struct{}
	init     bool
	config   *viper.Viper
	store    *postgres.Storage
	perm     *permissions.Svc
	sa       *svcacct.Svc

	uid, oid string
}

func (b *bs) Org() string {
	<-b.complete
	return b.oid
}

func Bootstrap(ctx context.Context, config *viper.Viper, log *logrus.Entry, store *postgres.Storage, perm *permissions.Svc, sa *svcacct.Svc) (*bs, error) {
	b := bs{
		complete: make(chan struct{}),
		config:   config,
		store:    store,
		perm:     perm,
		sa:       sa,
	}
	return &b, nil
}

func (b *bs) Init(ctx context.Context) error {
	log := logging.FromContext(ctx).WithField("method", "bootstrap")

	uid, orgID, err := b.bootstrapUser(ctx, b.config)
	if err != nil {
		return err
	}
	defer close(b.complete)
	b.uid, b.oid = uid, orgID
	perm := b.perm
	ctx = metautils.ExtractIncoming(ctx).Set(hydra.ClientIDKey, uid).
		Set(hydra.OrgIDKey, orgID).ToIncoming(ctx)
	bid, err := perm.CreatePermission(ctx, core.ServicePermission{
		Service: core.Service{
			Name:        "Bootstrap Permissions",
			Description: "Initialize and manage the system permissions and services.",
		},
		Res: []core.ServiceResource{
			{
				Name:        "Permission",
				Description: "Service Permission",
				Resource:    "RBAC:permission",
				Actions: []string{
					"create", "view", "delete", "grantPermission", "delegatePermission",
				},
			},
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("create RBAC:permission")
		return err
	}

	pid, err := perm.CreatePermission(ctx, core.ServicePermission{
		Service: core.Service{
			Name:        "Service/Product Management",
			Description: "Manage service availability to other organizations.",
		},
		Res: []core.ServiceResource{
			{
				Name:        "Service",
				Description: "Service Management",
				Resource:    "RBAC:service",
				Actions: []string{
					"create", "view", "publish", "delete", "grantPermission", "delegatePermission",
				},
			},
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("create RBAC:service")
		return err
	}

	pids, err := bootstrapPermissions(ctx, perm)
	if err != nil {
		logging.WithError(err, log).Error("create permissions")
		return err
	}
	for _, id := range append(pids, bid.Service.ID, pid.Service.ID) {
		if _, err := perm.GrantService(ctx, core.Grant{
			RoleID:  orgID,
			GrantID: id,
		}); err != nil {
			logging.WithError(err, log).Error("create RBAC:service")
			return err
		}
		svc, err := b.store.GetService(ctx, id)
		if err != nil {
			logging.WithError(err, log).Error("getting service")
			return err
		}
		if svc.Default {
			if err := perm.PublicService(ctx, core.Grant{
				RoleID:  uid,
				GrantID: id,
				Default: true,
			}); err != nil {
				logging.WithError(err, log).WithField("service", svc.Name).Error("publishing")
				return err
			}
		}
	}
	b.setDefaults(logging.WithLogger(ctx, log), b.config, uid, orgID)

	role := ""
	if b.init {
		role, err = perm.CreateRole(ctx, core.Role{
			OrgID:     orgID,
			Name:      "Bootstrap",
			Desc:      "Bootstrap",
			CreateUID: uid,
			Members:   []string{uid},
		})
		if err != nil {
			return err
		}

		id, s, err := b.sa.CreateSvcAccount(ctx, storage.SvcAccount{
			AuthType:     storage.OAuth2,
			OrgID:        orgID,
			ClientName:   "BootstrapClient",
			CreateUserID: uid,
		})
		if err != nil {
			return err
		}
		if _, err := perm.AssignRole(ctx, core.Grant{RoleID: role, GrantID: id}); err != nil {
			return err
		}
		if _, err := perm.AssignRole(ctx, core.Grant{RoleID: role, GrantID: uid}); err != nil {
			return err
		}
		log.WithFields(logrus.Fields{
			"clientID":     id,
			"clientSecret": s,
		}).Println("bootstrap account created")
	} else {
		r, err := perm.ListRole(ctx, core.ListRoleFilter{OrgID: orgID})
		if err != nil {
			return err
		}
		if len(r) == 0 {
			return fmt.Errorf("missing bootstrap role")
		}
		sort.Slice(r, func(i, j int) bool { return r[i].Created.Before(r[j].Created) })
		if r[0].Name != "Bootstrap" {
			return fmt.Errorf("invalid bootstrap role %v", r[0])
		}
		role = r[0].ID
	}

	pms, err := perm.ListPermission(ctx, core.ListPermissionFilter{OrgID: orgID})
	if err != nil {
		return err
	}
	for _, p := range pms {
		if _, err := perm.RoleGrant(ctx, core.Grant{
			RoleID: role, GrantID: p.ID,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (b *bs) bootstrapUser(ctx context.Context, config *viper.Viper) (user, org string, err error) {
	log := logging.FromContext(ctx).WithField("method", "bootstrap.bootstrapuser")

	o, err := b.store.GetOrgs(ctx)
	if err != nil {
		return "", "", err
	}
	if len(o) > 0 {
		log.WithField("orgs", len(o)).Info("bootstrap previously completed")
		sort.Slice(o, func(i, j int) bool { return o[i].Created.Before(o[j].Created) })
		u, err := b.store.GetUsersByOrg(ctx, o[0].ID)
		if err != nil {
			return "", "", err
		}
		sort.Slice(u, func(i, j int) bool { return u[i].Created.Before(u[j].Created) })
		return u[0].ID, o[0].ID, nil
	}
	b.init = true

	orgID, err := b.store.CreateOrg(ctx, storage.Organization{
		OrgName:      config.GetString("bootstrap.name"),
		ContactEmail: config.GetString("bootstrap.email"),
		Active:       true,
	})
	if err != nil {
		logging.WithError(err, log).Error("create org")
		return "", "", err
	}
	ctx = metautils.ExtractIncoming(ctx).Set(hydra.OrgIDKey, orgID).ToIncoming(ctx)
	u, err := b.store.CreateUser(ctx, storage.User{
		OrgID:         orgID,
		Username:      config.GetString("bootstrap.email"),
		FirstName:     "bootstrap",
		LastName:      "admin",
		Email:         config.GetString("bootstrap.email"),
		EmailVerified: false,
		InviteStatus:  storage.Approved,
	}, storage.Credential{
		Password: "bootstrap",
	})
	if err != nil {
		logging.WithError(err, log).Error("create user")
		return "", "", err
	}
	return u.ID, orgID, nil
}

func bootstrapPermissions(ctx context.Context, perm *permissions.Svc) ([]string, error) {
	log := logging.FromContext(ctx).WithField("method", "bootstrap.permissions")

	ids := []string{}
	rid, err := perm.CreatePermission(ctx, core.ServicePermission{
		Service: core.Service{
			Name:         "RBAC Roles",
			Description:  "Manage the system permissions and services.",
			GrantDefault: true,
		},
		Res: []core.ServiceResource{
			{
				Name:        "Role",
				Description: "Role Management",
				Resource:    "RBAC:role",
				Actions: []string{
					"create", "view", "assign", "delete",
					"grantPermission", "delegatePermission",
				},
			},
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("create RBAC:role")
		return nil, err
	}
	ids = append(ids, rid.Service.ID)

	pid, err := perm.CreatePermission(ctx, core.ServicePermission{
		Service: core.Service{
			Name:         "RBAC Permissions",
			Description:  "View and assign system permissions.",
			GrantDefault: true,
		},
		Res: []core.ServiceResource{
			{
				Name:        "Permission",
				Description: "Permisson Assignment",
				Resource:    "RBAC:permission",
				Actions:     []string{"create", "view", "assign", "grantPermission", "delegatePermission"},
			},
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("create RBAC:permission")
		return nil, err
	}
	ids = append(ids, pid.Service.ID)

	auid, err := perm.CreatePermission(ctx, core.ServicePermission{
		Service: core.Service{
			Name:         "RBAC Accounts",
			Description:  "Manage organization user accounts",
			GrantDefault: true,
		},
		Res: []core.ServiceResource{
			{
				Name:        "User Account",
				Description: "User account management",
				Resource:    "ACCOUNT:user",
				Actions: []string{
					"create", "invite", "view", "grantPermission", "delegatePermission", "delete", "update",
				},
			},
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("create ACCOUNT:user")
		return nil, err
	}
	ids = append(ids, auid.Service.ID)

	asid, err := perm.CreatePermission(ctx, core.ServicePermission{
		Service: core.Service{
			Name:         "RBAC Service Accounts",
			Description:  "Manage organization service accounts",
			GrantDefault: true,
		},
		Res: []core.ServiceResource{
			{
				Name:        "Service Account",
				Description: "Service account management",
				Resource:    "ACCOUNT:service",
				Actions:     []string{"create", "view", "grantPermission", "delegatePermission"},
			},
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("create ACCOUNT:service")
		return nil, err
	}
	ids = append(ids, asid.Service.ID)

	aoid, err := perm.CreatePermission(ctx, core.ServicePermission{
		Service: core.Service{
			Name:         "RBAC Organization",
			Description:  "Manage organization",
			GrantDefault: true,
		},
		Res: []core.ServiceResource{
			{
				Name:        "Organization",
				Description: "Organization management",
				Resource:    "ACCOUNT:org",
				Actions:     []string{"create", "view", "grantPermission", "delegatePermission"},
			},
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("create ACCOUNT:org")
		return nil, err
	}
	ids = append(ids, aoid.Service.ID)

	mfaid, err := perm.CreatePermission(ctx, core.ServicePermission{
		Service: core.Service{
			Name:         "MFA Management",
			Description:  "Manage enabled MFA",
			GrantDefault: true,
		},
		Res: []core.ServiceResource{
			{
				Name:        "MFA",
				Description: "Multi-Factor Authentication management",
				Resource:    "ACCOUNT:mfa",
				Actions:     []string{"create", "view", "validate", "grantPermission", "delegatePermission"},
			},
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("create ACCOUNT:org")
		return nil, err
	}
	ids = append(ids, mfaid.Service.ID)
	return ids, nil
}

func (b *bs) setDefaults(ctx context.Context, conf *viper.Viper, usr, org string) error {
	log := logging.FromContext(ctx).WithField("method", "bootstrap.setdefaults")

	svcs, err := b.store.ListService(ctx)
	if err != nil {
		return err
	}
	for _, s := range svcs {
		if !strings.HasPrefix(s.Name, "RBAC") {
			fmt.Println("skipping", s.Name)
			continue
		}
		if err := b.perm.PublicService(ctx, core.Grant{
			RoleID:      usr,
			GrantID:     s.ID,
			Environment: "", // System environment
			Default:     true,
		}); err != nil {
			logging.WithError(err, log).Error("publish service")
		}
	}
	return nil
}
