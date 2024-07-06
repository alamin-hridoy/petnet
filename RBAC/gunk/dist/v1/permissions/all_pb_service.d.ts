// package: permissions
// file: brank.as/rbac/gunk/v1/permissions/all.proto

import * as brank_as_rbac_gunk_v1_permissions_all_pb from "./all_pb";
import {grpc} from "@improbable-eng/grpc-web";

type PermissionServiceCreatePermission = {
  readonly methodName: string;
  readonly service: typeof PermissionService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.CreatePermissionRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.CreatePermissionResponse;
};

type PermissionServiceListPermission = {
  readonly methodName: string;
  readonly service: typeof PermissionService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.ListPermissionRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.ListPermissionResponse;
};

type PermissionServiceDeletePermission = {
  readonly methodName: string;
  readonly service: typeof PermissionService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.DeletePermissionRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.DeletePermissionResponse;
};

type PermissionServiceAssignPermission = {
  readonly methodName: string;
  readonly service: typeof PermissionService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.AssignPermissionRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.AssignPermissionResponse;
};

type PermissionServiceRevokePermission = {
  readonly methodName: string;
  readonly service: typeof PermissionService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.RevokePermissionRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.RevokePermissionResponse;
};

export class PermissionService {
  static readonly serviceName: string;
  static readonly CreatePermission: PermissionServiceCreatePermission;
  static readonly ListPermission: PermissionServiceListPermission;
  static readonly DeletePermission: PermissionServiceDeletePermission;
  static readonly AssignPermission: PermissionServiceAssignPermission;
  static readonly RevokePermission: PermissionServiceRevokePermission;
}

type RoleServiceCreateRole = {
  readonly methodName: string;
  readonly service: typeof RoleService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.CreateRoleRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.CreateRoleResponse;
};

type RoleServiceListRole = {
  readonly methodName: string;
  readonly service: typeof RoleService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.ListRoleRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.ListRoleResponse;
};

type RoleServiceListUserRoles = {
  readonly methodName: string;
  readonly service: typeof RoleService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.ListUserRolesRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.ListUserRolesResponse;
};

type RoleServiceUpdateRole = {
  readonly methodName: string;
  readonly service: typeof RoleService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.UpdateRoleRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.UpdateRoleResponse;
};

type RoleServiceDeleteRole = {
  readonly methodName: string;
  readonly service: typeof RoleService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.DeleteRoleRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.DeleteRoleResponse;
};

type RoleServiceAssignRolePermission = {
  readonly methodName: string;
  readonly service: typeof RoleService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.AssignRolePermissionRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.AssignRolePermissionResponse;
};

type RoleServiceRevokeRolePermission = {
  readonly methodName: string;
  readonly service: typeof RoleService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.RevokeRolePermissionRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.RevokeRolePermissionResponse;
};

type RoleServiceAddUser = {
  readonly methodName: string;
  readonly service: typeof RoleService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.AddUserRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.AddUserResponse;
};

type RoleServiceRemoveUser = {
  readonly methodName: string;
  readonly service: typeof RoleService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.RemoveUserRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.RemoveUserResponse;
};

export class RoleService {
  static readonly serviceName: string;
  static readonly CreateRole: RoleServiceCreateRole;
  static readonly ListRole: RoleServiceListRole;
  static readonly ListUserRoles: RoleServiceListUserRoles;
  static readonly UpdateRole: RoleServiceUpdateRole;
  static readonly DeleteRole: RoleServiceDeleteRole;
  static readonly AssignRolePermission: RoleServiceAssignRolePermission;
  static readonly RevokeRolePermission: RoleServiceRevokeRolePermission;
  static readonly AddUser: RoleServiceAddUser;
  static readonly RemoveUser: RoleServiceRemoveUser;
}

type ProductServiceGrantService = {
  readonly methodName: string;
  readonly service: typeof ProductService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.GrantServiceRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.GrantServiceResponse;
};

type ProductServiceListServiceAssignments = {
  readonly methodName: string;
  readonly service: typeof ProductService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.ListServiceAssignmentsRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.ListServiceAssignmentsResponse;
};

type ProductServiceRevokeService = {
  readonly methodName: string;
  readonly service: typeof ProductService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.RevokeServiceRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.RevokeServiceResponse;
};

type ProductServicePublicService = {
  readonly methodName: string;
  readonly service: typeof ProductService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.PublicServiceRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.PublicServiceResponse;
};

type ProductServiceListServices = {
  readonly methodName: string;
  readonly service: typeof ProductService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.ListServicesRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.ListServicesResponse;
};

export class ProductService {
  static readonly serviceName: string;
  static readonly GrantService: ProductServiceGrantService;
  static readonly ListServiceAssignments: ProductServiceListServiceAssignments;
  static readonly RevokeService: ProductServiceRevokeService;
  static readonly PublicService: ProductServicePublicService;
  static readonly ListServices: ProductServiceListServices;
}

