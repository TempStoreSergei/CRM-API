/* eslint-disable */
import { GrpcMethod, GrpcStreamMethod } from "@nestjs/microservices";
import { Observable } from "rxjs";
import { Timestamp } from "../../../../google/protobuf/timestamp";

export const protobufPackage = "messenger";

export enum MessageType {
  UNKNOWN = 0,
  TEXT = 1,
  IMAGE = 2,
  VIDEO = 3,
  UNRECOGNIZED = -1,
}

export interface User {
  id: string;
  email: string;
}

export interface Chat {
  id: string;
  user: User | undefined;
  userId: string;
  text: string;
  mediaUrl: string;
  type: MessageType;
  replyTo: string;
  readAt: Timestamp | undefined;
  createdAt: Timestamp | undefined;
  updatedAt: Timestamp | undefined;
  deletedAt: Timestamp | undefined;
}

export interface CreateChatRequest {
  user: User | undefined;
  text: string;
  mediaUrl: string;
  type: MessageType;
  replyTo: string;
}

export interface GetChatRequest {
  id: string;
}

export const MESSENGER_PACKAGE_NAME = "messenger";

export interface ChatServiceClient {
  createChat(request: Observable<CreateChatRequest>): Observable<Chat>;

  getChat(request: GetChatRequest): Observable<Chat>;
}

export interface ChatServiceController {
  createChat(request: Observable<CreateChatRequest>): Observable<Chat>;

  getChat(request: GetChatRequest): Observable<Chat>;
}

export function ChatServiceControllerMethods() {
  return function (constructor: Function) {
    const grpcMethods: string[] = ["getChat"];
    for (const method of grpcMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcMethod("ChatService", method)(constructor.prototype[method], method, descriptor);
    }
    const grpcStreamMethods: string[] = ["createChat"];
    for (const method of grpcStreamMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcStreamMethod("ChatService", method)(constructor.prototype[method], method, descriptor);
    }
  };
}

export const CHAT_SERVICE_NAME = "ChatService";
