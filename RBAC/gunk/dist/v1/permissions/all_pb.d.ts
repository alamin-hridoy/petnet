// package: permissions
// file: brank.as/rbac/gunk/v1/permissions/all.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class CreatePermissionRequest extends jspb.Message {
  getServicename(): string;
  setServicename(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  clearPermissionsList(): void;
  getPermissionsList(): Array<ServicePermission>;
  setPermissionsList(value: Array<ServicePermission>): void;
  addPermissions(value?: ServicePermission, index?: number): ServicePermission;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePermissionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreatePermissionRequest): CreatePermissionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreatePermissionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreatePermissionRequest;
  static deserializeBinaryFromReader(message: CreatePermissionRequest, reader: jspb.BinaryReader): CreatePermissionRequest;
}

export namespace CreatePermissionRequest {
  export type AsObject = {
    servicename: string,
    description: string,
    permissionsList: Array<ServicePermission.AsObject>,
  }
}

export class ServicePermission extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  clearActionsList(): void;
  getActionsList(): Array<string>;
  setActionsList(value: Array<string>): void;
  addActions(value: string, index?: number): string;

  getResource(): string;
  setResource(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServicePermission.AsObject;
  static toObject(includeInstance: boolean, msg: ServicePermission): ServicePermission.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ServicePermission, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServicePermission;
  static deserializeBinaryFromReader(message: ServicePermission, reader: jspb.BinaryReader): ServicePermission;
}

export namespace ServicePermission {
  export type AsObject = {
    name: string,
    description: string,
    actionsList: Array<string>,
    resource: string,
  }
}

export class CreatePermissionResponse extends jspb.Message {
  getIdMap(): jspb.Map<string, string>;
  clearIdMap(): void;
  getServiceid(): string;
  setServiceid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePermissionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreatePermissionResponse): CreatePermissionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreatePermissionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreatePermissionResponse;
  static deserializeBinaryFromReader(message: CreatePermissionResponse, reader: jspb.BinaryReader): CreatePermissionResponse;
}

export namespace CreatePermissionResponse {
  export type AsObject = {
    idMap: Array<[string, string]>,
    serviceid: string,
  }
}

export class ListPermissionRequest extends jspb.Message {
  getOrgid(): string;
  setOrgid(value: string): void;

  getEnvironment(): string;
  setEnvironment(value: string): void;

  clearIdList(): void;
  getIdList(): Array<string>;
  setIdList(value: Array<string>): void;
  addId(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPermissionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListPermissionRequest): ListPermissionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListPermissionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPermissionRequest;
  static deserializeBinaryFromReader(message: ListPermissionRequest, reader: jspb.BinaryReader): ListPermissionRequest;
}

export namespace ListPermissionRequest {
  export type AsObject = {
    orgid: string,
    environment: string,
    idList: Array<string>,
  }
}

export class Permission extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getServicename(): string;
  setServicename(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getEnvironment(): string;
  setEnvironment(value: string): void;

  getRestrict(): boolean;
  setRestrict(value: boolean): void;

  getAction(): string;
  setAction(value: string): void;

  getResource(): string;
  setResource(value: string): void;

  clearGroupsList(): void;
  getGroupsList(): Array<string>;
  setGroupsList(value: Array<string>): void;
  addGroups(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Permission.AsObject;
  static toObject(includeInstance: boolean, msg: Permission): Permission.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Permission, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Permission;
  static deserializeBinaryFromReader(message: Permission, reader: jspb.BinaryReader): Permission;
}

export namespace Permission {
  export type AsObject = {
    id: string,
    servicename: string,
    name: string,
    description: string,
    environment: string,
    restrict: boolean,
    action: string,
    resource: string,
    groupsList: Array<string>,
  }
}

export class ListPermissionResponse extends jspb.Message {
  clearPermissionsList(): void;
  getPermissionsList(): Array<Permission>;
  setPermissionsList(value: Array<Permission>): void;
  addPermissions(value?: Permission, index?: number): Permission;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPermissionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListPermissionResponse): ListPermissionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListPermissionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPermissionResponse;
  static deserializeBinaryFromReader(message: ListPermissionResponse, reader: jspb.BinaryReader): ListPermissionResponse;
}

export namespace ListPermissionResponse {
  export type AsObject = {
    permissionsList: Array<Permission.AsObject>,
  }
}

