// package: authenticate
// file: brank.as/rbac/gunk/v1/authenticate/all.proto

var brank_as_rbac_gunk_v1_authenticate_all_pb = require("./all_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var SessionService = (function () {
  function SessionService() {}
  SessionService.serviceName = "authenticate.SessionService";
  return SessionService;
}());

SessionService.Login = {
  methodName: "Login",
  service: SessionService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_authenticate_all_pb.LoginRequest,
  responseType: brank_as_rbac_gunk_v1_authenticate_all_pb.Session
};

SessionService.GetSession = {
  methodName: "GetSession",
  service: SessionService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_authenticate_all_pb.GetSessionRequest,
  responseType: brank_as_rbac_gunk_v1_authenticate_all_pb.Session
};

SessionService.RetryMFA = {
  methodName: "RetryMFA",
  service: SessionService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_authenticate_all_pb.RetryMFARequest,
  responseType: brank_as_rbac_gunk_v1_authenticate_all_pb.Session
};

exports.SessionService = SessionService;

function SessionServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

SessionServiceClient.prototype.login = function login(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(SessionService.Login, {
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

SessionServiceClient.prototype.getSession = function getSession(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(SessionService.GetSession, {
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

SessionServiceClient.prototype.retryMFA = function retryMFA(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(SessionService.RetryMFA, {
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

exports.SessionServiceClient = SessionServiceClient;

