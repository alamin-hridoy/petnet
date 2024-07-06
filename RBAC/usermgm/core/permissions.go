package core

import "time"

type OrgPermission struct {
	// ID is the permission identifier.
	ID    string
	OrgID string
	// SvcID is the Service that created the permission.
	SvcID string
	// ServiceName
	SvcName string
	// SvcPermID is the Service Permission ID
	SvcPermID string
	GrantID   string
	// CreateUID is the user that created the permission
	CreateUID string
	// DeleteUID is the deleted user ID
	DeleteUID string
	// Name is the name of the permission.
	// The name must be unique and immutable after creation, to avoid duplicate entries.
	Name string
	// Description gives a user understandable summary of the permission.
	Description string
	// Environment is the target OrgID, prod/sandbox runtime, etc.
	// Which define the operating scope of the permission within the entire data/operation space.
	Environment string
	// Permission should allow (grant) or deny (restrict) the action.
	Allow bool
	// Actions to be taken on the resource(s).
	// Can be data operations like read/write/delete or
	// more basic endpoint action like call/use.
	Action string
	// Resources should uniquely identify the data or endpoint that is being acted upon.
	// Can be multiple data objects, especially if an endpoint operates on them simultaneously.
	Resource string
	// Groups or Users assigned the permission.
	// In most cases, these should correspond to permission Groups.
	// Assignment to user(s) should be limited to Admin or other highly-restricted permissions.
	Groups []string
}

type ServicePermission struct {
	Service Service
	Res     []ServiceResource
}

type ServiceResource struct {
	ID          string
	Name        string
	Description string
	Resource    string
	Actions     []string
}

type Service struct {
	ID           string
	Name         string
	Description  string
	GrantDefault bool
}

type Role struct {
	ID    string
	OrgID string
	Name  string
	Desc  string
	// CreateUID is the user that created the role
	CreateUID string
	// DeleteUID is the user that deleted the role
	DeleteUID   string
	UpdatedUID  string
	Permissions []string
	Members     []string
	Created     time.Time
	Updated     time.Time
	Count       int
}

type Grant struct {
	RoleID      string
	GrantID     string
	Environment string
	Default     bool
}

type Validation struct {
	// Environment is the target OrgID, prod/sandbox runtime, etc.
	// Which define the operating scope of the permission within the entire data/operation space.
	Environment string
	// Action to be taken on the resource(s).
	Action string
	// Resources should identify the data or endpoint that is being acted upon.
	Resource string
	// ID is the user or service account requesting the action.
	ID string
	// OrgID is the organization that owns the resource being requested.
	OrgID string
}

type ListPermissionFilter struct {
	ID          []string
	OrgID       string
	Environment string
}

type ListRoleFilter struct {
	ID     []string
	OrgID  string
	SortBy string
	Name   string
	UserID string
	Limit  int32
	Offset int32
}

type ListUserRolesRequest struct {
	UserID []string
}
