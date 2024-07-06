// package: organization
// file: brank.as/rbac/gunk/v1/organization/all.proto

import * as brank_as_rbac_gunk_v1_organization_all_pb from "./all_pb";
import {grpc} from "@improbable-eng/grpc-web";

type OrganizationServiceGetOrganization = {
  readonly methodName: string;
  readonly service: typeof OrganizationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_organization_all_pb.GetOrganizationRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_organization_all_pb.GetOrganizationResponse;
};

type OrganizationServiceUpdateOrganization = {
  readonly methodName: string;
  readonly service: typeof OrganizationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_organization_all_pb.UpdateOrganizationRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_organization_all_pb.UpdateOrganizationResponse;
};

type OrganizationServiceConfirmUpdate = {
  readonly methodName: string;
  readonly service: typeof OrganizationService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof brank_as_rbac_gunk_v1_organization_all_pb.ConfirmUpdateRequest;
  readonly responseType: typeof brank_as_rbac_gunk_v1_organization_all_pb.ConfirmUpdateResponse;
};

export class OrganizationService {
  static readonly serviceName: string;
  static readonly GetOrganization: OrganizationServiceGetOrganization;
  static readonly UpdateOrganization: OrganizationServiceUpdateOrganization;
  static readonly ConfirmUpdate: OrganizationServiceConfirmUpdate;
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

export class OrganizationServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getOrganization(
    requestMessage: brank_as_rbac_gunk_v1_organization_all_pb.GetOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_organization_all_pb.GetOrganizationResponse|null) => void
  ): UnaryResponse;
  getOrganization(
    requestMessage: brank_as_rbac_gunk_v1_organization_all_pb.GetOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_organization_all_pb.GetOrganizationResponse|null) => void
  ): UnaryResponse;
  updateOrganization(
    requestMessage: brank_as_rbac_gunk_v1_organization_all_pb.UpdateOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_organization_all_pb.UpdateOrganizationResponse|null) => void
  ): UnaryResponse;
  updateOrganization(
    requestMessage: brank_as_rbac_gunk_v1_organization_all_pb.UpdateOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_organization_all_pb.UpdateOrganizationResponse|null) => void
  ): UnaryResponse;
  confirmUpdate(
    requestMessage: brank_as_rbac_gunk_v1_organization_all_pb.ConfirmUpdateRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_organization_all_pb.ConfirmUpdateResponse|null) => void
  ): UnaryResponse;
  confirmUpdate(
    requestMessage: brank_as_rbac_gunk_v1_organization_all_pb.ConfirmUpdateRequest,
    callback: (error: ServiceError|null, responseMessage: brank_as_rbac_gunk_v1_organization_all_pb.ConfirmUpdateResponse|null) => void
  ): UnaryResponse;
}

