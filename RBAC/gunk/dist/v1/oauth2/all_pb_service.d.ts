// package: brankas.rbac.v1.oauth2
// file: brank.as/rbac/gunk/v1/oauth2/all.proto

import * as brank_as_rbac_gunk_v1_oauth2_all_pb from "./all_pb";
import {grpc} from "@improbable-eng/grpc-web";

type AuthClientServiceCreateClient = {
  readonly methodName: string;
  readonly service: typeof AuthClientService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_oauth2_all_pb.CreateClientRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_oauth2_all_pb.CreateClientResponse;
};

type AuthClientServiceUpdateClient = {
  readonly methodName: string;
  readonly service: typeof AuthClientService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_oauth2_all_pb.UpdateClientRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_oauth2_all_pb.UpdateClientResponse;
};

type AuthClientServiceListClients = {
  readonly methodName: string;
  readonly service: typeof AuthClientService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_oauth2_all_pb.ListClientsRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_oauth2_all_pb.ListClientsResponse;
};

type AuthClientServiceDisableClient = {
  readonly methodName: string;
  readonly service: typeof AuthClientService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_oauth2_all_pb.DisableClientRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_oauth2_all_pb.DisableClientResponse;
};

export class AuthClientService {
  static readonly serviceName: string;
  static readonly CreateClient: AuthClientServiceCreateClient;
  static readonly UpdateClient: AuthClientServiceUpdateClient;
  static readonly ListClients: AuthClientServiceListClients;
  static readonly DisableClient: AuthClientServiceDisableClient;
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

export class AuthClientServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  createClient(
    requestMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.CreateClientRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.CreateClientResponse|null) => void
  ): UnaryResponse;
  createClient(
    requestMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.CreateClientRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.CreateClientResponse|null) => void
  ): UnaryResponse;
  updateClient(
    requestMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.UpdateClientRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.UpdateClientResponse|null) => void
  ): UnaryResponse;
  updateClient(
    requestMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.UpdateClientRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.UpdateClientResponse|null) => void
  ): UnaryResponse;
  listClients(
    requestMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.ListClientsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.ListClientsResponse|null) => void
  ): UnaryResponse;
  listClients(
    requestMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.ListClientsRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.ListClientsResponse|null) => void
  ): UnaryResponse;
  disableClient(
    requestMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.DisableClientRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.DisableClientResponse|null) => void
  ): UnaryResponse;
  disableClient(
    requestMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.DisableClientRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_oauth2_all_pb.DisableClientResponse|null) => void
  ): UnaryResponse;
}

