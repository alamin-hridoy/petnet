// package: rbac.brankas.consent
// file: brank.as/rbac/gunk/v1/consent/all.proto

import * as brank_as_rbac_gunk_v1_consent_all_pb from "./all_pb";
import {grpc} from "@improbable-eng/grpc-web";

type GrantServiceServeGrant = {
  readonly methodName: string;
  readonly service: typeof GrantService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_consent_all_pb.ServeGrantRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_consent_all_pb.ServeGrantResponse;
};

type GrantServiceGrant = {
  readonly methodName: string;
  readonly service: typeof GrantService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_consent_all_pb.GrantRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_consent_all_pb.GrantResponse;
};

export class GrantService {
  static readonly serviceName: string;
  static readonly ServeGrant: GrantServiceServeGrant;
  static readonly Grant: GrantServiceGrant;
}

type ScopeServiceUpsertScope = {
  readonly methodName: string;
  readonly service: typeof ScopeService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_consent_all_pb.UpsertScopeRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_consent_all_pb.UpsertScopeResponse;
};

type ScopeServiceUpdateGroup = {
  readonly methodName: string;
  readonly service: typeof ScopeService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_consent_all_pb.UpdateGroupRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_consent_all_pb.UpdateGroupResponse;
};

type ScopeServiceGetScope = {
  readonly methodName: string;
  readonly service: typeof ScopeService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_consent_all_pb.GetScopeRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_consent_all_pb.GetScopeResponse;
};

export class ScopeService {
  static readonly serviceName: string;
  static readonly UpsertScope: ScopeServiceUpsertScope;
  static readonly UpdateGroup: ScopeServiceUpdateGroup;
  static readonly GetScope: ScopeServiceGetScope;
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

export class GrantServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  serveGrant(
    requestMessage: brank_as_rbac_gunk_v1_consent_all_pb.ServeGrantRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_consent_all_pb.ServeGrantResponse|null) => void
  ): UnaryResponse;
  serveGrant(
    requestMessage: brank_as_rbac_gunk_v1_consent_all_pb.ServeGrantRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_consent_all_pb.ServeGrantResponse|null) => void
  ): UnaryResponse;
  grant(
    requestMessage: brank_as_rbac_gunk_v1_consent_all_pb.GrantRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_consent_all_pb.GrantResponse|null) => void
  ): UnaryResponse;
  grant(
    requestMessage: brank_as_rbac_gunk_v1_consent_all_pb.GrantRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_consent_all_pb.GrantResponse|null) => void
  ): UnaryResponse;
}

export class ScopeServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  upsertScope(
    requestMessage: brank_as_rbac_gunk_v1_consent_all_pb.UpsertScopeRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_consent_all_pb.UpsertScopeResponse|null) => void
  ): UnaryResponse;
  upsertScope(
    requestMessage: brank_as_rbac_gunk_v1_consent_all_pb.UpsertScopeRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_consent_all_pb.UpsertScopeResponse|null) => void
  ): UnaryResponse;
  updateGroup(
    requestMessage: brank_as_rbac_gunk_v1_consent_all_pb.UpdateGroupRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_consent_all_pb.UpdateGroupResponse|null) => void
  ): UnaryResponse;
  updateGroup(
    requestMessage: brank_as_rbac_gunk_v1_consent_all_pb.UpdateGroupRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_consent_all_pb.UpdateGroupResponse|null) => void
  ): UnaryResponse;
  getScope(
    requestMessage: brank_as_rbac_gunk_v1_consent_all_pb.GetScopeRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_consent_all_pb.GetScopeResponse|null) => void
  ): UnaryResponse;
  getScope(
    requestMessage: brank_as_rbac_gunk_v1_consent_all_pb.GetScopeRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_consent_all_pb.GetScopeResponse|null) => void
  ): UnaryResponse;
}