export class DeletePermissionRequest extends jspb.Message {
  getServiceid(): string;
  setServiceid(value: string): void;

  getServicename(): string;
  setServicename(value: string): void;

  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePermissionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePermissionRequest): DeletePermissionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeletePermissionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePermissionRequest;
  static deserializeBinaryFromReader(message: DeletePermissionRequest, reader: jspb.BinaryReader): DeletePermissionRequest;
}

export namespace DeletePermissionRequest {
  export type AsObject = {
    serviceid: string,
    servicename: string,
    id: string,
    name: string,
  }
}

export class DeletePermissionResponse extends jspb.Message {
  hasDeleted(): boolean;
  clearDeleted(): void;
  getDeleted(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setDeleted(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePermissionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePermissionResponse): DeletePermissionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeletePermissionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePermissionResponse;
  static deserializeBinaryFromReader(message: DeletePermissionResponse, reader: jspb.BinaryReader): DeletePermissionResponse;
}

export namespace DeletePermissionResponse {
  export type AsObject = {
    deleted?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class AssignPermissionRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getPermission(): string;
  setPermission(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AssignPermissionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AssignPermissionRequest): AssignPermissionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AssignPermissionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AssignPermissionRequest;
  static deserializeBinaryFromReader(message: AssignPermissionRequest, reader: jspb.BinaryReader): AssignPermissionRequest;
}

export namespace AssignPermissionRequest {
  export type AsObject = {
    userid: string,
    permission: string,
  }
}

export class AssignPermissionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AssignPermissionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AssignPermissionResponse): AssignPermissionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AssignPermissionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AssignPermissionResponse;
  static deserializeBinaryFromReader(message: AssignPermissionResponse, reader: jspb.BinaryReader): AssignPermissionResponse;
}

export namespace AssignPermissionResponse {
  export type AsObject = {
  }
}

export class RevokePermissionRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getPermission(): string;
  setPermission(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RevokePermissionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RevokePermissionRequest): RevokePermissionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RevokePermissionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RevokePermissionRequest;
  static deserializeBinaryFromReader(message: RevokePermissionRequest, reader: jspb.BinaryReader): RevokePermissionRequest;
}

export namespace RevokePermissionRequest {
  export type AsObject = {
    userid: string,
    permission: string,
  }
}

export class RevokePermissionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RevokePermissionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RevokePermissionResponse): RevokePermissionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RevokePermissionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RevokePermissionResponse;
  static deserializeBinaryFromReader(message: RevokePermissionResponse, reader: jspb.BinaryReader): RevokePermissionResponse;
}

export namespace RevokePermissionResponse {
  export type AsObject = {
  }
}

export class CreateRoleRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRoleRequest): CreateRoleRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRoleRequest;
  static deserializeBinaryFromReader(message: CreateRoleRequest, reader: jspb.BinaryReader): CreateRoleRequest;
}

export namespace CreateRoleRequest {
  export type AsObject = {
    name: string,
    description: string,
  }
}

export class CreateRoleResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRoleResponse): CreateRoleResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRoleResponse;
  static deserializeBinaryFromReader(message: CreateRoleResponse, reader: jspb.BinaryReader): CreateRoleResponse;
}

export namespace CreateRoleResponse {
  export type AsObject = {
    id: string,
  }
}

export class ListUserRolesRequest extends jspb.Message {
  clearUseridList(): void;
  getUseridList(): Array<string>;
  setUseridList(value: Array<string>): void;
  addUserid(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserRolesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserRolesRequest): ListUserRolesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListUserRolesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserRolesRequest;
  static deserializeBinaryFromReader(message: ListUserRolesRequest, reader: jspb.BinaryReader): ListUserRolesRequest;
}

export namespace ListUserRolesRequest {
  export type AsObject = {
    useridList: Array<string>,
  }
}

export class ListUserRolesResponse extends jspb.Message {
  getRolesMap(): jspb.Map<string, UserRoles>;
  clearRolesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserRolesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserRolesResponse): ListUserRolesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListUserRolesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserRolesResponse;
  static deserializeBinaryFromReader(message: ListUserRolesResponse, reader: jspb.BinaryReader): ListUserRolesResponse;
}

export namespace ListUserRolesResponse {
  export type AsObject = {
    rolesMap: Array<[string, UserRoles.AsObject]>,
  }
}

