// package: brankas.rbac.v1.serviceaccount
// file: brank.as/rbac/gunk/v1/serviceaccount/all.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class CreateAccountRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getEnv(): string;
  setEnv(value: string): void;

  getRole(): string;
  setRole(value: string): void;

  getAuthtype(): AuthTypeMap[keyof AuthTypeMap];
  setAuthtype(value: AuthTypeMap[keyof AuthTypeMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAccountRequest): CreateAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountRequest;
  static deserializeBinaryFromReader(message: CreateAccountRequest, reader: jspb.BinaryReader): CreateAccountRequest;
}

export namespace CreateAccountRequest {
  export type AsObject = {
    name: string,
    env: string,
    role: string,
    authtype: AuthTypeMap[keyof AuthTypeMap],
  }
}

export class CreateAccountResponse extends jspb.Message {
  getClientid(): string;
  setClientid(value: string): void;

  getSecret(): string;
  setSecret(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAccountResponse): CreateAccountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountResponse;
  static deserializeBinaryFromReader(message: CreateAccountResponse, reader: jspb.BinaryReader): CreateAccountResponse;
}

export namespace CreateAccountResponse {
  export type AsObject = {
    clientid: string,
    secret: string,
  }
}

export class ListAccountsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAccountsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAccountsRequest): ListAccountsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAccountsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAccountsRequest;
  static deserializeBinaryFromReader(message: ListAccountsRequest, reader: jspb.BinaryReader): ListAccountsRequest;
}

export namespace ListAccountsRequest {
  export type AsObject = {
  }
}

export class ServiceAccount extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getEnv(): string;
  setEnv(value: string): void;

  getClientid(): string;
  setClientid(value: string): void;

  getCreator(): string;
  setCreator(value: string): void;

  hasCreated(): boolean;
  clearCreated(): void;
  getCreated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setDisabled(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServiceAccount.AsObject;
  static toObject(includeInstance: boolean, msg: ServiceAccount): ServiceAccount.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ServiceAccount, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServiceAccount;
  static deserializeBinaryFromReader(message: ServiceAccount, reader: jspb.BinaryReader): ServiceAccount;
}

export namespace ServiceAccount {
  export type AsObject = {
    name: string,
    env: string,
    clientid: string,
    creator: string,
    created?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    disabled?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ListAccountsResponse extends jspb.Message {
  clearAccountsList(): void;
  getAccountsList(): Array<ServiceAccount>;
  setAccountsList(value: Array<ServiceAccount>): void;
  addAccounts(value?: ServiceAccount, index?: number): ServiceAccount;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAccountsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAccountsResponse): ListAccountsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListAccountsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAccountsResponse;
  static deserializeBinaryFromReader(message: ListAccountsResponse, reader: jspb.BinaryReader): ListAccountsResponse;
}

export namespace ListAccountsResponse {
  export type AsObject = {
    accountsList: Array<ServiceAccount.AsObject>,
  }
}

export class DisableAccountRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAccountRequest): DisableAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAccountRequest;
  static deserializeBinaryFromReader(message: DisableAccountRequest, reader: jspb.BinaryReader): DisableAccountRequest;
}

export namespace DisableAccountRequest {
  export type AsObject = {
    name: string,
  }
}

export class DisableAccountResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAccountResponse): DisableAccountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAccountResponse;
  static deserializeBinaryFromReader(message: DisableAccountResponse, reader: jspb.BinaryReader): DisableAccountResponse;
}

export namespace DisableAccountResponse {
  export type AsObject = {
  }
}

export class ValidateAccountRequest extends jspb.Message {
  getClientid(): string;
  setClientid(value: string): void;

  getOperation(): string;
  setOperation(value: string): void;

  getResource(): string;
  setResource(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateAccountRequest): ValidateAccountRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValidateAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateAccountRequest;
  static deserializeBinaryFromReader(message: ValidateAccountRequest, reader: jspb.BinaryReader): ValidateAccountRequest;
}

export namespace ValidateAccountRequest {
  export type AsObject = {
    clientid: string,
    operation: string,
    resource: string,
  }
}

export class ValidateAccountResponse extends jspb.Message {
  getEnvironment(): string;
  setEnvironment(value: string): void;

  getClientname(): string;
  setClientname(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateAccountResponse): ValidateAccountResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValidateAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateAccountResponse;
  static deserializeBinaryFromReader(message: ValidateAccountResponse, reader: jspb.BinaryReader): ValidateAccountResponse;
}

export namespace ValidateAccountResponse {
  export type AsObject = {
    environment: string,
    clientname: string,
    orgid: string,
  }
}

export class ValidateAPIKeyRequest extends jspb.Message {
  getApikey(): string;
  setApikey(value: string): void;

  getEndpoint(): string;
  setEndpoint(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateAPIKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateAPIKeyRequest): ValidateAPIKeyRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValidateAPIKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateAPIKeyRequest;
  static deserializeBinaryFromReader(message: ValidateAPIKeyRequest, reader: jspb.BinaryReader): ValidateAPIKeyRequest;
}

export namespace ValidateAPIKeyRequest {
  export type AsObject = {
    apikey: string,
    endpoint: string,
  }
}

export class ValidateAPIKeyResponse extends jspb.Message {
  getOrgid(): string;
  setOrgid(value: string): void;

  getEnv(): string;
  setEnv(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateAPIKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateAPIKeyResponse): ValidateAPIKeyResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValidateAPIKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateAPIKeyResponse;
  static deserializeBinaryFromReader(message: ValidateAPIKeyResponse, reader: jspb.BinaryReader): ValidateAPIKeyResponse;
}

export namespace ValidateAPIKeyResponse {
  export type AsObject = {
    orgid: string,
    env: string,
  }
}

export interface AuthTypeMap {
  OAUTH2: 0;
  APIKEY: 1;
}

export const AuthType: AuthTypeMap;

