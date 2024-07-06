// package: rbac.brankas.errors
// file: brank.as/rbac/gunk/v1/errors/all.proto

import * as jspb from "google-protobuf";

export class Details extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getError(): string;
  setError(value: string): void;

  getMessagesMap(): jspb.Map<string, string>;
  clearMessagesMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Details.AsObject;
  static toObject(includeInstance: boolean, msg: Details): Details.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Details, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Details;
  static deserializeBinaryFromReader(message: Details, reader: jspb.BinaryReader): Details;
}

export namespace Details {
  export type AsObject = {
    id: string,
    error: string,
    messagesMap: Array<[string, string]>,
  }
}

