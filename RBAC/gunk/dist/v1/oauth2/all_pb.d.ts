// package: brankas.rbac.v1.oauth2
// file: brank.as/rbac/gunk/v1/oauth2/all.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class CreateClientRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getAudience(): string;
  setAudience(value: string): void;

  getEnv(): string;
  setEnv(value: string): void;

  getRole(): string;
  setRole(value: string): void;

  getClienttype(): ClientTypeMap[keyof ClientTypeMap];
  setClienttype(value: ClientTypeMap[keyof ClientTypeMap]): void;

  clearCorsList(): void;
  getCorsList(): Array<string>;
  setCorsList(value: Array<string>): void;
  addCors(value: string, index?: number): string;

  getLogourl(): string;
  setLogourl(value: string): void;

  clearScopesList(): void;
  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): void;
  addScopes(value: string, index?: number): string;

  clearRedirecturlList(): void;
  getRedirecturlList(): Array<string>;
  setRedirecturlList(value: Array<string>): void;
  addRedirecturl(value: string, index?: number): string;

  getLogoutredirecturl(): string;
  setLogoutredirecturl(value: string): void;

  getIdentitysource(): string;
  setIdentitysource(value: string): void;

  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): ClientConfig | undefined;
  setConfig(value?: ClientConfig): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateClientRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateClientRequest): CreateClientRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateClientRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateClientRequest;
  static deserializeBinaryFromReader(message: CreateClientRequest, reader: jspb.BinaryReader): CreateClientRequest;
}

export namespace CreateClientRequest {
  export type AsObject = {
    name: string,
    audience: string,
    env: string,
    role: string,
    clienttype: ClientTypeMap[keyof ClientTypeMap],
    corsList: Array<string>,
    logourl: string,
    scopesList: Array<string>,
    redirecturlList: Array<string>,
    logoutredirecturl: string,
    identitysource: string,
    config?: ClientConfig.AsObject,
  }
}

export class ClientConfig extends jspb.Message {
  getLogintemplate(): string;
  setLogintemplate(value: string): void;

  getOtptemplate(): string;
  setOtptemplate(value: string): void;

  getConsenttemplate(): string;
  setConsenttemplate(value: string): void;

  getForceconsent(): boolean;
  setForceconsent(value: boolean): void;

  hasSessionduration(): boolean;
  clearSessionduration(): void;
  getSessionduration(): google_protobuf_duration_pb.Duration | undefined;
  setSessionduration(value?: google_protobuf_duration_pb.Duration): void;

  getIdentitysource(): string;
  setIdentitysource(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClientConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ClientConfig): ClientConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ClientConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClientConfig;
  static deserializeBinaryFromReader(message: ClientConfig, reader: jspb.BinaryReader): ClientConfig;
}

export namespace ClientConfig {
  export type AsObject = {
    logintemplate: string,
    otptemplate: string,
    consenttemplate: string,
    forceconsent: boolean,
    sessionduration?: google_protobuf_duration_pb.Duration.AsObject,
    identitysource: string,
  }
}

export class CreateClientResponse extends jspb.Message {
  getClientid(): string;
  setClientid(value: string): void;

  getSecret(): string;
  setSecret(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateClientResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateClientResponse): CreateClientResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateClientResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateClientResponse;
  static deserializeBinaryFromReader(message: CreateClientResponse, reader: jspb.BinaryReader): CreateClientResponse;
}

export namespace CreateClientResponse {
  export type AsObject = {
    clientid: string,
    secret: string,
  }
}

export class UpdateClientRequest extends jspb.Message {
  getClientid(): string;
  setClientid(value: string): void;

  getName(): string;
  setName(value: string): void;

  clearCorsList(): void;
  getCorsList(): Array<string>;
  setCorsList(value: Array<string>): void;
  addCors(value: string, index?: number): string;

  getLogourl(): string;
  setLogourl(value: string): void;

  clearScopesList(): void;
  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): void;
  addScopes(value: string, index?: number): string;

  clearRedirecturlList(): void;
  getRedirecturlList(): Array<string>;
  setRedirecturlList(value: Array<string>): void;
  addRedirecturl(value: string, index?: number): string;

  getLogoutredirecturl(): string;
  setLogoutredirecturl(value: string): void;

  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): ClientConfig | undefined;
  setConfig(value?: ClientConfig): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateClientRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateClientRequest): UpdateClientRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateClientRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateClientRequest;
  static deserializeBinaryFromReader(message: UpdateClientRequest, reader: jspb.BinaryReader): UpdateClientRequest;
}

