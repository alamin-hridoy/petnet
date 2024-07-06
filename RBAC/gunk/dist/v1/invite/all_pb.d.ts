// package: brankas.rbac.v1.invite
// file: brank.as/rbac/gunk/v1/invite/all.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class InviteUserRequest extends jspb.Message {
  getOrgid(): string;
  setOrgid(value: string): void;

  getOrgname(): string;
  setOrgname(value: string): void;

  getFirstname(): string;
  setFirstname(value: string): void;

  getLastname(): string;
  setLastname(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getPhone(): string;
  setPhone(value: string): void;

  getRole(): string;
  setRole(value: string): void;

  getCustomemaildataMap(): jspb.Map<string, string>;
  clearCustomemaildataMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InviteUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: InviteUserRequest): InviteUserRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: InviteUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InviteUserRequest;
  static deserializeBinaryFromReader(message: InviteUserRequest, reader: jspb.BinaryReader): InviteUserRequest;
}

export namespace InviteUserRequest {
  export type AsObject = {
    orgid: string,
    orgname: string,
    firstname: string,
    lastname: string,
    email: string,
    phone: string,
    role: string,
    customemaildataMap: Array<[string, string]>,
  }
}

export class InviteUserResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getInvitationcode(): string;
  setInvitationcode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InviteUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: InviteUserResponse): InviteUserResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: InviteUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InviteUserResponse;
  static deserializeBinaryFromReader(message: InviteUserResponse, reader: jspb.BinaryReader): InviteUserResponse;
}

export namespace InviteUserResponse {
  export type AsObject = {
    id: string,
    orgid: string,
    invitationcode: string,
  }
}

export class RetrieveInviteRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getCode(): string;
  setCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RetrieveInviteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RetrieveInviteRequest): RetrieveInviteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RetrieveInviteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RetrieveInviteRequest;
  static deserializeBinaryFromReader(message: RetrieveInviteRequest, reader: jspb.BinaryReader): RetrieveInviteRequest;
}

export namespace RetrieveInviteRequest {
  export type AsObject = {
    id: string,
    code: string,
  }
}

export class RetrieveInviteResponse extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getCountrycode(): string;
  setCountrycode(value: string): void;

  getPhone(): string;
  setPhone(value: string): void;

  getCompanyname(): string;
  setCompanyname(value: string): void;

  getActive(): boolean;
  setActive(value: boolean): void;

  getInviteemail(): string;
  setInviteemail(value: string): void;

  getInvitestatus(): string;
  setInvitestatus(value: string): void;

  hasInvited(): boolean;
  clearInvited(): void;
  getInvited(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setInvited(value?: google_protobuf_timestamp_pb.Timestamp): void;

  getFirstname(): string;
  setFirstname(value: string): void;

  getLastname(): string;
  setLastname(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RetrieveInviteResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RetrieveInviteResponse): RetrieveInviteResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RetrieveInviteResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RetrieveInviteResponse;
  static deserializeBinaryFromReader(message: RetrieveInviteResponse, reader: jspb.BinaryReader): RetrieveInviteResponse;
}

export namespace RetrieveInviteResponse {
  export type AsObject = {
    id: string,
    orgid: string,
    email: string,
    countrycode: string,
    phone: string,
    companyname: string,
    active: boolean,
    inviteemail: string,
    invitestatus: string,
    invited?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    firstname: string,
    lastname: string,
  }
}

export class ResendRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getCustomemaildataMap(): jspb.Map<string, string>;
  clearCustomemaildataMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResendRequest): ResendRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResendRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendRequest;
  static deserializeBinaryFromReader(message: ResendRequest, reader: jspb.BinaryReader): ResendRequest;
}

export namespace ResendRequest {
  export type AsObject = {
    id: string,
    customemaildataMap: Array<[string, string]>,
  }
}

export class ResendResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResendResponse): ResendResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResendResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendResponse;
  static deserializeBinaryFromReader(message: ResendResponse, reader: jspb.BinaryReader): ResendResponse;
}

export namespace ResendResponse {
  export type AsObject = {
  }
}

export class ListInviteRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListInviteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListInviteRequest): ListInviteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListInviteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListInviteRequest;
  static deserializeBinaryFromReader(message: ListInviteRequest, reader: jspb.BinaryReader): ListInviteRequest;
}

export namespace ListInviteRequest {
  export type AsObject = {
  }
}

export class Invite extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getContactemail(): string;
  setContactemail(value: string): void;

  hasInvited(): boolean;
  clearInvited(): void;
  getInvited(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setInvited(value?: google_protobuf_timestamp_pb.Timestamp): void;

  getRole(): string;
  setRole(value: string): void;

  getStatus(): string;
  setStatus(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Invite.AsObject;
  static toObject(includeInstance: boolean, msg: Invite): Invite.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Invite, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Invite;
  static deserializeBinaryFromReader(message: Invite, reader: jspb.BinaryReader): Invite;
}

export namespace Invite {
  export type AsObject = {
    id: string,
    orgid: string,
    contactemail: string,
    invited?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    role: string,
    status: string,
  }
}

export class ListInviteResponse extends jspb.Message {
  clearInvitesList(): void;
  getInvitesList(): Array<Invite>;
  setInvitesList(value: Array<Invite>): void;
  addInvites(value?: Invite, index?: number): Invite;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListInviteResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListInviteResponse): ListInviteResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListInviteResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListInviteResponse;
  static deserializeBinaryFromReader(message: ListInviteResponse, reader: jspb.BinaryReader): ListInviteResponse;
}

export namespace ListInviteResponse {
  export type AsObject = {
    invitesList: Array<Invite.AsObject>,
  }
}

export class CancelInviteRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CancelInviteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CancelInviteRequest): CancelInviteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CancelInviteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CancelInviteRequest;
  static deserializeBinaryFromReader(message: CancelInviteRequest, reader: jspb.BinaryReader): CancelInviteRequest;
}

export namespace CancelInviteRequest {
  export type AsObject = {
    id: string,
  }
}

export class CancelInviteResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CancelInviteResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CancelInviteResponse): CancelInviteResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CancelInviteResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CancelInviteResponse;
  static deserializeBinaryFromReader(message: CancelInviteResponse, reader: jspb.BinaryReader): CancelInviteResponse;
}

export namespace CancelInviteResponse {
  export type AsObject = {
  }
}

export class ApproveRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApproveRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ApproveRequest): ApproveRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApproveRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApproveRequest;
  static deserializeBinaryFromReader(message: ApproveRequest, reader: jspb.BinaryReader): ApproveRequest;
}

export namespace ApproveRequest {
  export type AsObject = {
    id: string,
  }
}

export class ApproveResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApproveResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ApproveResponse): ApproveResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ApproveResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApproveResponse;
  static deserializeBinaryFromReader(message: ApproveResponse, reader: jspb.BinaryReader): ApproveResponse;
}

export namespace ApproveResponse {
  export type AsObject = {
  }
}

export class RevokeRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RevokeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RevokeRequest): RevokeRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RevokeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RevokeRequest;
  static deserializeBinaryFromReader(message: RevokeRequest, reader: jspb.BinaryReader): RevokeRequest;
}

export namespace RevokeRequest {
  export type AsObject = {
    id: string,
  }
}

export class RevokeResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RevokeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RevokeResponse): RevokeResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RevokeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RevokeResponse;
  static deserializeBinaryFromReader(message: RevokeResponse, reader: jspb.BinaryReader): RevokeResponse;
}

export namespace RevokeResponse {
  export type AsObject = {
  }
}

