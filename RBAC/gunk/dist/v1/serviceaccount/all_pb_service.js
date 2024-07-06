// package: brankas.rbac.v1.serviceaccount
// file: brank.as/rbac/gunk/v1/serviceaccount/all.proto

var brank_as_rbac_gunk_v1_serviceaccount_all_pb = require("./all_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var SvcAccountService = (function () {
  function SvcAccountService() {}
  SvcAccountService.serviceName = "brankas.rbac.v1.serviceaccount.SvcAccountService";
  return SvcAccountService;
}());

SvcAccountService.CreateAccount = {
  methodName: "CreateAccount",
  service: SvcAccountService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_serviceaccount_all_pb.CreateAccountRequest,
  responseType: brank_as_rbac_gunk_v1_serviceaccount_all_pb.CreateAccountResponse
};

SvcAccountService.ListAccounts = {
  methodName: "ListAccounts",
  service: SvcAccountService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_serviceaccount_all_pb.ListAccountsRequest,
  responseType: brank_as_rbac_gunk_v1_serviceaccount_all_pb.ListAccountsResponse
};

SvcAccountService.DisableAccount = {
  methodName: "DisableAccount",
  service: SvcAccountService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_serviceaccount_all_pb.DisableAccountRequest,
  responseType: brank_as_rbac_gunk_v1_serviceaccount_all_pb.DisableAccountResponse
};

exports.SvcAccountService = SvcAccountService;

function SvcAccountServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

SvcAccountServiceClient.prototype.createAccount = function createAccount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(SvcAccountService.CreateAccount, {
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

SvcAccountServiceClient.prototype.listAccounts = function listAccounts(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(SvcAccountService.ListAccounts, {
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

SvcAccountServiceClient.prototype.disableAccount = function disableAccount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(SvcAccountService.DisableAccount, {
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

exports.SvcAccountServiceClient = SvcAccountServiceClient;

var ValidationService = (function () {
  function ValidationService() {}
  ValidationService.serviceName = "brankas.rbac.v1.serviceaccount.ValidationService";
  return ValidationService;
}());

ValidationService.ValidateAccount = {
  methodName: "ValidateAccount",
  service: ValidationService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_serviceaccount_all_pb.ValidateAccountRequest,
  responseType: brank_as_rbac_gunk_v1_serviceaccount_all_pb.ValidateAccountResponse
};

exports.ValidationService = ValidationService;

function ValidationServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ValidationServiceClient.prototype.validateAccount = function validateAccount(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ValidationService.ValidateAccount, {
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

exports.ValidationServiceClient = ValidationServiceClient;

