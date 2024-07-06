// package: organization
// file: brank.as/rbac/gunk/v1/organization/all.proto

var brank_as_rbac_gunk_v1_organization_all_pb = require("./all_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var OrganizationService = (function () {
  function OrganizationService() {}
  OrganizationService.serviceName = "organization.OrganizationService";
  return OrganizationService;
}());

OrganizationService.GetOrganization = {
  methodName: "GetOrganization",
  service: OrganizationService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_organization_all_pb.GetOrganizationRequest,
  responseType: brank_as_rbac_gunk_v1_organization_all_pb.GetOrganizationResponse
};

OrganizationService.UpdateOrganization = {
  methodName: "UpdateOrganization",
  service: OrganizationService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_organization_all_pb.UpdateOrganizationRequest,
  responseType: brank_as_rbac_gunk_v1_organization_all_pb.UpdateOrganizationResponse
};

OrganizationService.ConfirmUpdate = {
  methodName: "ConfirmUpdate",
  service: OrganizationService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_organization_all_pb.ConfirmUpdateRequest,
  responseType: brank_as_rbac_gunk_v1_organization_all_pb.ConfirmUpdateResponse
};

exports.OrganizationService = OrganizationService;

function OrganizationServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

OrganizationServiceClient.prototype.getOrganization = function getOrganization(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(OrganizationService.GetOrganization, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

OrganizationServiceClient.prototype.updateOrganization = function updateOrganization(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(OrganizationService.UpdateOrganization, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

OrganizationServiceClient.prototype.confirmUpdate = function confirmUpdate(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(OrganizationService.ConfirmUpdate, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.OrganizationServiceClient = OrganizationServiceClient;