export class UserRoles extends jspb.Message {
  clearUserrolesList(): void;
  getUserrolesList(): Array<string>;
  setUserrolesList(value: Array<string>): void;
  addUserroles(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserRoles.AsObject;
  static toObject(includeInstance: boolean, msg: UserRoles): UserRoles.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserRoles, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserRoles;
  static deserializeBinaryFromReader(message: UserRoles, reader: jspb.BinaryReader): UserRoles;
}

export namespace UserRoles {
  export type AsObject = {
    userrolesList: Array<string>,
  }
}

export class ListRoleRequest extends jspb.Message {
  getOrgid(): string;
  setOrgid(value: string): void;

  clearIdList(): void;
  getIdList(): Array<string>;
  setIdList(value: Array<string>): void;
  addId(value: string, index?: number): string;

  getSortby(): SortByMap[keyof SortByMap];
  setSortby(value: SortByMap[keyof SortByMap]): void;

  getName(): string;
  setName(value: string): void;

  getUserid(): string;
  setUserid(value: string): void;

  getLimit(): number;
  setLimit(value: number): void;

  getOffset(): number;
  setOffset(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListRoleRequest): ListRoleRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListRoleRequest;
  static deserializeBinaryFromReader(message: ListRoleRequest, reader: jspb.BinaryReader): ListRoleRequest;
}

export namespace ListRoleRequest {
  export type AsObject = {
    orgid: string,
    idList: Array<string>,
    sortby: SortByMap[keyof SortByMap],
    name: string,
    userid: string,
    limit: number,
    offset: number,
  }
}

export class Role extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  clearMembersList(): void;
  getMembersList(): Array<string>;
  setMembersList(value: Array<string>): void;
  addMembers(value: string, index?: number): string;

  clearPermissionsList(): void;
  getPermissionsList(): Array<string>;
  setPermissionsList(value: Array<string>): void;
  addPermissions(value: string, index?: number): string;

  getCreateuid(): string;
  setCreateuid(value: string): void;

  getDeleteuid(): string;
  setDeleteuid(value: string): void;

  hasCreated(): boolean;
  clearCreated(): void;
  getCreated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  getUpdateduid(): string;
  setUpdateduid(value: string): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Role.AsObject;
  static toObject(includeInstance: boolean, msg: Role): Role.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Role, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Role;
  static deserializeBinaryFromReader(message: Role, reader: jspb.BinaryReader): Role;
}

export namespace Role {
  export type AsObject = {
    id: string,
    orgid: string,
    name: string,
    description: string,
    membersList: Array<string>,
    permissionsList: Array<string>,
    createuid: string,
    deleteuid: string,
    created?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    updateduid: string,
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ListRoleResponse extends jspb.Message {
  clearRolesList(): void;
  getRolesList(): Array<Role>;
  setRolesList(value: Array<Role>): void;
  addRoles(value?: Role, index?: number): Role;

  getTotal(): number;
  setTotal(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListRoleResponse): ListRoleResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListRoleResponse;
  static deserializeBinaryFromReader(message: ListRoleResponse, reader: jspb.BinaryReader): ListRoleResponse;
}

export namespace ListRoleResponse {
  export type AsObject = {
    rolesList: Array<Role.AsObject>,
    total: number,
  }
}

export class UpdateRoleRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRoleRequest): UpdateRoleRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRoleRequest;
  static deserializeBinaryFromReader(message: UpdateRoleRequest, reader: jspb.BinaryReader): UpdateRoleRequest;
}

export namespace UpdateRoleRequest {
  export type AsObject = {
    id: string,
    name: string,
    description: string,
  }
}

export class UpdateRoleResponse extends jspb.Message {
  hasRole(): boolean;
  clearRole(): void;
  getRole(): Role | undefined;
  setRole(value?: Role): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRoleResponse): UpdateRoleResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRoleResponse;
  static deserializeBinaryFromReader(message: UpdateRoleResponse, reader: jspb.BinaryReader): UpdateRoleResponse;
}

export namespace UpdateRoleResponse {
  export type AsObject = {
    role?: Role.AsObject,
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class DeleteRoleRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRoleRequest): DeleteRoleRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRoleRequest;
  static deserializeBinaryFromReader(message: DeleteRoleRequest, reader: jspb.BinaryReader): DeleteRoleRequest;
}

