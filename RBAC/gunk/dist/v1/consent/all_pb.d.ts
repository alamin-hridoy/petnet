// package: rbac.brankas.consent
// file: brank.as/rbac/gunk/v1/consent/all.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class ServeGrantRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getClientid(): string;
  setClientid(value: string): void;

  getOwnerid(): string;
  setOwnerid(value: string): void;

  clearRequestedList(): void;
  getRequestedList(): Array<string>;
  setRequestedList(value: Array<string>): void;
  addRequested(value: string, index?: number): string;

  clearGrantedList(): void;
  getGrantedList(): Array<string>;
  setGrantedList(value: Array<string>): void;
  addGranted(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServeGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ServeGrantRequest): ServeGrantRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ServeGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServeGrantRequest;
  static deserializeBinaryFromReader(message: ServeGrantRequest, reader: jspb.BinaryReader): ServeGrantRequest;
}

export namespace ServeGrantRequest {
  export type AsObject = {
    userid: string,
    clientid: string,
    ownerid: string,
    requestedList: Array<string>,
    grantedList: Array<string>,
  }
}

export class ServeGrantResponse extends jspb.Message {
  getNewscopesMap(): jspb.Map<string, ScopeDetail>;
  clearNewscopesMap(): void;
  getGrantedscopesMap(): jspb.Map<string, ScopeDetail>;
  clearGrantedscopesMap(): void;
  getGroupsMap(): jspb.Map<string, GroupDetail>;
  clearGroupsMap(): void;
  getSkip(): boolean;
  setSkip(value: boolean): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getOrgname(): string;
  setOrgname(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServeGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ServeGrantResponse): ServeGrantResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ServeGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServeGrantResponse;
  static deserializeBinaryFromReader(message: ServeGrantResponse, reader: jspb.BinaryReader): ServeGrantResponse;
}

export namespace ServeGrantResponse {
  export type AsObject = {
    newscopesMap: Array<[string, ScopeDetail.AsObject]>,
    grantedscopesMap: Array<[string, ScopeDetail.AsObject]>,
    groupsMap: Array<[string, GroupDetail.AsObject]>,
    skip: boolean,
    orgid: string,
    orgname: string,
  }
}

export class GrantRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getClientid(): string;
  setClientid(value: string): void;

  getOwnerid(): string;
  setOwnerid(value: string): void;

  clearScopesList(): void;
  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): void;
  addScopes(value: string, index?: number): string;

  hasTimestamp(): boolean;
  clearTimestamp(): void;
  getTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GrantRequest): GrantRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrantRequest;
  static deserializeBinaryFromReader(message: GrantRequest, reader: jspb.BinaryReader): GrantRequest;
}

export namespace GrantRequest {
  export type AsObject = {
    userid: string,
    clientid: string,
    ownerid: string,
    scopesList: Array<string>,
    timestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class GrantResponse extends jspb.Message {
  getGrantid(): string;
  setGrantid(value: string): void;

  getSinglegrant(): boolean;
  setSinglegrant(value: boolean): void;

  clearGrantsList(): void;
  getGrantsList(): Array<string>;
  setGrantsList(value: Array<string>): void;
  addGrants(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GrantResponse): GrantResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrantResponse;
  static deserializeBinaryFromReader(message: GrantResponse, reader: jspb.BinaryReader): GrantResponse;
}

export namespace GrantResponse {
  export type AsObject = {
    grantid: string,
    singlegrant: boolean,
    grantsList: Array<string>,
  }
}

export class ConsentError extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): void;

  getErrordetailsMap(): jspb.Map<string, string>;
  clearErrordetailsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConsentError.AsObject;
  static toObject(includeInstance: boolean, msg: ConsentError): ConsentError.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConsentError, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConsentError;
  static deserializeBinaryFromReader(message: ConsentError, reader: jspb.BinaryReader): ConsentError;
}

export namespace ConsentError {
  export type AsObject = {
    message: string,
    errordetailsMap: Array<[string, string]>,
  }
}

