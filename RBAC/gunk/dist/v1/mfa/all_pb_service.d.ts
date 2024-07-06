// package: mfa
// file: brank.as/rbac/gunk/v1/mfa/all.proto

import * as brank_as_rbac_gunk_v1_mfa_all_pb from "./all_pb";
import {grpc} from "@improbable-eng/grpc-web";

type MFAServiceGetRegisteredMFA = {
  readonly methodName: string;
  readonly service: typeof MFAService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.GetRegisteredMFARequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.GetRegisteredMFAResponse;
};

type MFAServiceEnableMFA = {
  readonly methodName: string;
  readonly service: typeof MFAService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.EnableMFARequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.EnableMFAResponse;
};

type MFAServiceDisableMFA = {
  readonly methodName: string;
  readonly service: typeof MFAService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.DisableMFARequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.DisableMFAResponse;
};

export class MFAService {
  static readonly serviceName: string;
  static readonly GetRegisteredMFA: MFAServiceGetRegisteredMFA;
  static readonly EnableMFA: MFAServiceEnableMFA;
  static readonly DisableMFA: MFAServiceDisableMFA;
}

type MFAAuthServiceInitiateMFA = {
  readonly methodName: string;
  readonly service: typeof MFAAuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.InitiateMFARequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.InitiateMFAResponse;
};

type MFAAuthServiceValidateMFA = {
  readonly methodName: string;
  readonly service: typeof MFAAuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.ValidateMFARequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.ValidateMFAResponse;
};

type MFAAuthServiceRetryMFA = {
  readonly methodName: string;
  readonly service: typeof MFAAuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.RetryMFARequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.RetryMFAResponse;
};

type MFAAuthServiceExternalMFA = {
  readonly methodName: string;
  readonly service: typeof MFAAuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.ExternalMFARequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_mfa_all_pb.ExternalMFAResponse;
};

export class MFAAuthService {
  static readonly serviceName: string;
  static readonly InitiateMFA: MFAAuthServiceInitiateMFA;
  static readonly ValidateMFA: MFAAuthServiceValidateMFA;
  static readonly RetryMFA: MFAAuthServiceRetryMFA;
  static readonly ExternalMFA: MFAAuthServiceExternalMFA;
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

export class MFAServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getRegisteredMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.GetRegisteredMFARequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.GetRegisteredMFAResponse|null) => void
  ): UnaryResponse;
  getRegisteredMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.GetRegisteredMFARequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.GetRegisteredMFAResponse|null) => void
  ): UnaryResponse;
  enableMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.EnableMFARequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.EnableMFAResponse|null) => void
  ): UnaryResponse;
  enableMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.EnableMFARequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.EnableMFAResponse|null) => void
  ): UnaryResponse;
  disableMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.DisableMFARequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.DisableMFAResponse|null) => void
  ): UnaryResponse;
  disableMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.DisableMFARequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.DisableMFAResponse|null) => void
  ): UnaryResponse;
}

export class MFAAuthServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  initiateMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.InitiateMFARequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.InitiateMFAResponse|null) => void
  ): UnaryResponse;
  initiateMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.InitiateMFARequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.InitiateMFAResponse|null) => void
  ): UnaryResponse;
  validateMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.ValidateMFARequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.ValidateMFAResponse|null) => void
  ): UnaryResponse;
  validateMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.ValidateMFARequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.ValidateMFAResponse|null) => void
  ): UnaryResponse;
  retryMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.RetryMFARequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.RetryMFAResponse|null) => void
  ): UnaryResponse;
  retryMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.RetryMFARequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.RetryMFAResponse|null) => void
  ): UnaryResponse;
  externalMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.ExternalMFARequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.ExternalMFAResponse|null) => void
  ): UnaryResponse;
  externalMFA(
    requestMessage: brank_as_rbac_gunk_v1_mfa_all_pb.ExternalMFARequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_mfa_all_pb.ExternalMFAResponse|null) => void
  ): UnaryResponse;
}