export namespace DeleteRoleRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeleteRoleResponse extends jspb.Message {
  hasDeleted(): boolean;
  clearDeleted(): void;
  getDeleted(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setDeleted(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRoleResponse): DeleteRoleResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeleteRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRoleResponse;
  static deserializeBinaryFromReader(message: DeleteRoleResponse, reader: jspb.BinaryReader): DeleteRoleResponse;
}

export namespace DeleteRoleResponse {
  export type AsObject = {
    deleted?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class AssignRolePermissionRequest extends jspb.Message {
  getRoleid(): string;
  setRoleid(value: string): void;

  getPermission(): string;
  setPermission(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AssignRolePermissionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AssignRolePermissionRequest): AssignRolePermissionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AssignRolePermissionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AssignRolePermissionRequest;
  static deserializeBinaryFromReader(message: AssignRolePermissionRequest, reader: jspb.BinaryReader): AssignRolePermissionRequest;
}

export namespace AssignRolePermissionRequest {
  export type AsObject = {
    roleid: string,
    permission: string,
  }
}

export class AssignRolePermissionResponse extends jspb.Message {
  hasRole(): boolean;
  clearRole(): void;
  getRole(): Role | undefined;
  setRole(value?: Role): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AssignRolePermissionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AssignRolePermissionResponse): AssignRolePermissionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AssignRolePermissionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AssignRolePermissionResponse;
  static deserializeBinaryFromReader(message: AssignRolePermissionResponse, reader: jspb.BinaryReader): AssignRolePermissionResponse;
}

export namespace AssignRolePermissionResponse {
  export type AsObject = {
    role?: Role.AsObject,
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class RevokeRolePermissionRequest extends jspb.Message {
  getRoleid(): string;
  setRoleid(value: string): void;

  getPermission(): string;
  setPermission(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RevokeRolePermissionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RevokeRolePermissionRequest): RevokeRolePermissionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RevokeRolePermissionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RevokeRolePermissionRequest;
  static deserializeBinaryFromReader(message: RevokeRolePermissionRequest, reader: jspb.BinaryReader): RevokeRolePermissionRequest;
}

export namespace RevokeRolePermissionRequest {
  export type AsObject = {
    roleid: string,
    permission: string,
  }
}

export class RevokeRolePermissionResponse extends jspb.Message {
  hasRole(): boolean;
  clearRole(): void;
  getRole(): Role | undefined;
  setRole(value?: Role): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RevokeRolePermissionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RevokeRolePermissionResponse): RevokeRolePermissionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RevokeRolePermissionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RevokeRolePermissionResponse;
  static deserializeBinaryFromReader(message: RevokeRolePermissionResponse, reader: jspb.BinaryReader): RevokeRolePermissionResponse;
}

export namespace RevokeRolePermissionResponse {
  export type AsObject = {
    role?: Role.AsObject,
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class AddUserRequest extends jspb.Message {
  getRoleid(): string;
  setRoleid(value: string): void;

  getUserid(): string;
  setUserid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddUserRequest): AddUserRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddUserRequest;
  static deserializeBinaryFromReader(message: AddUserRequest, reader: jspb.BinaryReader): AddUserRequest;
}

export namespace AddUserRequest {
  export type AsObject = {
    roleid: string,
    userid: string,
  }
}

export class AddUserResponse extends jspb.Message {
  hasRole(): boolean;
  clearRole(): void;
  getRole(): Role | undefined;
  setRole(value?: Role): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddUserResponse): AddUserResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddUserResponse;
  static deserializeBinaryFromReader(message: AddUserResponse, reader: jspb.BinaryReader): AddUserResponse;
}

export namespace AddUserResponse {
  export type AsObject = {
    role?: Role.AsObject,
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class RemoveUserRequest extends jspb.Message {
  getRoleid(): string;
  setRoleid(value: string): void;

  getUserid(): string;
  setUserid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveUserRequest): RemoveUserRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RemoveUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveUserRequest;
  static deserializeBinaryFromReader(message: RemoveUserRequest, reader: jspb.BinaryReader): RemoveUserRequest;
}

export namespace RemoveUserRequest {
  export type AsObject = {
    roleid: string,
    userid: string,
  }
}

export class RemoveUserResponse extends jspb.Message {
  hasRole(): boolean;
  clearRole(): void;
  getRole(): Role | undefined;
  setRole(value?: Role): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveUserResponse): RemoveUserResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RemoveUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveUserResponse;
  static deserializeBinaryFromReader(message: RemoveUserResponse, reader: jspb.BinaryReader): RemoveUserResponse;
}

export namespace RemoveUserResponse {
  export type AsObject = {
    role?: Role.AsObject,
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class GrantServiceRequest extends jspb.Message {
  getServiceid(): string;
  setServiceid(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getEnvironment(): string;
  setEnvironment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrantServiceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GrantServiceRequest): GrantServiceRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrantServiceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrantServiceRequest;
  static deserializeBinaryFromReader(message: GrantServiceRequest, reader: jspb.BinaryReader): GrantServiceRequest;
}

export namespace GrantServiceRequest {
  export type AsObject = {
    serviceid: string,
    orgid: string,
    environment: string,
  }
}

export class GrantServiceResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrantServiceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GrantServiceResponse): GrantServiceResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrantServiceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrantServiceResponse;
  static deserializeBinaryFromReader(message: GrantServiceResponse, reader: jspb.BinaryReader): GrantServiceResponse;
}

export namespace GrantServiceResponse {
  export type AsObject = {
    id: string,
  }
}

export class ListServiceAssignmentsRequest extends jspb.Message {
  getOrgid(): string;
  setOrgid(value: string): void;

