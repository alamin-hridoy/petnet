// package: brankas.rbac.v1.invite
// file: brank.as/rbac/gunk/v1/invite/all.proto

import * as brank_as_rbac_gunk_v1_invite_all_pb from "./all_pb";
import {grpc} from "@improbable-eng/grpc-web";

type InviteServiceInviteUser = {
  readonly methodName: string;
  readonly service: typeof InviteService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_invite_all_pb.InviteUserRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_invite_all_pb.InviteUserResponse;
};

type InviteServiceResend = {
  readonly methodName: string;
  readonly service: typeof InviteService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_invite_all_pb.ResendRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_invite_all_pb.ResendResponse;
};

type InviteServiceListInvite = {
  readonly methodName: string;
  readonly service: typeof InviteService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_invite_all_pb.ListInviteRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_invite_all_pb.ListInviteResponse;
};

type InviteServiceRetrieveInvite = {
  readonly methodName: string;
  readonly service: typeof InviteService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_invite_all_pb.RetrieveInviteRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_invite_all_pb.RetrieveInviteResponse;
};

type InviteServiceCancelInvite = {
  readonly methodName: string;
  readonly service: typeof InviteService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_invite_all_pb.CancelInviteRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_invite_all_pb.CancelInviteResponse;
};

type InviteServiceApprove = {
  readonly methodName: string;
  readonly service: typeof InviteService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_invite_all_pb.ApproveRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_invite_all_pb.ApproveResponse;
};

type InviteServiceRevoke = {
  readonly methodName: string;
  readonly service: typeof InviteService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_invite_all_pb.RevokeRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_invite_all_pb.RevokeResponse;
};

export class InviteService {
  static readonly serviceName: string;
  static readonly InviteUser: InviteServiceInviteUser;
  static readonly Resend: InviteServiceResend;
  static readonly ListInvite: InviteServiceListInvite;
  static readonly RetrieveInvite: InviteServiceRetrieveInvite;
  static readonly CancelInvite: InviteServiceCancelInvite;
  static readonly Approve: InviteServiceApprove;
  static readonly Revoke: InviteServiceRevoke;
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

export class InviteServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  inviteUser(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.InviteUserRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.InviteUserResponse|null) => void
  ): UnaryResponse;
  inviteUser(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.InviteUserRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.InviteUserResponse|null) => void
  ): UnaryResponse;
  resend(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.ResendRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.ResendResponse|null) => void
  ): UnaryResponse;
  resend(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.ResendRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.ResendResponse|null) => void
  ): UnaryResponse;
  listInvite(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.ListInviteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.ListInviteResponse|null) => void
  ): UnaryResponse;
  listInvite(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.ListInviteRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.ListInviteResponse|null) => void
  ): UnaryResponse;
  retrieveInvite(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.RetrieveInviteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.RetrieveInviteResponse|null) => void
  ): UnaryResponse;
  retrieveInvite(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.RetrieveInviteRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.RetrieveInviteResponse|null) => void
  ): UnaryResponse;
  cancelInvite(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.CancelInviteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.CancelInviteResponse|null) => void
  ): UnaryResponse;
  cancelInvite(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.CancelInviteRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.CancelInviteResponse|null) => void
  ): UnaryResponse;
  approve(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.ApproveRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.ApproveResponse|null) => void
  ): UnaryResponse;
  approve(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.ApproveRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.ApproveResponse|null) => void
  ): UnaryResponse;
  revoke(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.RevokeRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.RevokeResponse|null) => void
  ): UnaryResponse;
  revoke(
    requestMessage: brank_as_rbac_gunk_v1_invite_all_pb.RevokeRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_invite_all_pb.RevokeResponse|null) => void
  ): UnaryResponse;
}

