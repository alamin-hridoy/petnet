// package: permissions
// file: brank.as/rbac/gunk/v1/permissions/all.proto

var brank_as_rbac_gunk_v1_permissions_all_pb = require("./all_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var PermissionService = (function () {
  function PermissionService() {}
  PermissionService.serviceName = "permissions.PermissionService";
  return PermissionService;
}());

PermissionService.CreatePermission = {
  methodName: "CreatePermission",
  service: PermissionService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.CreatePermissionRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.CreatePermissionResponse
};

PermissionService.ListPermission = {
  methodName: "ListPermission",
  service: PermissionService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.ListPermissionRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.ListPermissionResponse
};

PermissionService.DeletePermission = {
  methodName: "DeletePermission",
  service: PermissionService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.DeletePermissionRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.DeletePermissionResponse
};

PermissionService.AssignPermission = {
  methodName: "AssignPermission",
  service: PermissionService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.AssignPermissionRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.AssignPermissionResponse
};

PermissionService.RevokePermission = {
  methodName: "RevokePermission",
  service: PermissionService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.RevokePermissionRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.RevokePermissionResponse
};

exports.PermissionService = PermissionService;

function PermissionServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

PermissionServiceClient.prototype.createPermission = function createPermission(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PermissionService.CreatePermission, {
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

PermissionServiceClient.prototype.listPermission = function listPermission(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PermissionService.ListPermission, {
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

PermissionServiceClient.prototype.deletePermission = function deletePermission(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PermissionService.DeletePermission, {
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

PermissionServiceClient.prototype.assignPermission = function assignPermission(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PermissionService.AssignPermission, {
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

PermissionServiceClient.prototype.revokePermission = function revokePermission(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PermissionService.RevokePermission, {
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

exports.PermissionServiceClient = PermissionServiceClient;

var RoleService = (function () {
  function RoleService() {}
  RoleService.serviceName = "permissions.RoleService";
  return RoleService;
}());

RoleService.CreateRole = {
  methodName: "CreateRole",
  service: RoleService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.CreateRoleRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.CreateRoleResponse
};

RoleService.ListRole = {
  methodName: "ListRole",
  service: RoleService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.ListRoleRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.ListRoleResponse
};

RoleService.ListUserRoles = {
  methodName: "ListUserRoles",
  service: RoleService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.ListUserRolesRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.ListUserRolesResponse
};

RoleService.UpdateRole = {
  methodName: "UpdateRole",
  service: RoleService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.UpdateRoleRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.UpdateRoleResponse
};

RoleService.DeleteRole = {
  methodName: "DeleteRole",
  service: RoleService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.DeleteRoleRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.DeleteRoleResponse
};

RoleService.AssignRolePermission = {
  methodName: "AssignRolePermission",
  service: RoleService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.AssignRolePermissionRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.AssignRolePermissionResponse
};

RoleService.RevokeRolePermission = {
  methodName: "RevokeRolePermission",
  service: RoleService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.RevokeRolePermissionRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.RevokeRolePermissionResponse
};

RoleService.AddUser = {
  methodName: "AddUser",
  service: RoleService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.AddUserRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.AddUserResponse
};

RoleService.RemoveUser = {
  methodName: "RemoveUser",
  service: RoleService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.RemoveUserRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.RemoveUserResponse
};

exports.RoleService = RoleService;

function RoleServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

RoleServiceClient.prototype.createRole = function createRole(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RoleService.CreateRole, {
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

RoleServiceClient.prototype.listRole = function listRole(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RoleService.ListRole, {
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

RoleServiceClient.prototype.listUserRoles = function listUserRoles(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RoleService.ListUserRoles, {
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

RoleServiceClient.prototype.updateRole = function updateRole(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RoleService.UpdateRole, {
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

RoleServiceClient.prototype.deleteRole = function deleteRole(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RoleService.DeleteRole, {
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

RoleServiceClient.prototype.assignRolePermission = function assignRolePermission(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RoleService.AssignRolePermission, {
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

RoleServiceClient.prototype.revokeRolePermission = function revokeRolePermission(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RoleService.RevokeRolePermission, {
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

RoleServiceClient.prototype.addUser = function addUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RoleService.AddUser, {
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

RoleServiceClient.prototype.removeUser = function removeUser(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RoleService.RemoveUser, {
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

exports.RoleServiceClient = RoleServiceClient;

var ProductService = (function () {
  function ProductService() {}
  ProductService.serviceName = "permissions.ProductService";
  return ProductService;
}());

ProductService.GrantService = {
  methodName: "GrantService",
  service: ProductService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.GrantServiceRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.GrantServiceResponse
};

ProductService.ListServiceAssignments = {
  methodName: "ListServiceAssignments",
  service: ProductService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.ListServiceAssignmentsRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.ListServiceAssignmentsResponse
};

ProductService.RevokeService = {
  methodName: "RevokeService",
  service: ProductService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.RevokeServiceRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.RevokeServiceResponse
};

ProductService.PublicService = {
  methodName: "PublicService",
  service: ProductService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.PublicServiceRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.PublicServiceResponse
};

ProductService.ListServices = {
  methodName: "ListServices",
  service: ProductService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.ListServicesRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.ListServicesResponse
};

exports.ProductService = ProductService;

function ProductServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ProductServiceClient.prototype.grantService = function grantService(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ProductService.GrantService, {
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

ProductServiceClient.prototype.listServiceAssignments = function listServiceAssignments(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ProductService.ListServiceAssignments, {
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

ProductServiceClient.prototype.revokeService = function revokeService(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ProductService.RevokeService, {
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

ProductServiceClient.prototype.publicService = function publicService(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ProductService.PublicService, {
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

ProductServiceClient.prototype.listServices = function listServices(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ProductService.ListServices, {
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

exports.ProductServiceClient = ProductServiceClient;

var ValidationService = (function () {
  function ValidationService() {}
  ValidationService.serviceName = "permissions.ValidationService";
  return ValidationService;
}());

ValidationService.ValidatePermission = {
  methodName: "ValidatePermission",
  service: ValidationService,
  requestStream: false,
  responseStream: false,
  requestType: brank_as_rbac_gunk_v1_permissions_all_pb.ValidatePermissionRequest,
  responseType: brank_as_rbac_gunk_v1_permissions_all_pb.ValidatePermissionResponse
};

exports.ValidationService = ValidationService;

function ValidationServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ValidationServiceClient.prototype.validatePermission = function validatePermission(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ValidationService.ValidatePermission, {
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

