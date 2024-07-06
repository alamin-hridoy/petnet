// package: user
// file: brank.as/rbac/gunk/v1/user/all.proto

import * as brank_as_rbac_gunk_v1_user_all_pb from "./all_pb";
import {grpc} from "@improbable-eng/grpc-web";

type SignupSignup = {
  readonly methodName: string;
  readonly service: typeof Signup;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.SignupRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.SignupResponse;
};

type SignupResendConfirmEmail = {
  readonly methodName: string;
  readonly service: typeof Signup;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.ResendConfirmEmailRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.ResendConfirmEmailResponse;
};

type SignupEmailConfirmation = {
  readonly methodName: string;
  readonly service: typeof Signup;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.EmailConfirmationRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.EmailConfirmationResponse;
};

type SignupForgotPassword = {
  readonly methodName: string;
  readonly service: typeof Signup;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.ForgotPasswordRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.ForgotPasswordResponse;
};

type SignupResetPassword = {
  readonly methodName: string;
  readonly service: typeof Signup;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.ResetPasswordRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.ResetPasswordResponse;
};

export class Signup {
  static readonly serviceName: string;
  static readonly Signup: SignupSignup;
  static readonly ResendConfirmEmail: SignupResendConfirmEmail;
  static readonly EmailConfirmation: SignupEmailConfirmation;
  static readonly ForgotPassword: SignupForgotPassword;
  static readonly ResetPassword: SignupResetPassword;
}

type UserServiceGetUser = {
  readonly methodName: string;
  readonly service: typeof UserService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.GetUserRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.GetUserResponse;
};

type UserServiceListUsers = {
  readonly methodName: string;
  readonly service: typeof UserService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.ListUsersRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.ListUsersResponse;
};

type UserServiceChangePassword = {
  readonly methodName: string;
  readonly service: typeof UserService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.ChangePasswordRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.ChangePasswordResponse;
};

type UserServiceConfirmUpdate = {
  readonly methodName: string;
  readonly service: typeof UserService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.ConfirmUpdateRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.ConfirmUpdateResponse;
};

type UserServiceUpdateUser = {
  readonly methodName: string;
  readonly service: typeof UserService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.UpdateUserRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.UpdateUserResponse;
};

type UserServiceDisableUser = {
  readonly methodName: string;
  readonly service: typeof UserService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.DisableUserRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.DisableUserResponse;
};

type UserServiceEnableUser = {
  readonly methodName: string;
  readonly service: typeof UserService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.EnableUserRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.EnableUserResponse;
};

export class UserService {
  static readonly serviceName: string;
  static readonly GetUser: UserServiceGetUser;
  static readonly ListUsers: UserServiceListUsers;
  static readonly ChangePassword: UserServiceChangePassword;
  static readonly ConfirmUpdate: UserServiceConfirmUpdate;
  static readonly UpdateUser: UserServiceUpdateUser;
  static readonly DisableUser: UserServiceDisableUser;
  static readonly EnableUser: UserServiceEnableUser;
}

type UserAuthServiceAuthenticateUser = {
  readonly methodName: string;
  readonly service: typeof UserAuthService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_user_all_pb.AuthenticateUserRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_user_all_pb.AuthenticateUserResponse;
};

export class UserAuthService {
  static readonly serviceName: string;
  static readonly AuthenticateUser: UserAuthServiceAuthenticateUser;
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

export class SignupClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  signup(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.SignupRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.SignupResponse|null) => void
  ): UnaryResponse;
  signup(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.SignupRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.SignupResponse|null) => void
  ): UnaryResponse;
  resendConfirmEmail(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.ResendConfirmEmailRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.ResendConfirmEmailResponse|null) => void
  ): UnaryResponse;
  resendConfirmEmail(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.ResendConfirmEmailRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.ResendConfirmEmailResponse|null) => void
  ): UnaryResponse;
  emailConfirmation(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.EmailConfirmationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.EmailConfirmationResponse|null) => void
  ): UnaryResponse;
  emailConfirmation(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.EmailConfirmationRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.EmailConfirmationResponse|null) => void
  ): UnaryResponse;
  forgotPassword(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.ForgotPasswordRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.ForgotPasswordResponse|null) => void
  ): UnaryResponse;
  forgotPassword(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.ForgotPasswordRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.ForgotPasswordResponse|null) => void
  ): UnaryResponse;
  resetPassword(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.ResetPasswordRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.ResetPasswordResponse|null) => void
  ): UnaryResponse;
  resetPassword(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.ResetPasswordRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.ResetPasswordResponse|null) => void
  ): UnaryResponse;
}

export class UserServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getUser(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.GetUserRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.GetUserResponse|null) => void
  ): UnaryResponse;
  getUser(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.GetUserRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.GetUserResponse|null) => void
  ): UnaryResponse;
  listUsers(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.ListUsersRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.ListUsersResponse|null) => void
  ): UnaryResponse;
  listUsers(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.ListUsersRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.ListUsersResponse|null) => void
  ): UnaryResponse;
  changePassword(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.ChangePasswordRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.ChangePasswordResponse|null) => void
  ): UnaryResponse;
  changePassword(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.ChangePasswordRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.ChangePasswordResponse|null) => void
  ): UnaryResponse;
  confirmUpdate(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.ConfirmUpdateRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.ConfirmUpdateResponse|null) => void
  ): UnaryResponse;
  confirmUpdate(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.ConfirmUpdateRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.ConfirmUpdateResponse|null) => void
  ): UnaryResponse;
  updateUser(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.UpdateUserRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.UpdateUserResponse|null) => void
  ): UnaryResponse;
  updateUser(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.UpdateUserRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.UpdateUserResponse|null) => void
  ): UnaryResponse;
  disableUser(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.DisableUserRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.DisableUserResponse|null) => void
  ): UnaryResponse;
  disableUser(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.DisableUserRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.DisableUserResponse|null) => void
  ): UnaryResponse;
  enableUser(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.EnableUserRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.EnableUserResponse|null) => void
  ): UnaryResponse;
  enableUser(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.EnableUserRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.EnableUserResponse|null) => void
  ): UnaryResponse;
}

export class UserAuthServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  authenticateUser(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.AuthenticateUserRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.AuthenticateUserResponse|null) => void
  ): UnaryResponse;
  authenticateUser(
    requestMessage: brank_as_rbac_gunk_v1_user_all_pb.AuthenticateUserRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_user_all_pb.AuthenticateUserResponse|null) => void
  ): UnaryResponse;
}

