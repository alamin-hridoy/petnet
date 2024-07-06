package org

import (
	"context"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

type OrgStore interface {
	CreateOrg(context.Context, storage.Organization) (string, error)
	ActivateOrg(ctx context.Context, id string) error
	GetOrgByID(context.Context, string) (*storage.Organization, error)
	UpdateOrgByID(context.Context, storage.Organization) (*storage.Organization, error)
	ListServicePublic(context.Context) ([]storage.DefaultService, error)
	ListOrgPermission(ctx context.Context, orgID string) ([]storage.OrgPermission, error)
}

type UserStore interface {
	GetUserByID(context.Context, string) (*storage.User, error)
}

type PermissionStore interface {
	CreateRole(context.Context, core.Role) (string, error)
	ListRole(context.Context, core.ListRoleFilter) ([]core.Role, error)
	GrantService(context.Context, core.Grant) (string, error)
	AssignService(context.Context, string, core.Grant) (string, error)
	RoleGrant(context.Context, core.Grant) (*core.Role, error)
}

type Svc struct {
	org OrgStore
	pm  PermissionStore
	usr UserStore
}

func New(org OrgStore, pm PermissionStore, usr UserStore) *Svc {
	return &Svc{org: org, pm: pm, usr: usr}
}
