// package: user
// file: brank.as/rbac/gunk/v1/user/all.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import * as brank_as_rbac_gunk_v1_mfa_all_pb from "./all_pb";

export class SignupRequest extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): void;

  getFirstname(): string;
  setFirstname(value: string): void;

  getLastname(): string;
  setLastname(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  getInvitecode(): string;
  setInvitecode(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SignupRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SignupRequest): SignupRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SignupRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SignupRequest;
  static deserializeBinaryFromReader(message: SignupRequest, reader: jspb.BinaryReader): SignupRequest;
}

export namespace SignupRequest {
  export type AsObject = {
    username: string,
    firstname: string,
    lastname: string,
    email: string,
    password: string,
    invitecode: string,
    orgid: string,
  }
}

export class SignupResponse extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SignupResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SignupResponse): SignupResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SignupResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SignupResponse;
  static deserializeBinaryFromReader(message: SignupResponse, reader: jspb.BinaryReader): SignupResponse;
}

export namespace SignupResponse {
  export type AsObject = {
    userid: string,
    orgid: string,
  }
}

export class ResendConfirmEmailRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendConfirmEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResendConfirmEmailRequest): ResendConfirmEmailRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResendConfirmEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendConfirmEmailRequest;
  static deserializeBinaryFromReader(message: ResendConfirmEmailRequest, reader: jspb.BinaryReader): ResendConfirmEmailRequest;
}

export namespace ResendConfirmEmailRequest {
  export type AsObject = {
    userid: string,
  }
}

export class ResendConfirmEmailResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendConfirmEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResendConfirmEmailResponse): ResendConfirmEmailResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResendConfirmEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendConfirmEmailResponse;
  static deserializeBinaryFromReader(message: ResendConfirmEmailResponse, reader: jspb.BinaryReader): ResendConfirmEmailResponse;
}

export namespace ResendConfirmEmailResponse {
  export type AsObject = {
  }
}

export class EmailConfirmationRequest extends jspb.Message {
  getCode(): string;
  setCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EmailConfirmationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EmailConfirmationRequest): EmailConfirmationRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EmailConfirmationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EmailConfirmationRequest;
  static deserializeBinaryFromReader(message: EmailConfirmationRequest, reader: jspb.BinaryReader): EmailConfirmationRequest;
}

export namespace EmailConfirmationRequest {
  export type AsObject = {
    code: string,
  }
}

export class EmailConfirmationResponse extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getUserid(): string;
  setUserid(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  getUsername(): string;
  setUsername(value: string): void;

  getFirstname(): string;
  setFirstname(value: string): void;

  getLastname(): string;
  setLastname(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EmailConfirmationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EmailConfirmationResponse): EmailConfirmationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EmailConfirmationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EmailConfirmationResponse;
  static deserializeBinaryFromReader(message: EmailConfirmationResponse, reader: jspb.BinaryReader): EmailConfirmationResponse;
}

export namespace EmailConfirmationResponse {
  export type AsObject = {
    email: string,
    userid: string,
    orgid: string,
    username: string,
    firstname: string,
    lastname: string,
  }
}

export class ForgotPasswordRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ForgotPasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ForgotPasswordRequest): ForgotPasswordRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ForgotPasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ForgotPasswordRequest;
  static deserializeBinaryFromReader(message: ForgotPasswordRequest, reader: jspb.BinaryReader): ForgotPasswordRequest;
}

export namespace ForgotPasswordRequest {
  export type AsObject = {
    email: string,
  }
}

export class ForgotPasswordResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ForgotPasswordResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ForgotPasswordResponse): ForgotPasswordResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ForgotPasswordResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ForgotPasswordResponse;
  static deserializeBinaryFromReader(message: ForgotPasswordResponse, reader: jspb.BinaryReader): ForgotPasswordResponse;
}

export namespace ForgotPasswordResponse {
  export type AsObject = {
  }
}

export class ResetPasswordRequest extends jspb.Message {
  getCode(): string;
  setCode(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetPasswordRequest): ResetPasswordRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResetPasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetPasswordRequest;
  static deserializeBinaryFromReader(message: ResetPasswordRequest, reader: jspb.BinaryReader): ResetPasswordRequest;
}

export namespace ResetPasswordRequest {
  export type AsObject = {
    code: string,
    password: string,
  }
}

export class ResetPasswordResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPasswordResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetPasswordResponse): ResetPasswordResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResetPasswordResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetPasswordResponse;
  static deserializeBinaryFromReader(message: ResetPasswordResponse, reader: jspb.BinaryReader): ResetPasswordResponse;
}

export namespace ResetPasswordResponse {
  export type AsObject = {
  }
}

