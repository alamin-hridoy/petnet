package profile

import (
	"context"
	"fmt"

	pm "brank.as/petnet/svcutil/permission"
	apb "brank.as/rbac/gunk/v1/authenticate"
	authpb "brank.as/rbac/gunk/v1/authenticate"
	ppb "brank.as/rbac/gunk/v1/permissions"
	upb "brank.as/rbac/gunk/v1/user"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

type Mock struct {
	RevokedID string
	Roles     []*ppb.Role
}

func (m Mock) ValidatePermission(ctx context.Context, in *ppb.ValidatePermissionRequest, opts ...grpc.CallOption) (*ppb.ValidatePermissionResponse, error) {
	return &ppb.ValidatePermissionResponse{}, nil
}

func (m Mock) GrantService(ctx context.Context, in *ppb.GrantServiceRequest, opts ...grpc.CallOption) (*ppb.GrantServiceResponse, error) {
	return &ppb.GrantServiceResponse{}, nil
}

func (m Mock) ListServiceAssignments(ctx context.Context, in *ppb.ListServiceAssignmentsRequest, opts ...grpc.CallOption) (*ppb.ListServiceAssignmentsResponse, error) {
	return &ppb.ListServiceAssignmentsResponse{}, nil
}

func (m Mock) RevokeService(ctx context.Context, in *ppb.RevokeServiceRequest, opts ...grpc.CallOption) (*ppb.RevokeServiceResponse, error) {
	return &ppb.RevokeServiceResponse{}, nil
}

func (m Mock) PublicService(ctx context.Context, in *ppb.PublicServiceRequest, opts ...grpc.CallOption) (*ppb.PublicServiceResponse, error) {
	return &ppb.PublicServiceResponse{}, nil
}

func (m Mock) ListServices(ctx context.Context, in *ppb.ListServicesRequest, opts ...grpc.CallOption) (*ppb.ListServicesResponse, error) {
	return &ppb.ListServicesResponse{}, nil
}

func (m Mock) CreatePermission(ctx context.Context, in *ppb.CreatePermissionRequest, opts ...grpc.CallOption) (*ppb.CreatePermissionResponse, error) {
	return &ppb.CreatePermissionResponse{}, nil
}

func (m Mock) RetryMFA(ctx context.Context, in *apb.RetryMFARequest, opts ...grpc.CallOption) (*apb.Session, error) {
	return &apb.Session{}, nil
}

func (m Mock) ListPermission(ctx context.Context, in *ppb.ListPermissionRequest, opts ...grpc.CallOption) (*ppb.ListPermissionResponse, error) {
	ps1 := &ppb.Permission{
		ID:          "10000000-0000-0000-0000-000000000000",
		ServiceName: pm.AdminServiceName,
		Name:        "dsaListAndDetails",
		Resource:    "dsaListDetail",
		Action:      "create",
	}
	ps2 := &ppb.Permission{
		ID:          "20000000-0000-0000-0000-000000000000",
		ServiceName: pm.AdminServiceName,
		Name:        "dsaListAndDetails",
		Resource:    "dsaListDetail",
		Action:      "read",
	}
	ps3 := &ppb.Permission{
		ID:       "30000000-0000-0000-0000-000000000000",
		Name:     "RBACRequired",
		Resource: "ACCOUNT:user",
		Action:   "view",
	}

	ps := []*ppb.Permission{}
	switch len(in.GetID()) {
	case 0:
		ps = append(ps, ps1, ps2, ps3)
	case 1:
		ps = append(ps, ps1)
	case 2:
		ps = append(ps, ps1, ps2)
	}
	return &ppb.ListPermissionResponse{Permissions: ps}, nil
}

func (m Mock) DeletePermission(ctx context.Context, in *ppb.DeletePermissionRequest, opts ...grpc.CallOption) (*ppb.DeletePermissionResponse, error) {
	return &ppb.DeletePermissionResponse{}, nil
}

func (m Mock) AssignPermission(ctx context.Context, in *ppb.AssignPermissionRequest, opts ...grpc.CallOption) (*ppb.AssignPermissionResponse, error) {
	return &ppb.AssignPermissionResponse{}, nil
}

func (m Mock) RevokePermission(ctx context.Context, in *ppb.RevokePermissionRequest, opts ...grpc.CallOption) (*ppb.RevokePermissionResponse, error) {
	return &ppb.RevokePermissionResponse{}, nil
}

func (m Mock) GetUser(ctx context.Context, in *upb.GetUserRequest, opts ...grpc.CallOption) (*upb.GetUserResponse, error) {
	usc := make([]*upb.User, len(us))
	copy(usc, us)
	return &upb.GetUserResponse{User: usc[0]}, nil
}

func (m Mock) ListUsers(ctx context.Context, in *upb.ListUsersRequest, opts ...grpc.CallOption) (*upb.ListUsersResponse, error) {
	usc := make([]*upb.User, len(us))
	copy(usc, us)
	return &upb.ListUsersResponse{Users: usc}, nil
}

func (m Mock) ChangePassword(ctx context.Context, in *upb.ChangePasswordRequest, opts ...grpc.CallOption) (*upb.ChangePasswordResponse, error) {
	return &upb.ChangePasswordResponse{}, nil
}

func (m Mock) ConfirmUpdate(ctx context.Context, in *upb.ConfirmUpdateRequest, opts ...grpc.CallOption) (*upb.ConfirmUpdateResponse, error) {
	return &upb.ConfirmUpdateResponse{}, nil
}

func (s Mock) UpdateUser(ctx context.Context, in *upb.UpdateUserRequest, opts ...grpc.CallOption) (*upb.UpdateUserResponse, error) {
	return &upb.UpdateUserResponse{}, nil
}

func (m Mock) AuthenticateUser(ctx context.Context, in *upb.AuthenticateUserRequest, opts ...grpc.CallOption) (*upb.AuthenticateUserResponse, error) {
	return &upb.AuthenticateUserResponse{
		UserID: "10100000-0000-0000-0000-000000000000",
		OrgID:  "10200000-0000-0000-0000-000000000000",
	}, nil
}

func (m Mock) Signup(ctx context.Context, in *upb.SignupRequest, opts ...grpc.CallOption) (*upb.SignupResponse, error) {
	return &upb.SignupResponse{}, nil
}

func (m Mock) ResendConfirmEmail(ctx context.Context, in *upb.ResendConfirmEmailRequest, opts ...grpc.CallOption) (*upb.ResendConfirmEmailResponse, error) {
	return &upb.ResendConfirmEmailResponse{}, nil
}

func (m Mock) EmailConfirmation(ctx context.Context, in *upb.EmailConfirmationRequest, opts ...grpc.CallOption) (*upb.EmailConfirmationResponse, error) {
	return &upb.EmailConfirmationResponse{}, nil
}

func (m Mock) ForgotPassword(ctx context.Context, in *upb.ForgotPasswordRequest, opts ...grpc.CallOption) (*upb.ForgotPasswordResponse, error) {
	return &upb.ForgotPasswordResponse{}, nil
}

func (m Mock) ResetPassword(ctx context.Context, in *upb.ResetPasswordRequest, opts ...grpc.CallOption) (*upb.ResetPasswordResponse, error) {
	return &upb.ResetPasswordResponse{}, nil
}

func (m Mock) DisableUser(ctx context.Context, in *upb.DisableUserRequest, opts ...grpc.CallOption) (*upb.DisableUserResponse, error) {
	return &upb.DisableUserResponse{}, nil
}

func (m Mock) EnableUser(ctx context.Context, in *upb.EnableUserRequest, opts ...grpc.CallOption) (*upb.EnableUserResponse, error) {
	return &upb.EnableUserResponse{}, nil
}

func (m Mock) GetSession(ctx context.Context, req *authpb.GetSessionRequest, opts ...grpc.CallOption) (*authpb.Session, error) {
	return &authpb.Session{}, nil
}

func (m Mock) Login(ctx context.Context, req *authpb.LoginRequest, opts ...grpc.CallOption) (*authpb.Session, error) {
	return &authpb.Session{}, nil
}

func (m Mock) CreateRole(ctx context.Context, in *ppb.CreateRoleRequest, opts ...grpc.CallOption) (*ppb.CreateRoleResponse, error) {
	return &ppb.CreateRoleResponse{ID: uuid.New().String()}, nil
}

func (m Mock) ListRole(ctx context.Context, in *ppb.ListRoleRequest, opts ...grpc.CallOption) (*ppb.ListRoleResponse, error) {
	return &ppb.ListRoleResponse{
		Roles: m.Roles,
	}, nil
}

func (m Mock) UpdateRole(ctx context.Context, in *ppb.UpdateRoleRequest, opts ...grpc.CallOption) (*ppb.UpdateRoleResponse, error) {
	return &ppb.UpdateRoleResponse{}, nil
}

func (m Mock) DeleteRole(ctx context.Context, in *ppb.DeleteRoleRequest, opts ...grpc.CallOption) (*ppb.DeleteRoleResponse, error) {
	return &ppb.DeleteRoleResponse{}, nil
}

func (m Mock) AssignRolePermission(ctx context.Context, in *ppb.AssignRolePermissionRequest, opts ...grpc.CallOption) (*ppb.AssignRolePermissionResponse, error) {
	return &ppb.AssignRolePermissionResponse{}, nil
}

func (m Mock) RevokeRolePermission(ctx context.Context, in *ppb.RevokeRolePermissionRequest, opts ...grpc.CallOption) (*ppb.RevokeRolePermissionResponse, error) {
	pid := in.GetPermission()
	if m.RevokedID != "" {
		if pid != m.RevokedID {
			return nil, fmt.Errorf("revoke id mismatch, got: %v, want: %v", pid, m.RevokedID)
		}
	}
	return &ppb.RevokeRolePermissionResponse{}, nil
}

func (m Mock) AddUser(ctx context.Context, in *ppb.AddUserRequest, opts ...grpc.CallOption) (*ppb.AddUserResponse, error) {
	return &ppb.AddUserResponse{}, nil
}

func (m Mock) RemoveUser(ctx context.Context, in *ppb.RemoveUserRequest, opts ...grpc.CallOption) (*ppb.RemoveUserResponse, error) {
	return &ppb.RemoveUserResponse{}, nil
}

var us = []*upb.User{
	{
		ID:           "10000000-0000-0000-0000-000000000000",
		OrgID:        "20000000-0000-0000-0000-000000000000",
		OrgName:      "org-name",
		FirstName:    "firstname",
		LastName:     "lastname",
		Email:        "user@mail.com",
		InviteStatus: "Onboarding",
		CountryCode:  "32",
		Phone:        "123456789",
		Created:      tspb.Now(),
		Updated:      tspb.Now(),
	},
	{
		ID:           "30000000-0000-0000-0000-000000000000",
		OrgID:        "40000000-0000-0000-0000-000000000000",
		OrgName:      "org-name-two",
		FirstName:    "firstnametwo",
		LastName:     "lastnametwo",
		Email:        "user2@mail.com",
		InviteStatus: "Onboarding",
		CountryCode:  "32",
		Phone:        "123456782",
		Created:      tspb.Now(),
		Updated:      tspb.Now(),
	},
}
