// package: brankas.rbac.v1.serviceaccount
// file: brank.as/rbac/gunk/v1/serviceaccount/all.proto

import * as brank_as_rbac_gunk_v1_serviceaccount_all_pb from "./all_pb";
import {grpc} from "@improbable-eng/grpc-web";

type SvcAccountServiceCreateAccount = {
  readonly methodName: string;
  readonly service: typeof SvcAccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_serviceaccount_all_pb.CreateAccountRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_serviceaccount_all_pb.CreateAccountResponse;
};

type SvcAccountServiceListAccounts = {
  readonly methodName: string;
  readonly service: typeof SvcAccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_serviceaccount_all_pb.ListAccountsRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_serviceaccount_all_pb.ListAccountsResponse;
};

type SvcAccountServiceDisableAccount = {
  readonly methodName: string;
  readonly service: typeof SvcAccountService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_serviceaccount_all_pb.DisableAccountRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_serviceaccount_all_pb.DisableAccountResponse;
};

export class SvcAccountService {
  static readonly serviceName: string;
  static readonly CreateAccount: SvcAccountServiceCreateAccount;
  static readonly ListAccounts: SvcAccountServiceListAccounts;
  static readonly DisableAccount: SvcAccountServiceDisableAccount;
}

type ValidationServiceValidateAccount = {
  readonly methodName: string;
  readonly service: typeof ValidationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_serviceaccount_all_pb.ValidateAccountRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_serviceaccount_all_pb.ValidateAccountResponse;
};

export class ValidationService {
  static readonly serviceName: string;
  static readonly ValidateAccount: ValidationServiceValidateAccount;
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

export class SvcAccountServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  createAccount(
    requestMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.CreateAccountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.CreateAccountResponse|null) => void
  ): UnaryResponse;
  createAccount(
    requestMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.CreateAccountRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.CreateAccountResponse|null) => void
  ): UnaryResponse;
  listAccounts(
    requestMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.ListAccountsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.ListAccountsResponse|null) => void
  ): UnaryResponse;
  listAccounts(
    requestMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.ListAccountsRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.ListAccountsResponse|null) => void
  ): UnaryResponse;
  disableAccount(
    requestMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.DisableAccountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.DisableAccountResponse|null) => void
  ): UnaryResponse;
  disableAccount(
    requestMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.DisableAccountRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.DisableAccountResponse|null) => void
  ): UnaryResponse;
}

export class ValidationServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  validateAccount(
    requestMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.ValidateAccountRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.ValidateAccountResponse|null) => void
  ): UnaryResponse;
  validateAccount(
    requestMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.ValidateAccountRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_serviceaccount_all_pb.ValidateAccountResponse|null) => void
  ): UnaryResponse;
}