export class UpsertScopeRequest extends jspb.Message {
  getScope(): string;
  setScope(value: string): void;

  getName(): string;
  setName(value: string): void;

  getGroupname(): string;
  setGroupname(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertScopeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertScopeRequest): UpsertScopeRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpsertScopeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertScopeRequest;
  static deserializeBinaryFromReader(message: UpsertScopeRequest, reader: jspb.BinaryReader): UpsertScopeRequest;
}

export namespace UpsertScopeRequest {
  export type AsObject = {
    scope: string,
    name: string,
    groupname: string,
    description: string,
  }
}

export class UpsertScopeResponse extends jspb.Message {
  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertScopeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertScopeResponse): UpsertScopeResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpsertScopeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertScopeResponse;
  static deserializeBinaryFromReader(message: UpsertScopeResponse, reader: jspb.BinaryReader): UpsertScopeResponse;
}

export namespace UpsertScopeResponse {
  export type AsObject = {
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class UpdateGroupRequest extends jspb.Message {
  getGroupname(): string;
  setGroupname(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGroupRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGroupRequest): UpdateGroupRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateGroupRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGroupRequest;
  static deserializeBinaryFromReader(message: UpdateGroupRequest, reader: jspb.BinaryReader): UpdateGroupRequest;
}

export namespace UpdateGroupRequest {
  export type AsObject = {
    groupname: string,
    description: string,
  }
}

export class UpdateGroupResponse extends jspb.Message {
  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGroupResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGroupResponse): UpdateGroupResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateGroupResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGroupResponse;
  static deserializeBinaryFromReader(message: UpdateGroupResponse, reader: jspb.BinaryReader): UpdateGroupResponse;
}

export namespace UpdateGroupResponse {
  export type AsObject = {
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class GetScopeRequest extends jspb.Message {
  clearScopesList(): void;
  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): void;
  addScopes(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetScopeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetScopeRequest): GetScopeRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetScopeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetScopeRequest;
  static deserializeBinaryFromReader(message: GetScopeRequest, reader: jspb.BinaryReader): GetScopeRequest;
}

export namespace GetScopeRequest {
  export type AsObject = {
    scopesList: Array<string>,
  }
}

export class GetScopeResponse extends jspb.Message {
  getScopesMap(): jspb.Map<string, ScopeDetail>;
  clearScopesMap(): void;
  getGroupsMap(): jspb.Map<string, GroupDetail>;
  clearGroupsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetScopeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetScopeResponse): GetScopeResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetScopeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetScopeResponse;
  static deserializeBinaryFromReader(message: GetScopeResponse, reader: jspb.BinaryReader): GetScopeResponse;
}

export namespace GetScopeResponse {
  export type AsObject = {
    scopesMap: Array<[string, ScopeDetail.AsObject]>,
    groupsMap: Array<[string, GroupDetail.AsObject]>,
  }
}

export class GroupDetail extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GroupDetail.AsObject;
  static toObject(includeInstance: boolean, msg: GroupDetail): GroupDetail.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GroupDetail, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GroupDetail;
  static deserializeBinaryFromReader(message: GroupDetail, reader: jspb.BinaryReader): GroupDetail;
}

export namespace GroupDetail {
  export type AsObject = {
    name: string,
    description: string,
  }
}

export class ScopeDetail extends jspb.Message {
  getScope(): string;
  setScope(value: string): void;

  getName(): string;
  setName(value: string): void;

  getGroup(): string;
  setGroup(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ScopeDetail.AsObject;
  static toObject(includeInstance: boolean, msg: ScopeDetail): ScopeDetail.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ScopeDetail, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ScopeDetail;
  static deserializeBinaryFromReader(message: ScopeDetail, reader: jspb.BinaryReader): ScopeDetail;
}

export namespace ScopeDetail {
  export type AsObject = {
    scope: string,
    name: string,
    group: string,
    description: string,
  }
}