export namespace UpdateClientRequest {
  export type AsObject = {
    clientid: string,
    name: string,
    corsList: Array<string>,
    logourl: string,
    scopesList: Array<string>,
    redirecturlList: Array<string>,
    logoutredirecturl: string,
    config?: ClientConfig.AsObject,
  }
}

export class UpdateClientResponse extends jspb.Message {
  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateClientResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateClientResponse): UpdateClientResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateClientResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateClientResponse;
  static deserializeBinaryFromReader(message: UpdateClientResponse, reader: jspb.BinaryReader): UpdateClientResponse;
}

export namespace UpdateClientResponse {
  export type AsObject = {
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ListClientsRequest extends jspb.Message {
  getClientid(): string;
  setClientid(value: string): void;

  getEnv(): string;
  setEnv(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getListdisable(): boolean;
  setListdisable(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListClientsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListClientsRequest): ListClientsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListClientsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListClientsRequest;
  static deserializeBinaryFromReader(message: ListClientsRequest, reader: jspb.BinaryReader): ListClientsRequest;
}

export namespace ListClientsRequest {
  export type AsObject = {
    clientid: string,
    env: string,
    orgid: string,
    listdisable: boolean,
  }
}

export class Oauth2Client extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getEnv(): string;
  setEnv(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getClientid(): string;
  setClientid(value: string): void;

  getLogourl(): string;
  setLogourl(value: string): void;

  clearScopesList(): void;
  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): void;
  addScopes(value: string, index?: number): string;

  clearRedirecturlList(): void;
  getRedirecturlList(): Array<string>;
  setRedirecturlList(value: Array<string>): void;
  addRedirecturl(value: string, index?: number): string;

  getLogoutredirecturl(): string;
  setLogoutredirecturl(value: string): void;

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

  hasConfig(): boolean;
  clearConfig(): void;
  getConfig(): ClientConfig | undefined;
  setConfig(value?: ClientConfig): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Oauth2Client.AsObject;
  static toObject(includeInstance: boolean, msg: Oauth2Client): Oauth2Client.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Oauth2Client, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Oauth2Client;
  static deserializeBinaryFromReader(message: Oauth2Client, reader: jspb.BinaryReader): Oauth2Client;
}

export namespace Oauth2Client {
  export type AsObject = {
    name: string,
    env: string,
    orgid: string,
    clientid: string,
    logourl: string,
    scopesList: Array<string>,
    redirecturlList: Array<string>,
    logoutredirecturl: string,
    creator: string,
    created?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    disabled?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    config?: ClientConfig.AsObject,
  }
}

export class ListClientsResponse extends jspb.Message {
  clearClientsList(): void;
  getClientsList(): Array<Oauth2Client>;
  setClientsList(value: Array<Oauth2Client>): void;
  addClients(value?: Oauth2Client, index?: number): Oauth2Client;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListClientsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListClientsResponse): ListClientsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListClientsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListClientsResponse;
  static deserializeBinaryFromReader(message: ListClientsResponse, reader: jspb.BinaryReader): ListClientsResponse;
}

export namespace ListClientsResponse {
  export type AsObject = {
    clientsList: Array<Oauth2Client.AsObject>,
  }
}

export class DisableClientRequest extends jspb.Message {
  getClientid(): string;
  setClientid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableClientRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableClientRequest): DisableClientRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableClientRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableClientRequest;
  static deserializeBinaryFromReader(message: DisableClientRequest, reader: jspb.BinaryReader): DisableClientRequest;
}

export namespace DisableClientRequest {
  export type AsObject = {
    clientid: string,
  }
}

export class DisableClientResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableClientResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableClientResponse): DisableClientResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableClientResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableClientResponse;
  static deserializeBinaryFromReader(message: DisableClientResponse, reader: jspb.BinaryReader): DisableClientResponse;
}

export namespace DisableClientResponse {
  export type AsObject = {
  }
}

export interface ClientTypeMap {
  PRIVATE: 0;
  PUBLIC: 1;
}

export const ClientType: ClientTypeMap;