  clearServiceidList(): void;
  getServiceidList(): Array<string>;
  setServiceidList(value: Array<string>): void;
  addServiceid(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListServiceAssignmentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListServiceAssignmentsRequest): ListServiceAssignmentsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListServiceAssignmentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListServiceAssignmentsRequest;
  static deserializeBinaryFromReader(message: ListServiceAssignmentsRequest, reader: jspb.BinaryReader): ListServiceAssignmentsRequest;
}

export namespace ListServiceAssignmentsRequest {
  export type AsObject = {
    orgid: string,
    serviceidList: Array<string>,
  }
}

export class ServiceAssignment extends jspb.Message {
  getGrant(): string;
  setGrant(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getServiceid(): string;
  setServiceid(value: string): void;

  getServicename(): string;
  setServicename(value: string): void;

  getEnvironment(): string;
  setEnvironment(value: string): void;

  getGrantedby(): string;
  setGrantedby(value: string): void;

  hasGranted(): boolean;
  clearGranted(): void;
  getGranted(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setGranted(value?: google_protobuf_timestamp_pb.Timestamp): void;

  getRevokedby(): string;
  setRevokedby(value: string): void;

  hasRevoked(): boolean;
  clearRevoked(): void;
  getRevoked(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setRevoked(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServiceAssignment.AsObject;
  static toObject(includeInstance: boolean, msg: ServiceAssignment): ServiceAssignment.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ServiceAssignment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServiceAssignment;
  static deserializeBinaryFromReader(message: ServiceAssignment, reader: jspb.BinaryReader): ServiceAssignment;
}

export namespace ServiceAssignment {
  export type AsObject = {
    grant: string,
    orgid: string,
    serviceid: string,
    servicename: string,
    environment: string,
    grantedby: string,
    granted?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    revokedby: string,
    revoked?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ListServiceAssignmentsResponse extends jspb.Message {
  clearServiceassignmentsList(): void;
  getServiceassignmentsList(): Array<ServiceAssignment>;
  setServiceassignmentsList(value: Array<ServiceAssignment>): void;
  addServiceassignments(value?: ServiceAssignment, index?: number): ServiceAssignment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListServiceAssignmentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListServiceAssignmentsResponse): ListServiceAssignmentsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListServiceAssignmentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListServiceAssignmentsResponse;
  static deserializeBinaryFromReader(message: ListServiceAssignmentsResponse, reader: jspb.BinaryReader): ListServiceAssignmentsResponse;
}

export namespace ListServiceAssignmentsResponse {
  export type AsObject = {
    serviceassignmentsList: Array<ServiceAssignment.AsObject>,
  }
}

export class RevokeServiceRequest extends jspb.Message {
  getServiceid(): string;
  setServiceid(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getEnvironment(): string;
  setEnvironment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RevokeServiceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RevokeServiceRequest): RevokeServiceRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RevokeServiceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RevokeServiceRequest;
  static deserializeBinaryFromReader(message: RevokeServiceRequest, reader: jspb.BinaryReader): RevokeServiceRequest;
}

export namespace RevokeServiceRequest {
  export type AsObject = {
    serviceid: string,
    orgid: string,
    environment: string,
  }
}

export class RevokeServiceResponse extends jspb.Message {
  hasRevoked(): boolean;
  clearRevoked(): void;
  getRevoked(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setRevoked(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RevokeServiceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RevokeServiceResponse): RevokeServiceResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RevokeServiceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RevokeServiceResponse;
  static deserializeBinaryFromReader(message: RevokeServiceResponse, reader: jspb.BinaryReader): RevokeServiceResponse;
}

export namespace RevokeServiceResponse {
  export type AsObject = {
    revoked?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class PublicServiceRequest extends jspb.Message {
  getServiceid(): string;
  setServiceid(value: string): void;

  getEnvironment(): string;
  setEnvironment(value: string): void;

  getEnabled(): boolean;
  setEnabled(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublicServiceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PublicServiceRequest): PublicServiceRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PublicServiceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublicServiceRequest;
  static deserializeBinaryFromReader(message: PublicServiceRequest, reader: jspb.BinaryReader): PublicServiceRequest;
}

export namespace PublicServiceRequest {
  export type AsObject = {
    serviceid: string,
    environment: string,
    enabled: boolean,
  }
}

export class PublicServiceResponse extends jspb.Message {
  hasProcessed(): boolean;
  clearProcessed(): void;
  getProcessed(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setProcessed(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublicServiceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PublicServiceResponse): PublicServiceResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PublicServiceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublicServiceResponse;
  static deserializeBinaryFromReader(message: PublicServiceResponse, reader: jspb.BinaryReader): PublicServiceResponse;
}

export namespace PublicServiceResponse {
  export type AsObject = {
    processed?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ListServicesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListServicesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListServicesRequest): ListServicesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListServicesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListServicesRequest;
  static deserializeBinaryFromReader(message: ListServicesRequest, reader: jspb.BinaryReader): ListServicesRequest;
}

export namespace ListServicesRequest {
  export type AsObject = {
  }
}

export class ListServicesResponse extends jspb.Message {
  clearServicesList(): void;
  getServicesList(): Array<Service>;
  setServicesList(value: Array<Service>): void;
  addServices(value?: Service, index?: number): Service;

  clearPublicList(): void;
  getPublicList(): Array<Service>;
  setPublicList(value: Array<Service>): void;
  addPublic(value?: Service, index?: number): Service;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListServicesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListServicesResponse): ListServicesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListServicesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListServicesResponse;
  static deserializeBinaryFromReader(message: ListServicesResponse, reader: jspb.BinaryReader): ListServicesResponse;
}

export namespace ListServicesResponse {
  export type AsObject = {
    servicesList: Array<Service.AsObject>,
    publicList: Array<Service.AsObject>,
  }
}

export class Service extends jspb.Message {
  getServiceid(): string;
  setServiceid(value: string): void;

  getServicename(): string;
  setServicename(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getEnvironment(): string;
  setEnvironment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Service.AsObject;
  static toObject(includeInstance: boolean, msg: Service): Service.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Service, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Service;
  static deserializeBinaryFromReader(message: Service, reader: jspb.BinaryReader): Service;
}

export namespace Service {
  export type AsObject = {
    serviceid: string,
    servicename: string,
    description: string,
    environment: string,
  }
}

export class ValidatePermissionRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getAction(): string;
  setAction(value: string): void;

  getResource(): string;
  setResource(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getEnvironment(): string;
  setEnvironment(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidatePermissionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ValidatePermissionRequest): ValidatePermissionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValidatePermissionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidatePermissionRequest;
  static deserializeBinaryFromReader(message: ValidatePermissionRequest, reader: jspb.BinaryReader): ValidatePermissionRequest;
}

export namespace ValidatePermissionRequest {
  export type AsObject = {
    id: string,
    action: string,
    resource: string,
    orgid: string,
    environment: string,
  }
}

export class ValidatePermissionResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  hasValidated(): boolean;
  clearValidated(): void;
  getValidated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setValidated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidatePermissionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ValidatePermissionResponse): ValidatePermissionResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValidatePermissionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidatePermissionResponse;
  static deserializeBinaryFromReader(message: ValidatePermissionResponse, reader: jspb.BinaryReader): ValidatePermissionResponse;
}

export namespace ValidatePermissionResponse {
  export type AsObject = {
    id: string,
    name: string,
    orgid: string,
    validated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export interface SortByMap {
  DESC: 0;
  ASC: 1;
}

export const SortBy: SortByMap;

