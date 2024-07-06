// package: rbac.brankas.consent
// file: brank.as/rbac/gunk/v1/consent/all.proto

var brank_as_rbac_gunk_v1_consent_all_pb = require("./all_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var GrantService = (function () {
  function GrantService() {}
  GrantService.serviceName = "rbac.brankas.consent.GrantService";
  return GrantService;
}());

GrantService.ServeGrant = {
  methodName: "ServeGrant",
  service: GrantService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_consent_all_pb.ServeGrantRequest,
  responseType: brank_as_rbac_gunk_v1_consent_all_pb.ServeGrantResponse
};

GrantService.Grant = {
  methodName: "Grant",
  service: GrantService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_consent_all_pb.GrantRequest,
  responseType: brank_as_rbac_gunk_v1_consent_all_pb.GrantResponse
};

exports.GrantService = GrantService;

function GrantServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

GrantServiceClient.prototype.serveGrant = function serveGrant(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GrantService.ServeGrant, {
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

GrantServiceClient.prototype.grant = function grant(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(GrantService.Grant, {
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

exports.GrantServiceClient = GrantServiceClient;

var ScopeService = (function () {
  function ScopeService() {}
  ScopeService.serviceName = "rbac.brankas.consent.ScopeService";
  return ScopeService;
}());

ScopeService.UpsertScope = {
  methodName: "UpsertScope",
  service: ScopeService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_consent_all_pb.UpsertScopeRequest,
  responseType: brank_as_rbac_gunk_v1_consent_all_pb.UpsertScopeResponse
};

ScopeService.UpdateGroup = {
  methodName: "UpdateGroup",
  service: ScopeService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_consent_all_pb.UpdateGroupRequest,
  responseType: brank_as_rbac_gunk_v1_consent_all_pb.UpdateGroupResponse
};

ScopeService.GetScope = {
  methodName: "GetScope",
  service: ScopeService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_consent_all_pb.GetScopeRequest,
  responseType: brank_as_rbac_gunk_v1_consent_all_pb.GetScopeResponse
};

exports.ScopeService = ScopeService;

function ScopeServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ScopeServiceClient.prototype.upsertScope = function upsertScope(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ScopeService.UpsertScope, {
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

ScopeServiceClient.prototype.updateGroup = function updateGroup(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ScopeService.UpdateGroup, {
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

ScopeServiceClient.prototype.getScope = function getScope(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ScopeService.GetScope, {
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

exports.ScopeServiceClient = ScopeServiceClient;

