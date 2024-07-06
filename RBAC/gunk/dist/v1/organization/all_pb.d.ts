// package: organization
// file: brank.as/rbac/gunk/v1/organization/all.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import * as brank_as_rbac_gunk_v1_mfa_all_pb from "./all_pb";

export class GetOrganizationRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrganizationRequest): GetOrganizationRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrganizationRequest;
  static deserializeBinaryFromReader(message: GetOrganizationRequest, reader: jspb.BinaryReader): GetOrganizationRequest;
}

export namespace GetOrganizationRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetOrganizationResponse extends jspb.Message {
  clearOrganizationList(): void;
  getOrganizationList(): Array<Organization>;
  setOrganizationList(value: Array<Organization>): void;
  addOrganization(value?: Organization, index?: number): Organization;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrganizationResponse): GetOrganizationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrganizationResponse;
  static deserializeBinaryFromReader(message: GetOrganizationResponse, reader: jspb.BinaryReader): GetOrganizationResponse;
}

export namespace GetOrganizationResponse {
  export type AsObject = {
    organizationList: Array<Organization.AsObject>,
  }
}

export class Organization extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getPhone(): string;
  setPhone(value: string): void;

  getActive(): boolean;
  setActive(value: boolean): void;

  getLoginmfa(): boolean;
  setLoginmfa(value: boolean): void;

  hasCreated(): boolean;
  clearCreated(): void;
  getCreated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Organization.AsObject;
  static toObject(includeInstance: boolean, msg: Organization): Organization.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Organization, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Organization;
  static deserializeBinaryFromReader(message: Organization, reader: jspb.BinaryReader): Organization;
}

export namespace Organization {
  export type AsObject = {
    id: string,
    name: string,
    email: string,
    phone: string,
    active: boolean,
    loginmfa: boolean,
    created?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class UpdateOrganizationRequest extends jspb.Message {
  getOrganizationid(): string;
  setOrganizationid(value: string): void;

  getName(): string;
  setName(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getPhone(): string;
  setPhone(value: string): void;

  getLoginmfa(): EnableOptMap[keyof EnableOptMap];
  setLoginmfa(value: EnableOptMap[keyof EnableOptMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrganizationRequest): UpdateOrganizationRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrganizationRequest;
  static deserializeBinaryFromReader(message: UpdateOrganizationRequest, reader: jspb.BinaryReader): UpdateOrganizationRequest;
}

export namespace UpdateOrganizationRequest {
  export type AsObject = {
    organizationid: string,
    name: string,
    email: string,
    phone: string,
    loginmfa: EnableOptMap[keyof EnableOptMap],
  }
}

export class UpdateOrganizationResponse extends jspb.Message {
  getMfaeventid(): string;
  setMfaeventid(value: string): void;

  getMfatype(): brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap];
  setMfatype(value: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap]): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrganizationResponse): UpdateOrganizationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrganizationResponse;
  static deserializeBinaryFromReader(message: UpdateOrganizationResponse, reader: jspb.BinaryReader): UpdateOrganizationResponse;
}

export namespace UpdateOrganizationResponse {
  export type AsObject = {
    mfaeventid: string,
    mfatype: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap],
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ConfirmUpdateRequest extends jspb.Message {
  getOrganizationid(): string;
  setOrganizationid(value: string): void;

  getMfaeventid(): string;
  setMfaeventid(value: string): void;

  getMfatype(): brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap];
  setMfatype(value: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap]): void;

  getMfatoken(): string;
  setMfatoken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfirmUpdateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ConfirmUpdateRequest): ConfirmUpdateRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConfirmUpdateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfirmUpdateRequest;
  static deserializeBinaryFromReader(message: ConfirmUpdateRequest, reader: jspb.BinaryReader): ConfirmUpdateRequest;
}

export namespace ConfirmUpdateRequest {
  export type AsObject = {
    organizationid: string,
    mfaeventid: string,
    mfatype: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap],
    mfatoken: string,
  }
}

export class ConfirmUpdateResponse extends jspb.Message {
  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfirmUpdateResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ConfirmUpdateResponse): ConfirmUpdateResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ConfirmUpdateResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfirmUpdateResponse;
  static deserializeBinaryFromReader(message: ConfirmUpdateResponse, reader: jspb.BinaryReader): ConfirmUpdateResponse;
}

export namespace ConfirmUpdateResponse {
  export type AsObject = {
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export interface EnableOptMap {
  NOCHANGE: 0;
  ENABLE: 1;
  DISABLE: 2;
}

export const EnableOpt: EnableOptMap;

