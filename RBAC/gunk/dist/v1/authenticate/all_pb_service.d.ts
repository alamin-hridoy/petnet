// package: authenticate
// file: brank.as/rbac/gunk/v1/authenticate/all.proto

import * as brank_as_rbac_gunk_v1_authenticate_all_pb from "./all_pb";
import {grpc} from "@improbable-eng/grpc-web";

type SessionServiceLogin = {
  readonly methodName: string;
  readonly service: typeof SessionService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_authenticate_all_pb.LoginRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_authenticate_all_pb.Session;
};

type SessionServiceGetSession = {
  readonly methodName: string;
  readonly service: typeof SessionService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_authenticate_all_pb.GetSessionRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_authenticate_all_pb.Session;
};

type SessionServiceRetryMFA = {
  readonly methodName: string;
  readonly service: typeof SessionService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_authenticate_all_pb.RetryMFARequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_authenticate_all_pb.Session;
};

export class SessionService {
  static readonly serviceName: string;
  static readonly Login: SessionServiceLogin;
  static readonly GetSession: SessionServiceGetSession;
  static readonly RetryMFA: SessionServiceRetryMFA;
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

export class SessionServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  login(
    requestMessage: brank_as_rbac_gunk_v1_authenticate_all_pb.LoginRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_authenticate_all_pb.Session|null) => void
  ): UnaryResponse;
  login(
    requestMessage: brank_as_rbac_gunk_v1_authenticate_all_pb.LoginRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_authenticate_all_pb.Session|null) => void
  ): UnaryResponse;
  getSession(
    requestMessage: brank_as_rbac_gunk_v1_authenticate_all_pb.GetSessionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_authenticate_all_pb.Session|null) => void
  ): UnaryResponse;
  getSession(
    requestMessage: brank_as_rbac_gunk_v1_authenticate_all_pb.GetSessionRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_authenticate_all_pb.Session|null) => void
  ): UnaryResponse;
  retryMFA(
    requestMessage: brank_as_rbac_gunk_v1_authenticate_all_pb.RetryMFARequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_authenticate_all_pb.Session|null) => void
  ): UnaryResponse;
  retryMFA(
    requestMessage: brank_as_rbac_gunk_v1_authenticate_all_pb.RetryMFARequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_authenticate_all_pb.Session|null) => void
  ): UnaryResponse;
}

