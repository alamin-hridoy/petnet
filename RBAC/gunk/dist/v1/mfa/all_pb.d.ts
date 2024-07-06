// package: mfa
// file: brank.as/rbac/gunk/v1/mfa/all.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class GetRegisteredMFARequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRegisteredMFARequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRegisteredMFARequest): GetRegisteredMFARequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetRegisteredMFARequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRegisteredMFARequest;
  static deserializeBinaryFromReader(message: GetRegisteredMFARequest, reader: jspb.BinaryReader): GetRegisteredMFARequest;
}

export namespace GetRegisteredMFARequest {
  export type AsObject = {
    userid: string,
  }
}

export class GetRegisteredMFAResponse extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  clearMfaList(): void;
  getMfaList(): Array<MFAEntry>;
  setMfaList(value: Array<MFAEntry>): void;
  addMfa(value?: MFAEntry, index?: number): MFAEntry;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRegisteredMFAResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetRegisteredMFAResponse): GetRegisteredMFAResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetRegisteredMFAResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRegisteredMFAResponse;
  static deserializeBinaryFromReader(message: GetRegisteredMFAResponse, reader: jspb.BinaryReader): GetRegisteredMFAResponse;
}

export namespace GetRegisteredMFAResponse {
  export type AsObject = {
    userid: string,
    mfaList: Array<MFAEntry.AsObject>,
  }
}

export class MFAEntry extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getType(): MFAMap[keyof MFAMap];
  setType(value: MFAMap[keyof MFAMap]): void;

  getSource(): string;
  setSource(value: string): void;

  hasEnabled(): boolean;
  clearEnabled(): void;
  getEnabled(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setEnabled(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setDisabled(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MFAEntry.AsObject;
  static toObject(includeInstance: boolean, msg: MFAEntry): MFAEntry.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MFAEntry, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MFAEntry;
  static deserializeBinaryFromReader(message: MFAEntry, reader: jspb.BinaryReader): MFAEntry;
}

export namespace MFAEntry {
  export type AsObject = {
    id: string,
    type: MFAMap[keyof MFAMap],
    source: string,
    enabled?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    disabled?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class EnableMFARequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getType(): MFAMap[keyof MFAMap];
  setType(value: MFAMap[keyof MFAMap]): void;

  getSource(): string;
  setSource(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableMFARequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableMFARequest): EnableMFARequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableMFARequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableMFARequest;
  static deserializeBinaryFromReader(message: EnableMFARequest, reader: jspb.BinaryReader): EnableMFARequest;
}

export namespace EnableMFARequest {
  export type AsObject = {
    userid: string,
    type: MFAMap[keyof MFAMap],
    source: string,
  }
}

export class EnableMFAResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getInitializecode(): string;
  setInitializecode(value: string): void;

  getEventid(): string;
  setEventid(value: string): void;

  clearCodesList(): void;
  getCodesList(): Array<string>;
  setCodesList(value: Array<string>): void;
  addCodes(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableMFAResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableMFAResponse): EnableMFAResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableMFAResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableMFAResponse;
  static deserializeBinaryFromReader(message: EnableMFAResponse, reader: jspb.BinaryReader): EnableMFAResponse;
}

export namespace EnableMFAResponse {
  export type AsObject = {
    id: string,
    initializecode: string,
    eventid: string,
    codesList: Array<string>,
  }
}

export class DisableMFARequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getMfaid(): string;
  setMfaid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableMFARequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableMFARequest): DisableMFARequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableMFARequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableMFARequest;
  static deserializeBinaryFromReader(message: DisableMFARequest, reader: jspb.BinaryReader): DisableMFARequest;
}

export namespace DisableMFARequest {
  export type AsObject = {
    userid: string,
    mfaid: string,
  }
}

