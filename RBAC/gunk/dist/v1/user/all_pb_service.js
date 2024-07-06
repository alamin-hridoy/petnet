// package: user
// file: brank.as/rbac/gunk/v1/user/all.proto

var brank_as_rbac_gunk_v1_user_all_pb = require("./all_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var Signup = (function () {
  function Signup() {}
  Signup.serviceName = "user.Signup";
  return Signup;
}());

Signup.Signup = {
  methodName: "Signup",
  service: Signup,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.SignupRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.SignupResponse
};

Signup.ResendConfirmEmail = {
  methodName: "ResendConfirmEmail",
  service: Signup,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.ResendConfirmEmailRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.ResendConfirmEmailResponse
};

Signup.EmailConfirmation = {
  methodName: "EmailConfirmation",
  service: Signup,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.EmailConfirmationRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.EmailConfirmationResponse
};

Signup.ForgotPassword = {
  methodName: "ForgotPassword",
  service: Signup,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.ForgotPasswordRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.ForgotPasswordResponse
};

Signup.ResetPassword = {
  methodName: "ResetPassword",
  service: Signup,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.ResetPasswordRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.ResetPasswordResponse
};

exports.Signup = Signup;

function SignupClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

SignupClient.prototype.signup = function signup(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(Signup.Signup, {
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

SignupClient.prototype.resendConfirmEmail = function resendConfirmEmail(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(Signup.ResendConfirmEmail, {
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

SignupClient.prototype.emailConfirmation = function emailConfirmation(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(Signup.EmailConfirmation, {
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

SignupClient.prototype.forgotPassword = function forgotPassword(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(Signup.ForgotPassword, {
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

SignupClient.prototype.resetPassword = function resetPassword(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(Signup.ResetPassword, {
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

exports.SignupClient = SignupClient;

var UserService = (function () {
  function UserService() {}
  UserService.serviceName = "user.UserService";
  return UserService;
}());

UserService.GetUser = {
  methodName: "GetUser",
  service: UserService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.GetUserRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.GetUserResponse
};

UserService.ListUsers = {
  methodName: "ListUsers",
  service: UserService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.ListUsersRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.ListUsersResponse
};

UserService.ChangePassword = {
  methodName: "ChangePassword",
  service: UserService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.ChangePasswordRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.ChangePasswordResponse
};

UserService.ConfirmUpdate = {
  methodName: "ConfirmUpdate",
  service: UserService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.ConfirmUpdateRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.ConfirmUpdateResponse
};

UserService.UpdateUser = {
  methodName: "UpdateUser",
  service: UserService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.UpdateUserRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.UpdateUserResponse
};

UserService.DisableUser = {
  methodName: "DisableUser",
  service: UserService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.DisableUserRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.DisableUserResponse
};

UserService.EnableUser = {
  methodName: "EnableUser",
  service: UserService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.EnableUserRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.EnableUserResponse
};

exports.UserService = UserService;

function UserServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

UserServiceClient.prototype.getUser = function getUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserService.GetUser, {
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

UserServiceClient.prototype.listUsers = function listUsers(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserService.ListUsers, {
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

UserServiceClient.prototype.changePassword = function changePassword(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserService.ChangePassword, {
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

UserServiceClient.prototype.confirmUpdate = function confirmUpdate(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserService.ConfirmUpdate, {
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

UserServiceClient.prototype.updateUser = function updateUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserService.UpdateUser, {
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

UserServiceClient.prototype.disableUser = function disableUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserService.DisableUser, {
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

UserServiceClient.prototype.enableUser = function enableUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserService.EnableUser, {
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

exports.UserServiceClient = UserServiceClient;

var UserAuthService = (function () {
  function UserAuthService() {}
  UserAuthService.serviceName = "user.UserAuthService";
  return UserAuthService;
}());

UserAuthService.AuthenticateUser = {
  methodName: "AuthenticateUser",
  service: UserAuthService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_user_all_pb.AuthenticateUserRequest,
  responseType: brank_as_rbac_gunk_v1_user_all_pb.AuthenticateUserResponse
};

exports.UserAuthService = UserAuthService;

function UserAuthServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

UserAuthServiceClient.prototype.authenticateUser = function authenticateUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(UserAuthService.AuthenticateUser, {
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

exports.UserAuthServiceClient = UserAuthServiceClient;