export class GetUserRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserRequest): GetUserRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserRequest;
  static deserializeBinaryFromReader(message: GetUserRequest, reader: jspb.BinaryReader): GetUserRequest;
}

export namespace GetUserRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetUserResponse extends jspb.Message {
  hasUser(): boolean;
  clearUser(): void;
  getUser(): User | undefined;
  setUser(value?: User): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserResponse): GetUserResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserResponse;
  static deserializeBinaryFromReader(message: GetUserResponse, reader: jspb.BinaryReader): GetUserResponse;
}

export namespace GetUserResponse {
  export type AsObject = {
    user?: User.AsObject,
  }
}

export class User extends jspb.Message {
  getId(): string;
  setId(value: string): void;

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

  getInvitestatus(): string;
  setInvitestatus(value: string): void;

  getCountrycode(): string;
  setCountrycode(value: string): void;

  getPhone(): string;
  setPhone(value: string): void;

  hasCreated(): boolean;
  clearCreated(): void;
  getCreated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasDeleted(): boolean;
  clearDeleted(): void;
  getDeleted(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setDeleted(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    id: string,
    orgid: string,
    orgname: string,
    firstname: string,
    lastname: string,
    email: string,
    invitestatus: string,
    countrycode: string,
    phone: string,
    created?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    deleted?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ListUsersRequest extends jspb.Message {
  getOrgid(): string;
  setOrgid(value: string): void;

  getName(): string;
  setName(value: string): void;

  getSortby(): SortByMap[keyof SortByMap];
  setSortby(value: SortByMap[keyof SortByMap]): void;

  getSortbycolumn(): SortByColumnMap[keyof SortByColumnMap];
  setSortbycolumn(value: SortByColumnMap[keyof SortByColumnMap]): void;

  clearStatusList(): void;
  getStatusList(): Array<StatusMap[keyof StatusMap]>;
  setStatusList(value: Array<StatusMap[keyof StatusMap]>): void;
  addStatus(value: StatusMap[keyof StatusMap], index?: number): StatusMap[keyof StatusMap];

  getLimit(): number;
  setLimit(value: number): void;

  getOffset(): number;
  setOffset(value: number): void;

  clearIdList(): void;
  getIdList(): Array<string>;
  setIdList(value: Array<string>): void;
  addId(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUsersRequest): ListUsersRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUsersRequest;
  static deserializeBinaryFromReader(message: ListUsersRequest, reader: jspb.BinaryReader): ListUsersRequest;
}

export namespace ListUsersRequest {
  export type AsObject = {
    orgid: string,
    name: string,
    sortby: SortByMap[keyof SortByMap],
    sortbycolumn: SortByColumnMap[keyof SortByColumnMap],
    statusList: Array<StatusMap[keyof StatusMap]>,
    limit: number,
    offset: number,
    idList: Array<string>,
  }
}

export class ListUsersResponse extends jspb.Message {
  clearUsersList(): void;
  getUsersList(): Array<User>;
  setUsersList(value: Array<User>): void;
  addUsers(value?: User, index?: number): User;

  getTotal(): number;
  setTotal(value: number): void;

  getUserMap(): jspb.Map<string, User>;
  clearUserMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUsersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUsersResponse): ListUsersResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListUsersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUsersResponse;
  static deserializeBinaryFromReader(message: ListUsersResponse, reader: jspb.BinaryReader): ListUsersResponse;
}

export namespace ListUsersResponse {
  export type AsObject = {
    usersList: Array<User.AsObject>,
    total: number,
    userMap: Array<[string, User.AsObject]>,
  }
}

export class ChangePasswordRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getEventid(): string;
  setEventid(value: string): void;

  getOldpassword(): string;
  setOldpassword(value: string): void;

  getNewpassword(): string;
  setNewpassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangePasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChangePasswordRequest): ChangePasswordRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangePasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangePasswordRequest;
  static deserializeBinaryFromReader(message: ChangePasswordRequest, reader: jspb.BinaryReader): ChangePasswordRequest;
}

export namespace ChangePasswordRequest {
  export type AsObject = {
    userid: string,
    eventid: string,
    oldpassword: string,
    newpassword: string,
  }
}

export class ChangePasswordResponse extends jspb.Message {
  getMfaeventid(): string;
  setMfaeventid(value: string): void;

  getMfatype(): brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap];
  setMfatype(value: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap]): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangePasswordResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ChangePasswordResponse): ChangePasswordResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ChangePasswordResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangePasswordResponse;
  static deserializeBinaryFromReader(message: ChangePasswordResponse, reader: jspb.BinaryReader): ChangePasswordResponse;
}

