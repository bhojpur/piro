// package: v1
// file: piro-ui.proto

import * as jspb from "google-protobuf";
import * as piro_pb from "./piro_pb";

export class ListJobSpecsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListJobSpecsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListJobSpecsRequest): ListJobSpecsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListJobSpecsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListJobSpecsRequest;
  static deserializeBinaryFromReader(message: ListJobSpecsRequest, reader: jspb.BinaryReader): ListJobSpecsRequest;
}

export namespace ListJobSpecsRequest {
  export type AsObject = {
  }
}

export class ListJobSpecsResponse extends jspb.Message {
  hasRepo(): boolean;
  clearRepo(): void;
  getRepo(): piro_pb.Repository | undefined;
  setRepo(value?: piro_pb.Repository): void;

  getName(): string;
  setName(value: string): void;

  getPath(): string;
  setPath(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  clearArgumentsList(): void;
  getArgumentsList(): Array<DesiredAnnotation>;
  setArgumentsList(value: Array<DesiredAnnotation>): void;
  addArguments(value?: DesiredAnnotation, index?: number): DesiredAnnotation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListJobSpecsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListJobSpecsResponse): ListJobSpecsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListJobSpecsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListJobSpecsResponse;
  static deserializeBinaryFromReader(message: ListJobSpecsResponse, reader: jspb.BinaryReader): ListJobSpecsResponse;
}

export namespace ListJobSpecsResponse {
  export type AsObject = {
    repo?: piro_pb.Repository.AsObject,
    name: string,
    path: string,
    description: string,
    argumentsList: Array<DesiredAnnotation.AsObject>,
  }
}

export class DesiredAnnotation extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getRequired(): boolean;
  setRequired(value: boolean): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DesiredAnnotation.AsObject;
  static toObject(includeInstance: boolean, msg: DesiredAnnotation): DesiredAnnotation.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DesiredAnnotation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DesiredAnnotation;
  static deserializeBinaryFromReader(message: DesiredAnnotation, reader: jspb.BinaryReader): DesiredAnnotation;
}

export namespace DesiredAnnotation {
  export type AsObject = {
    name: string,
    required: boolean,
    description: string,
  }
}

export class IsReadOnlyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IsReadOnlyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: IsReadOnlyRequest): IsReadOnlyRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: IsReadOnlyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IsReadOnlyRequest;
  static deserializeBinaryFromReader(message: IsReadOnlyRequest, reader: jspb.BinaryReader): IsReadOnlyRequest;
}

export namespace IsReadOnlyRequest {
  export type AsObject = {
  }
}

export class IsReadOnlyResponse extends jspb.Message {
  getReadonly(): boolean;
  setReadonly(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IsReadOnlyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: IsReadOnlyResponse): IsReadOnlyResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: IsReadOnlyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IsReadOnlyResponse;
  static deserializeBinaryFromReader(message: IsReadOnlyResponse, reader: jspb.BinaryReader): IsReadOnlyResponse;
}

export namespace IsReadOnlyResponse {
  export type AsObject = {
    readonly: boolean,
  }
}
