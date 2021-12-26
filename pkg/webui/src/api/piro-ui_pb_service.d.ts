// package: v1
// file: piro-ui.proto

import * as piro_ui_pb from "./piro-ui_pb";
import {grpc} from "@improbable-eng/grpc-web";

type PiroUIListJobSpecs = {
  readonly methodName: string;
  readonly service: typeof PiroUI;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof piro_ui_pb.ListJobSpecsRequest;
  readonly responseType: typeof piro_ui_pb.ListJobSpecsResponse;
};

type PiroUIIsReadOnly = {
  readonly methodName: string;
  readonly service: typeof PiroUI;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof piro_ui_pb.IsReadOnlyRequest;
  readonly responseType: typeof piro_ui_pb.IsReadOnlyResponse;
};

export class PiroUI {
  static readonly serviceName: string;
  static readonly ListJobSpecs: PiroUIListJobSpecs;
  static readonly IsReadOnly: PiroUIIsReadOnly;
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

export class PiroUIClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  listJobSpecs(requestMessage: piro_ui_pb.ListJobSpecsRequest, metadata?: grpc.Metadata): ResponseStream<piro_ui_pb.ListJobSpecsResponse>;
  isReadOnly(
    requestMessage: piro_ui_pb.IsReadOnlyRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: piro_ui_pb.IsReadOnlyResponse|null) => void
  ): UnaryResponse;
  isReadOnly(
    requestMessage: piro_ui_pb.IsReadOnlyRequest,
    callback: (error: ServiceError|null, responseMessage: piro_ui_pb.IsReadOnlyResponse|null) => void
  ): UnaryResponse;
}