export namespace ChangePasswordResponse {
  export type AsObject = {
    mfaeventid: string,
    mfatype: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap],
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ConfirmUpdateRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

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
    userid: string,
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

export class DisableUserRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getCustomemaildataMap(): jspb.Map<string, string>;
  clearCustomemaildataMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableUserRequest): DisableUserRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableUserRequest;
  static deserializeBinaryFromReader(message: DisableUserRequest, reader: jspb.BinaryReader): DisableUserRequest;
}

export namespace DisableUserRequest {
  export type AsObject = {
    userid: string,
    customemaildataMap: Array<[string, string]>,
  }
}

export class DisableUserResponse extends jspb.Message {
  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableUserResponse): DisableUserResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DisableUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableUserResponse;
  static deserializeBinaryFromReader(message: DisableUserResponse, reader: jspb.BinaryReader): DisableUserResponse;
}

export namespace DisableUserResponse {
  export type AsObject = {
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class EnableUserRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getCustomemaildataMap(): jspb.Map<string, string>;
  clearCustomemaildataMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableUserRequest): EnableUserRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableUserRequest;
  static deserializeBinaryFromReader(message: EnableUserRequest, reader: jspb.BinaryReader): EnableUserRequest;
}

export namespace EnableUserRequest {
  export type AsObject = {
    userid: string,
    customemaildataMap: Array<[string, string]>,
  }
}

export class EnableUserResponse extends jspb.Message {
  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableUserResponse): EnableUserResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnableUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableUserResponse;
  static deserializeBinaryFromReader(message: EnableUserResponse, reader: jspb.BinaryReader): EnableUserResponse;
}

export namespace EnableUserResponse {
  export type AsObject = {
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class UpdateUserRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getFirstname(): string;
  setFirstname(value: string): void;

  getLastname(): string;
  setLastname(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getMfatype(): brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap];
  setMfatype(value: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap]): void;

  getLoginmfa(): EnableOptMap[keyof EnableOptMap];
  setLoginmfa(value: EnableOptMap[keyof EnableOptMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserRequest): UpdateUserRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserRequest;
  static deserializeBinaryFromReader(message: UpdateUserRequest, reader: jspb.BinaryReader): UpdateUserRequest;
}

export namespace UpdateUserRequest {
  export type AsObject = {
    userid: string,
    firstname: string,
    lastname: string,
    email: string,
    mfatype: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap],
    loginmfa: EnableOptMap[keyof EnableOptMap],
  }
}

export class UpdateUserResponse extends jspb.Message {
  getMfaeventid(): string;
  setMfaeventid(value: string): void;

  getMfatype(): brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap];
  setMfatype(value: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap]): void;

  hasUpdated(): boolean;
  clearUpdated(): void;
  getUpdated(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdated(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserResponse): UpdateUserResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpdateUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserResponse;
  static deserializeBinaryFromReader(message: UpdateUserResponse, reader: jspb.BinaryReader): UpdateUserResponse;
}

export namespace UpdateUserResponse {
  export type AsObject = {
    mfaeventid: string,
    mfatype: brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap[keyof brank_as_rbac_gunk_v1_mfa_all_pb.MFAMap],
    updated?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class AuthenticateUserRequest extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthenticateUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AuthenticateUserRequest): AuthenticateUserRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuthenticateUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthenticateUserRequest;
  static deserializeBinaryFromReader(message: AuthenticateUserRequest, reader: jspb.BinaryReader): AuthenticateUserRequest;
}

export namespace AuthenticateUserRequest {
  export type AsObject = {
    username: string,
    password: string,
  }
}

export class AuthenticateUserResponse extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getOrgid(): string;
  setOrgid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthenticateUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AuthenticateUserResponse): AuthenticateUserResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuthenticateUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthenticateUserResponse;
  static deserializeBinaryFromReader(message: AuthenticateUserResponse, reader: jspb.BinaryReader): AuthenticateUserResponse;
}

export namespace AuthenticateUserResponse {
  export type AsObject = {
    userid: string,
    orgid: string,
  }
}

export interface SortByMap {
  DESC: 0;
  ASC: 1;
}

export const SortBy: SortByMap;

export interface SortByColumnMap {
  CREATEDDATE: 0;
  USERNAME: 1;
}

export const SortByColumn: SortByColumnMap;

export interface StatusMap {
  INVITED: 0;
  INVITESENT: 1;
  EXPIRED: 2;
  INPROGRESS: 3;
  REVOKED: 4;
  APPROVED: 5;
}

export const Status: StatusMap;

export interface EnableOptMap {
  NOCHANGE: 0;
  ENABLE: 1;
  DISABLE: 2;
}

export const EnableOpt: EnableOptMap;