export class DisableMFAResponse extends jspb.Message {
  hasDisabled(): boolean;
  clearDisabled(): void;
  getDisabled(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setDisabled(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableMFAResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableMFAResponse): DisableMFAResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableMFAResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableMFAResponse;
  static deserializeBinaryFromReader(message: DisableMFAResponse, reader: jspb.BinaryReader): DisableMFAResponse;
}

export namespace DisableMFAResponse {
  export type AsObject = {
    disabled?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class InitiateMFARequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getType(): MFAMap[keyof MFAMap];
  setType(value: MFAMap[keyof MFAMap]): void;

  getSourceid(): string;
  setSourceid(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitiateMFARequest.AsObject;
  static toObject(includeInstance: boolean, msg: InitiateMFARequest): InitiateMFARequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: InitiateMFARequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitiateMFARequest;
  static deserializeBinaryFromReader(message: InitiateMFARequest, reader: jspb.BinaryReader): InitiateMFARequest;
}

export namespace InitiateMFARequest {
  export type AsObject = {
    userid: string,
    type: MFAMap[keyof MFAMap],
    sourceid: string,
    description: string,
  }
}

export class InitiateMFAResponse extends jspb.Message {
  getEventid(): string;
  setEventid(value: string): void;

  clearSourcesList(): void;
  getSourcesList(): Array<MFAEntry>;
  setSourcesList(value: Array<MFAEntry>): void;
  addSources(value?: MFAEntry, index?: number): MFAEntry;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitiateMFAResponse.AsObject;
  static toObject(includeInstance: boolean, msg: InitiateMFAResponse): InitiateMFAResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: InitiateMFAResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitiateMFAResponse;
  static deserializeBinaryFromReader(message: InitiateMFAResponse, reader: jspb.BinaryReader): InitiateMFAResponse;
}

export namespace InitiateMFAResponse {
  export type AsObject = {
    eventid: string,
    sourcesList: Array<MFAEntry.AsObject>,
    value: string,
  }
}

export class RetryMFARequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getEventid(): string;
  setEventid(value: string): void;

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
    userid: string,
    eventid: string,
  }
}

export class RetryMFAResponse extends jspb.Message {
  getEventid(): string;
  setEventid(value: string): void;

  clearSourcesList(): void;
  getSourcesList(): Array<MFAEntry>;
  setSourcesList(value: Array<MFAEntry>): void;
  addSources(value?: MFAEntry, index?: number): MFAEntry;

  getValue(): string;
  setValue(value: string): void;

  getAttempt(): number;
  setAttempt(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RetryMFAResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RetryMFAResponse): RetryMFAResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RetryMFAResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RetryMFAResponse;
  static deserializeBinaryFromReader(message: RetryMFAResponse, reader: jspb.BinaryReader): RetryMFAResponse;
}

export namespace RetryMFAResponse {
  export type AsObject = {
    eventid: string,
    sourcesList: Array<MFAEntry.AsObject>,
    value: string,
    attempt: number,
  }
}

export class ExternalMFARequest extends jspb.Message {
  getEventid(): string;
  setEventid(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  getSourceid(): string;
  setSourceid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExternalMFARequest.AsObject;
  static toObject(includeInstance: boolean, msg: ExternalMFARequest): ExternalMFARequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExternalMFARequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExternalMFARequest;
  static deserializeBinaryFromReader(message: ExternalMFARequest, reader: jspb.BinaryReader): ExternalMFARequest;
}

export namespace ExternalMFARequest {
  export type AsObject = {
    eventid: string,
    value: string,
    sourceid: string,
  }
}

export class ExternalMFAResponse extends jspb.Message {
  getEventid(): string;
  setEventid(value: string): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExternalMFAResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ExternalMFAResponse): ExternalMFAResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExternalMFAResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExternalMFAResponse;
  static deserializeBinaryFromReader(message: ExternalMFAResponse, reader: jspb.BinaryReader): ExternalMFAResponse;
}

export namespace ExternalMFAResponse {
  export type AsObject = {
    eventid: string,
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ValidateMFARequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getType(): MFAMap[keyof MFAMap];
  setType(value: MFAMap[keyof MFAMap]): void;

  getToken(): string;
  setToken(value: string): void;

  getEventid(): string;
  setEventid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateMFARequest.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateMFARequest): ValidateMFARequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValidateMFARequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateMFARequest;
  static deserializeBinaryFromReader(message: ValidateMFARequest, reader: jspb.BinaryReader): ValidateMFARequest;
}

export namespace ValidateMFARequest {
  export type AsObject = {
    userid: string,
    type: MFAMap[keyof MFAMap],
    token: string,
    eventid: string,
  }
}

export class ValidateMFAResponse extends jspb.Message {
  getEventid(): string;
  setEventid(value: string): void;

  hasValidated(): boolean;
  clearValidated(): void;
  getValidated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setValidated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  getExternalid(): string;
  setExternalid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateMFAResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateMFAResponse): ValidateMFAResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ValidateMFAResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateMFAResponse;
  static deserializeBinaryFromReader(message: ValidateMFAResponse, reader: jspb.BinaryReader): ValidateMFAResponse;
}

export namespace ValidateMFAResponse {
  export type AsObject = {
    eventid: string,
    validated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    externalid: string,
  }
}

export interface MFAMap {
  PASS: 0;
  TOTP: 1;
  CODE: 2;
  SMS: 3;
  RECOVERY: 4;
  EMAIL: 5;
}

export const MFA: MFAMap;

