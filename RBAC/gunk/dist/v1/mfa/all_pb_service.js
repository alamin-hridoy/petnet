// package: mfa
// file: brank.as/rbac/gunk/v1/mfa/all.proto

var brank_as_rbac_gunk_v1_mfa_all_pb = require("./all_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var MFAService = (function () {
  function MFAService() {}
  MFAService.serviceName = "mfa.MFAService";
  return MFAService;
}());

MFAService.GetRegisteredMFA = {
  methodName: "GetRegisteredMFA",
  service: MFAService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_mfa_all_pb.GetRegisteredMFARequest,
  responseType: brank_as_rbac_gunk_v1_mfa_all_pb.GetRegisteredMFAResponse
};

MFAService.EnableMFA = {
  methodName: "EnableMFA",
  service: MFAService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_mfa_all_pb.EnableMFARequest,
  responseType: brank_as_rbac_gunk_v1_mfa_all_pb.EnableMFAResponse
};

MFAService.DisableMFA = {
  methodName: "DisableMFA",
  service: MFAService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_mfa_all_pb.DisableMFARequest,
  responseType: brank_as_rbac_gunk_v1_mfa_all_pb.DisableMFAResponse
};

exports.MFAService = MFAService;

function MFAServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

MFAServiceClient.prototype.getRegisteredMFA = function getRegisteredMFA(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(MFAService.GetRegisteredMFA, {
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

MFAServiceClient.prototype.enableMFA = function enableMFA(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(MFAService.EnableMFA, {
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

MFAServiceClient.prototype.disableMFA = function disableMFA(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(MFAService.DisableMFA, {
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

exports.MFAServiceClient = MFAServiceClient;

var MFAAuthService = (function () {
  function MFAAuthService() {}
  MFAAuthService.serviceName = "mfa.MFAAuthService";
  return MFAAuthService;
}());

MFAAuthService.InitiateMFA = {
  methodName: "InitiateMFA",
  service: MFAAuthService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_mfa_all_pb.InitiateMFARequest,
  responseType: brank_as_rbac_gunk_v1_mfa_all_pb.InitiateMFAResponse
};

MFAAuthService.ValidateMFA = {
  methodName: "ValidateMFA",
  service: MFAAuthService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_mfa_all_pb.ValidateMFARequest,
  responseType: brank_as_rbac_gunk_v1_mfa_all_pb.ValidateMFAResponse
};

MFAAuthService.RetryMFA = {
  methodName: "RetryMFA",
  service: MFAAuthService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_mfa_all_pb.RetryMFARequest,
  responseType: brank_as_rbac_gunk_v1_mfa_all_pb.RetryMFAResponse
};

MFAAuthService.ExternalMFA = {
  methodName: "ExternalMFA",
  service: MFAAuthService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_mfa_all_pb.ExternalMFARequest,
  responseType: brank_as_rbac_gunk_v1_mfa_all_pb.ExternalMFAResponse
};

exports.MFAAuthService = MFAAuthService;

function MFAAuthServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

MFAAuthServiceClient.prototype.initiateMFA = function initiateMFA(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(MFAAuthService.InitiateMFA, {
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

MFAAuthServiceClient.prototype.validateMFA = function validateMFA(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(MFAAuthService.ValidateMFA, {
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

MFAAuthServiceClient.prototype.retryMFA = function retryMFA(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(MFAAuthService.RetryMFA, {
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

MFAAuthServiceClient.prototype.externalMFA = function externalMFA(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(MFAAuthService.ExternalMFA, {
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

exports.MFAAuthServiceClient = MFAAuthServiceClient;