type ValidationServiceValidatePermission = {
  readonly methodName: string;
  readonly service: typeof ValidationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.ValidatePermissionRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_permissions_all_pb.ValidatePermissionResponse;
};

export class ValidationService {
  static readonly serviceName: string;
  static readonly ValidatePermission: ValidationServiceValidatePermission;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class PermissionServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  createPermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.CreatePermissionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.CreatePermissionResponse|null) => void
  ): UnaryResponse;
  createPermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.CreatePermissionRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.CreatePermissionResponse|null) => void
  ): UnaryResponse;
  listPermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListPermissionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListPermissionResponse|null) => void
  ): UnaryResponse;
  listPermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListPermissionRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListPermissionResponse|null) => void
  ): UnaryResponse;
  deletePermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.DeletePermissionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.DeletePermissionResponse|null) => void
  ): UnaryResponse;
  deletePermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.DeletePermissionRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.DeletePermissionResponse|null) => void
  ): UnaryResponse;
  assignPermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.AssignPermissionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.AssignPermissionResponse|null) => void
  ): UnaryResponse;
  assignPermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.AssignPermissionRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.AssignPermissionResponse|null) => void
  ): UnaryResponse;
  revokePermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RevokePermissionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RevokePermissionResponse|null) => void
  ): UnaryResponse;
  revokePermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RevokePermissionRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RevokePermissionResponse|null) => void
  ): UnaryResponse;
}

export class RoleServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  createRole(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.CreateRoleRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.CreateRoleResponse|null) => void
  ): UnaryResponse;
  createRole(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.CreateRoleRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.CreateRoleResponse|null) => void
  ): UnaryResponse;
  listRole(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListRoleRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListRoleResponse|null) => void
  ): UnaryResponse;
  listRole(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListRoleRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListRoleResponse|null) => void
  ): UnaryResponse;
  listUserRoles(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListUserRolesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListUserRolesResponse|null) => void
  ): UnaryResponse;
  listUserRoles(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListUserRolesRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListUserRolesResponse|null) => void
  ): UnaryResponse;
  updateRole(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.UpdateRoleRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.UpdateRoleResponse|null) => void
  ): UnaryResponse;
  updateRole(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.UpdateRoleRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.UpdateRoleResponse|null) => void
  ): UnaryResponse;
  deleteRole(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.DeleteRoleRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.DeleteRoleResponse|null) => void
  ): UnaryResponse;
  deleteRole(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.DeleteRoleRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.DeleteRoleResponse|null) => void
  ): UnaryResponse;
  assignRolePermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.AssignRolePermissionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.AssignRolePermissionResponse|null) => void
  ): UnaryResponse;
  assignRolePermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.AssignRolePermissionRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.AssignRolePermissionResponse|null) => void
  ): UnaryResponse;
  revokeRolePermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RevokeRolePermissionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RevokeRolePermissionResponse|null) => void
  ): UnaryResponse;
  revokeRolePermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RevokeRolePermissionRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RevokeRolePermissionResponse|null) => void
  ): UnaryResponse;
  addUser(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.AddUserRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.AddUserResponse|null) => void
  ): UnaryResponse;
  addUser(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.AddUserRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.AddUserResponse|null) => void
  ): UnaryResponse;
  removeUser(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RemoveUserRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RemoveUserResponse|null) => void
  ): UnaryResponse;
  removeUser(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RemoveUserRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RemoveUserResponse|null) => void
  ): UnaryResponse;
}

export class ProductServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  grantService(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.GrantServiceRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.GrantServiceResponse|null) => void
  ): UnaryResponse;
  grantService(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.GrantServiceRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.GrantServiceResponse|null) => void
  ): UnaryResponse;
  listServiceAssignments(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListServiceAssignmentsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListServiceAssignmentsResponse|null) => void
  ): UnaryResponse;
  listServiceAssignments(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListServiceAssignmentsRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListServiceAssignmentsResponse|null) => void
  ): UnaryResponse;
  revokeService(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RevokeServiceRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RevokeServiceResponse|null) => void
  ): UnaryResponse;
  revokeService(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RevokeServiceRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.RevokeServiceResponse|null) => void
  ): UnaryResponse;
  publicService(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.PublicServiceRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.PublicServiceResponse|null) => void
  ): UnaryResponse;
  publicService(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.PublicServiceRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.PublicServiceResponse|null) => void
  ): UnaryResponse;
  listServices(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListServicesRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListServicesResponse|null) => void
  ): UnaryResponse;
  listServices(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListServicesRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ListServicesResponse|null) => void
  ): UnaryResponse;
}

export class ValidationServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  validatePermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ValidatePermissionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ValidatePermissionResponse|null) => void
  ): UnaryResponse;
  validatePermission(
    requestMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ValidatePermissionRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_permissions_all_pb.ValidatePermissionResponse|null) => void
  ): UnaryResponse;
}

