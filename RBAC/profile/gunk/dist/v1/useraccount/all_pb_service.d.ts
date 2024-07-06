// package: useraccount
// file: brank.as/rbac/profile/gunk/v1/useraccount/all.proto

import * as brank_as_rbac_profile_gunk_v1_useraccount_all_pb from "./all_pb";
import {grpc} from "@improbable-eng/grpc-web";

type UserServiceGetUser = {
  readonly methodName: string;
  readonly service: typeof UserService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_profile_gunk_v1_useraccount_all_pb.GetUserRequest;
  readonly responseType: typeof brank_as_rbac_profile_gunk_v1_useraccount_all_pb.GetUserResponse;
};

type UserServiceListUsers = {
  readonly methodName: string;
  readonly service: typeof UserService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_profile_gunk_v1_useraccount_all_pb.ListUsersRequest;
  readonly responseType: typeof brank_as_rbac_profile_gunk_v1_useraccount_all_pb.ListUsersResponse;
};

export class UserService {
  static readonly serviceName: string;
  static readonly GetUser: UserServiceGetUser;
  static readonly ListUsers: UserServiceListUsers;
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

export class UserServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getUser(
    requestMessage: brank_as_rbac_profile_gunk_v1_useraccount_all_pb.GetUserRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_profile_gunk_v1_useraccount_all_pb.GetUserResponse|null) => void
  ): UnaryResponse;
  getUser(
    requestMessage: brank_as_rbac_profile_gunk_v1_useraccount_all_pb.GetUserRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_profile_gunk_v1_useraccount_all_pb.GetUserResponse|null) => void
  ): UnaryResponse;
  listUsers(
    requestMessage: brank_as_rbac_profile_gunk_v1_useraccount_all_pb.ListUsersRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_profile_gunk_v1_useraccount_all_pb.ListUsersResponse|null) => void
  ): UnaryResponse;
  listUsers(
    requestMessage: brank_as_rbac_profile_gunk_v1_useraccount_all_pb.ListUsersRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_profile_gunk_v1_useraccount_all_pb.ListUsersResponse|null) => void
  ): UnaryResponse;
}

