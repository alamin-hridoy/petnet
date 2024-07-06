// package: brankas.rbac.v1.oauth2
// file: brank.as/rbac/gunk/v1/oauth2/all.proto

var brank_as_rbac_gunk_v1_oauth2_all_pb = require("./all_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var AuthClientService = (function () {
  function AuthClientService() {}
  AuthClientService.serviceName = "brankas.rbac.v1.oauth2.AuthClientService";
  return AuthClientService;
}());

AuthClientService.CreateClient = {
  methodName: "CreateClient",
  service: AuthClientService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_oauth2_all_pb.CreateClientRequest,
  responseType: brank_as_rbac_gunk_v1_oauth2_all_pb.CreateClientResponse
};

AuthClientService.UpdateClient = {
  methodName: "UpdateClient",
  service: AuthClientService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_oauth2_all_pb.UpdateClientRequest,
  responseType: brank_as_rbac_gunk_v1_oauth2_all_pb.UpdateClientResponse
};

AuthClientService.ListClients = {
  methodName: "ListClients",
  service: AuthClientService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_oauth2_all_pb.ListClientsRequest,
  responseType: brank_as_rbac_gunk_v1_oauth2_all_pb.ListClientsResponse
};

AuthClientService.DisableClient = {
  methodName: "DisableClient",
  service: AuthClientService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_oauth2_all_pb.DisableClientRequest,
  responseType: brank_as_rbac_gunk_v1_oauth2_all_pb.DisableClientResponse
};

exports.AuthClientService = AuthClientService;

function AuthClientServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

AuthClientServiceClient.prototype.createClient = function createClient(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AuthClientService.CreateClient, {
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

AuthClientServiceClient.prototype.updateClient = function updateClient(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AuthClientService.UpdateClient, {
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

AuthClientServiceClient.prototype.listClients = function listClients(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AuthClientService.ListClients, {
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

AuthClientServiceClient.prototype.disableClient = function disableClient(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(AuthClientService.DisableClient, {
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

exports.AuthClientServiceClient = AuthClientServiceClient;

