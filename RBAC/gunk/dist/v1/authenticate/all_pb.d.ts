// package: authenticate
// file: brank.as/rbac/gunk/v1/authenticate/all.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import * as brank_as_rbac_gunk_v1_mfa_all_pb from "./all_pb";

export class LoginRequest extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  getClientid(): string;
  setClientid(value: string): void;

  getExtraMap(): jspb.Map<string, string>;
  clearExtraMap(): void;
  getMfaeventid(): string;
  setMfaeventid(value: string): void;

  getMfatype(): brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap];
  setMfatype(value: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap]): void;

  getMfatoken(): string;
  setMfatoken(value: string): void;

  getSubject(): string;
  setSubject(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LoginRequest): LoginRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LoginRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginRequest;
  static deserializeBinaryFromReader(message: LoginRequest, reader: jspb.BinaryReader): LoginRequest;
}

export namespace LoginRequest {
  export type AsObject = {
    username: string,
    password: string,
    clientid: string,
    extraMap: Array<[string, string]>,
    mfaeventid: string,
    mfatype: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap],
    mfatoken: string,
    subject: string,
  }
}

export class Session extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getSessionMap(): jspb.Map<string, string>;
  clearSessionMap(): void;
  getOpenidMap(): jspb.Map<string, string>;
  clearOpenidMap(): void;
  getMfaeventid(): string;
  setMfaeventid(value: string): void;

  getMfatype(): brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap];
  setMfatype(value: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap]): void;

  hasPasswordexpiry(): boolean;
  clearPasswordexpiry(): void;
  getPasswordexpiry(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setPasswordexpiry(value?: google_protobuf_timestamp_pb.Timestamp): void;

  getResetrequired(): boolean;
  setResetrequired(value: boolean): void;

  getAttempt(): number;
  setAttempt(value: number): void;

  getForcelogin(): boolean;
  setForcelogin(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Session.AsObject;
  static toObject(includeInstance: boolean, msg: Session): Session.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Session, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Session;
  static deserializeBinaryFromReader(message: Session, reader: jspb.BinaryReader): Session;
}

export namespace Session {
  export type AsObject = {
    userid: string,
    orgid: string,
    sessionMap: Array<[string, string]>,
    openidMap: Array<[string, string]>,
    mfaeventid: string,
    mfatype: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap],
    passwordexpiry?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    resetrequired: boolean,
    attempt: number,
    forcelogin: boolean,
  }
}

export class GetSessionRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getClientid(): string;
  setClientid(value: string): void;

  getExtraMap(): jspb.Map<string, string>;
  clearExtraMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSessionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSessionRequest): GetSessionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetSessionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSessionRequest;
  static deserializeBinaryFromReader(message: GetSessionRequest, reader: jspb.BinaryReader): GetSessionRequest;
}

export namespace GetSessionRequest {
  export type AsObject = {
    userid: string,
    clientid: string,
    extraMap: Array<[string, string]>,
  }
}

export class RetryMFARequest extends jspb.Message {
  getSubject(): string;
  setSubject(value: string): void;

  getClientid(): string;
  setClientid(value: string): void;

  getMfaeventid(): string;
  setMfaeventid(value: string): void;

  getMfatype(): brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap];
  setMfatype(value: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RetryMFARequest.AsObject;
  static toObject(includeInstance: boolean, msg: RetryMFARequest): RetryMFARequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RetryMFARequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RetryMFARequest;
  static deserializeBinaryFromReader(message: RetryMFARequest, reader: jspb.BinaryReader): RetryMFARequest;
}

export namespace RetryMFARequest {
  export type AsObject = {
    subject: string,
    clientid: string,
    mfaeventid: string,
    mfatype: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap],
  }
}

export class SessionError extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): void;

  getTrackingattempts(): boolean;
  setTrackingattempts(value: boolean): void;

  getRemainingattempts(): number;
  setRemainingattempts(value: number): void;

  getErrordetailsMap(): jspb.Map<string, string>;
  clearErrordetailsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SessionError.AsObject;
  static toObject(includeInstance: boolean, msg: SessionError): SessionError.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SessionError, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SessionError;
  static deserializeBinaryFromReader(message: SessionError, reader: jspb.BinaryReader): SessionError;
}

export namespace SessionError {
  export type AsObject = {
    message: string,
    trackingattempts: boolean,
    remainingattempts: number,
    errordetailsMap: Array<[string, string]>,
  }
}

