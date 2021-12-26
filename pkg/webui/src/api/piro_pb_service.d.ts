// package: v1
// file: piro.proto

import * as piro_pb from "./piro_pb";
import {grpc} from "@improbable-eng/grpc-web";

type PiroServiceStartLocalJob = {
  readonly methodName: string;
  readonly service: typeof PiroService;
  readonly requestStream: true;
  readonly responseStream: false;
  readonly requestType: typeof piro_pb.StartLocalJobRequest;
  readonly responseType: typeof piro_pb.StartJobResponse;
};

type PiroServiceStartGitHubJob = {
  readonly methodName: string;
  readonly service: typeof PiroService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof piro_pb.StartGitHubJobRequest;
  readonly responseType: typeof piro_pb.StartJobResponse;
};

type PiroServiceStartFromPreviousJob = {
  readonly methodName: string;
  readonly service: typeof PiroService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof piro_pb.StartFromPreviousJobRequest;
  readonly responseType: typeof piro_pb.StartJobResponse;
};

type PiroServiceStartJob = {
  readonly methodName: string;
  readonly service: typeof PiroService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof piro_pb.StartJobRequest;
  readonly responseType: typeof piro_pb.StartJobResponse;
};

type PiroServiceListJobs = {
  readonly methodName: string;
  readonly service: typeof PiroService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof piro_pb.ListJobsRequest;
  readonly responseType: typeof piro_pb.ListJobsResponse;
};

type PiroServiceSubscribe = {
  readonly methodName: string;
  readonly service: typeof PiroService;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof piro_pb.SubscribeRequest;
  readonly responseType: typeof piro_pb.SubscribeResponse;
};

type PiroServiceGetJob = {
  readonly methodName: string;
  readonly service: typeof PiroService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof piro_pb.GetJobRequest;
  readonly responseType: typeof piro_pb.GetJobResponse;
};

type PiroServiceListen = {
  readonly methodName: string;
  readonly service: typeof PiroService;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof piro_pb.ListenRequest;
  readonly responseType: typeof piro_pb.ListenResponse;
};

type PiroServiceStopJob = {
  readonly methodName: string;
  readonly service: typeof PiroService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof piro_pb.StopJobRequest;
  readonly responseType: typeof piro_pb.StopJobResponse;
};

export class PiroService {
  static readonly serviceName: string;
  static readonly StartLocalJob: PiroServiceStartLocalJob;
  static readonly StartGitHubJob: PiroServiceStartGitHubJob;
  static readonly StartFromPreviousJob: PiroServiceStartFromPreviousJob;
  static readonly StartJob: PiroServiceStartJob;
  static readonly ListJobs: PiroServiceListJobs;
  static readonly Subscribe: PiroServiceSubscribe;
  static readonly GetJob: PiroServiceGetJob;
  static readonly Listen: PiroServiceListen;
  static readonly StopJob: PiroServiceStopJob;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class PiroServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  startLocalJob(metadata?: grpc.Metadata): RequestStream<piro_pb.StartLocalJobRequest>;
  startGitHubJob(
    requestMessage: piro_pb.StartGitHubJobRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: piro_pb.StartJobResponse|null) => void
  ): UnaryResponse;
  startGitHubJob(
    requestMessage: piro_pb.StartGitHubJobRequest,
    callback: (error: ServiceError|null, responseMessage: piro_pb.StartJobResponse|null) => void
  ): UnaryResponse;
  startFromPreviousJob(
    requestMessage: piro_pb.StartFromPreviousJobRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: piro_pb.StartJobResponse|null) => void
  ): UnaryResponse;
  startFromPreviousJob(
    requestMessage: piro_pb.StartFromPreviousJobRequest,
    callback: (error: ServiceError|null, responseMessage: piro_pb.StartJobResponse|null) => void
  ): UnaryResponse;
  startJob(
    requestMessage: piro_pb.StartJobRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: piro_pb.StartJobResponse|null) => void
  ): UnaryResponse;
  startJob(
    requestMessage: piro_pb.StartJobRequest,
    callback: (error: ServiceError|null, responseMessage: piro_pb.StartJobResponse|null) => void
  ): UnaryResponse;
  listJobs(
    requestMessage: piro_pb.ListJobsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: piro_pb.ListJobsResponse|null) => void
  ): UnaryResponse;
  listJobs(
    requestMessage: piro_pb.ListJobsRequest,
    callback: (error: ServiceError|null, responseMessage: piro_pb.ListJobsResponse|null) => void
  ): UnaryResponse;
  subscribe(requestMessage: piro_pb.SubscribeRequest, metadata?: grpc.Metadata): ResponseStream<piro_pb.SubscribeResponse>;
  getJob(
    requestMessage: piro_pb.GetJobRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: piro_pb.GetJobResponse|null) => void
  ): UnaryResponse;
  getJob(
    requestMessage: piro_pb.GetJobRequest,
    callback: (error: ServiceError|null, responseMessage: piro_pb.GetJobResponse|null) => void
  ): UnaryResponse;
  listen(requestMessage: piro_pb.ListenRequest, metadata?: grpc.Metadata): ResponseStream<piro_pb.ListenResponse>;
  stopJob(
    requestMessage: piro_pb.StopJobRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: piro_pb.StopJobResponse|null) => void
  ): UnaryResponse;
  stopJob(
    requestMessage: piro_pb.StopJobRequest,
    callback: (error: ServiceError|null, responseMessage: piro_pb.StopJobResponse|null) => void
  ): UnaryResponse;
}
