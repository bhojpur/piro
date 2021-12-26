// package: v1
// file: piro-ui.proto

var piro_ui_pb = require("./piro-ui_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var PiroUI = (function () {
  function PiroUI() {}
  PiroUI.serviceName = "v1.PiroUI";
  return PiroUI;
}());

PiroUI.ListJobSpecs = {
  methodName: "ListJobSpecs",
  service: PiroUI,
  requestStream: false,
  responseStream: true,
  requestType: piro_ui_pb.ListJobSpecsRequest,
  responseType: piro_ui_pb.ListJobSpecsResponse
};

PiroUI.IsReadOnly = {
  methodName: "IsReadOnly",
  service: PiroUI,
  requestStream: false,
  responseStream: false,
  requestType: piro_ui_pb.IsReadOnlyRequest,
  responseType: piro_ui_pb.IsReadOnlyResponse
};

exports.PiroUI = PiroUI;

function PiroUIClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

PiroUIClient.prototype.listJobSpecs = function listJobSpecs(requestMessage, metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.invoke(PiroUI.ListJobSpecs, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onMessage: function (responseMessage) {
      listeners.data.forEach(function (handler) {
        handler(responseMessage);
      });
    },
    onEnd: function (status, statusMessage, trailers) {
      listeners.status.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners.end.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners = null;
    }
  });
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

PiroUIClient.prototype.isReadOnly = function isReadOnly(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PiroUI.IsReadOnly, {
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

exports.PiroUIClient = PiroUIClient;
