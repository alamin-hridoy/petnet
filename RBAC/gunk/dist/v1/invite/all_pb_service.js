// package: brankas.rbac.v1.invite
// file: brank.as/rbac/gunk/v1/invite/all.proto

var brank_as_rbac_gunk_v1_invite_all_pb = require("./all_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var InviteService = (function () {
  function InviteService() {}
  InviteService.serviceName = "brankas.rbac.v1.invite.InviteService";
  return InviteService;
}());

InviteService.InviteUser = {
  methodName: "InviteUser",
  service: InviteService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_invite_all_pb.InviteUserRequest,
  responseType: brank_as_rbac_gunk_v1_invite_all_pb.InviteUserResponse
};

InviteService.Resend = {
  methodName: "Resend",
  service: InviteService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_invite_all_pb.ResendRequest,
  responseType: brank_as_rbac_gunk_v1_invite_all_pb.ResendResponse
};

InviteService.ListInvite = {
  methodName: "ListInvite",
  service: InviteService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_invite_all_pb.ListInviteRequest,
  responseType: brank_as_rbac_gunk_v1_invite_all_pb.ListInviteResponse
};

InviteService.RetrieveInvite = {
  methodName: "RetrieveInvite",
  service: InviteService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_invite_all_pb.RetrieveInviteRequest,
  responseType: brank_as_rbac_gunk_v1_invite_all_pb.RetrieveInviteResponse
};

InviteService.CancelInvite = {
  methodName: "CancelInvite",
  service: InviteService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_invite_all_pb.CancelInviteRequest,
  responseType: brank_as_rbac_gunk_v1_invite_all_pb.CancelInviteResponse
};

InviteService.Approve = {
  methodName: "Approve",
  service: InviteService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_invite_all_pb.ApproveRequest,
  responseType: brank_as_rbac_gunk_v1_invite_all_pb.ApproveResponse
};

InviteService.Revoke = {
  methodName: "Revoke",
  service: InviteService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_invite_all_pb.RevokeRequest,
  responseType: brank_as_rbac_gunk_v1_invite_all_pb.RevokeResponse
};

exports.InviteService = InviteService;

function InviteServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

InviteServiceClient.prototype.inviteUser = function inviteUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(InviteService.InviteUser, {
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

InviteServiceClient.prototype.resend = function resend(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(InviteService.Resend, {
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

InviteServiceClient.prototype.listInvite = function listInvite(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(InviteService.ListInvite, {
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

InviteServiceClient.prototype.retrieveInvite = function retrieveInvite(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(InviteService.RetrieveInvite, {
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

InviteServiceClient.prototype.cancelInvite = function cancelInvite(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(InviteService.CancelInvite, {
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

InviteServiceClient.prototype.approve = function approve(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(InviteService.Approve, {
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

InviteServiceClient.prototype.revoke = function revoke(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(InviteService.Revoke, {
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

exports.InviteServiceClient = InviteServiceClient;

