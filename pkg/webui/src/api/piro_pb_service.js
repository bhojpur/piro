// package: v1
// file: piro.proto

var piro_pb = require("./piro_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var PiroService = (function () {
  function PiroService() {}
  PiroService.serviceName = "v1.PiroService";
  return PiroService;
}());

PiroService.StartLocalJob = {
  methodName: "StartLocalJob",
  service: PiroService,
  requestStream: true,
  responseStream: false,
  requestType: piro_pb.StartLocalJobRequest,
  responseType: piro_pb.StartJobResponse
};

PiroService.StartGitHubJob = {
  methodName: "StartGitHubJob",
  service: PiroService,
  requestStream: false,
  responseStream: false,
  requestType: piro_pb.StartGitHubJobRequest,
  responseType: piro_pb.StartJobResponse
};

PiroService.StartFromPreviousJob = {
  methodName: "StartFromPreviousJob",
  service: PiroService,
  requestStream: false,
  responseStream: false,
  requestType: piro_pb.StartFromPreviousJobRequest,
  responseType: piro_pb.StartJobResponse
};

PiroService.StartJob = {
  methodName: "StartJob",
  service: PiroService,
  requestStream: false,
  responseStream: false,
  requestType: piro_pb.StartJobRequest,
  responseType: piro_pb.StartJobResponse
};

PiroService.ListJobs = {
  methodName: "ListJobs",
  service: PiroService,
  requestStream: false,
  responseStream: false,
  requestType: piro_pb.ListJobsRequest,
  responseType: piro_pb.ListJobsResponse
};

PiroService.Subscribe = {
  methodName: "Subscribe",
  service: PiroService,
  requestStream: false,
  responseStream: true,
  requestType: piro_pb.SubscribeRequest,
  responseType: piro_pb.SubscribeResponse
};

PiroService.GetJob = {
  methodName: "GetJob",
  service: PiroService,
  requestStream: false,
  responseStream: false,
  requestType: piro_pb.GetJobRequest,
  responseType: piro_pb.GetJobResponse
};

PiroService.Listen = {
  methodName: "Listen",
  service: PiroService,
  requestStream: false,
  responseStream: true,
  requestType: piro_pb.ListenRequest,
  responseType: piro_pb.ListenResponse
};

PiroService.StopJob = {
  methodName: "StopJob",
  service: PiroService,
  requestStream: false,
  responseStream: false,
  requestType: piro_pb.StopJobRequest,
  responseType: piro_pb.StopJobResponse
};

exports.PiroService = PiroService;

function PiroServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

PiroServiceClient.prototype.startLocalJob = function startLocalJob(metadata) {
  var listeners = {
    end: [],
    status: []
  };
  var client = grpc.client(PiroService.StartLocalJob, {
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport
  });
  client.onEnd(function (status, statusMessage, trailers) {
    listeners.status.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners.end.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners = null;
  });
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    write: function (requestMessage) {
      if (!client.started) {
        client.start(metadata);
      }
      client.send(requestMessage);
      return this;
    },
    end: function () {
      client.finishSend();
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

PiroServiceClient.prototype.startGitHubJob = function startGitHubJob(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PiroService.StartGitHubJob, {
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

PiroServiceClient.prototype.startFromPreviousJob = function startFromPreviousJob(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PiroService.StartFromPreviousJob, {
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

PiroServiceClient.prototype.startJob = function startJob(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PiroService.StartJob, {
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

PiroServiceClient.prototype.listJobs = function listJobs(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PiroService.ListJobs, {
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

PiroServiceClient.prototype.subscribe = function subscribe(requestMessage, metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.invoke(PiroService.Subscribe, {
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

PiroServiceClient.prototype.getJob = function getJob(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PiroService.GetJob, {
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

PiroServiceClient.prototype.listen = function listen(requestMessage, metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.invoke(PiroService.Listen, {
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

PiroServiceClient.prototype.stopJob = function stopJob(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PiroService.StopJob, {
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

exports.PiroServiceClient = PiroServiceClient;
