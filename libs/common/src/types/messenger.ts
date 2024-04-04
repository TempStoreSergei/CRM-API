/* eslint-disable */
import { GrpcMethod, GrpcStreamMethod } from "@nestjs/microservices";
import { Observable } from "rxjs";
import { Empty } from "../../../../google/protobuf/empty";
import { Timestamp } from "../../../../google/protobuf/timestamp";

export const protobufPackage = "messenger";

export interface Conversation {
  id: string;
  title: string;
  type: Conversation_ConversationType;
  /** Collection with member ids */
  memberIds: string[];
  image: Uint8Array;
  creationTime:
    | Timestamp
    | undefined;
  /** Last update time (changes on conversations or last message time sent) */
  updateTime: Timestamp | undefined;
}

export enum Conversation_ConversationType {
  INDIVIDUAL = 0,
  GROUP = 1,
  UNRECOGNIZED = -1,
}

export interface Message {
  id: string;
  senderId: string;
  event?: ActionEvent | undefined;
  text?: string | undefined;
  audio?: Uint8Array | undefined;
  image?: Uint8Array | undefined;
  creationTime: Timestamp | undefined;
  deliveryTime: Timestamp | undefined;
}

export interface ActionEvent {
  event: ActionEvent_EventType;
  /**
   * Extra key:value data to provide more context to events
   * this will vary by EventType
   * example: `{ "oldTitle": "Foo", "newTitle":"Bar" }`
   */
  extraParams: { [key: string]: string };
}

export enum ActionEvent_EventType {
  CREATED_GROUP = 0,
  UPDATED_GROUP = 1,
  JOIN_GROUP = 2,
  LEAVE_GROUP = 3,
  ADD_MEMBER = 4,
  REMOVE_MEMBER = 5,
  EDIT_TITLE = 6,
  UNRECOGNIZED = -1,
}

export interface ActionEvent_ExtraParamsEntry {
  key: string;
  value: string;
}

export interface TypingMessage {
  senderId: string;
  type: TypingMessage_ActionType;
}

export enum TypingMessage_ActionType {
  CANCEL_TYPING = 0,
  TYPING = 1,
  RECORDING_AUDIO = 2,
  UPLOADING_AUDIO = 3,
  UPLOADING_PHOTO = 4,
  UNRECOGNIZED = -1,
}

export interface CreateConversationRequest {
  title: string;
  type: Conversation_ConversationType;
  memberIds: string[];
  image: Uint8Array;
}

export interface CreateConversationResponse {
  conversation: Conversation | undefined;
}

export interface GetConversationsRequest {
  limit: number;
  offset: number;
  id: string[];
}

export interface GetConversationsResponse {
  conversations: Conversation[];
}

export interface LeaveConversationRequest {
  conversationId: string;
}

export interface MemberRequest {
  conversationId: string;
  userId: string;
}

export interface ChatMessage {
  conversationId: string;
  typing?: TypingMessage | undefined;
  message?: Message | undefined;
}

export interface GetHistoryRequest {
  limit: number;
  offset: number;
  conversationId: string;
}

export interface GetHistoryResponse {
  messages: Message[];
}

export interface ReadHistoryRequest {
  lastMessageId: string;
  conversationId: string;
}

export const MESSENGER_PACKAGE_NAME = "messenger";

export interface MessageServiceClient {
  /** Return the history of messages for an conversation in DESC order. */

  getHistory(request: GetHistoryRequest): Observable<GetHistoryResponse>;

  /** Notifies the reading of messages from a channel or a user. */

  readHistory(request: ReadHistoryRequest): Observable<Empty>;

  /** Send and receive Messages or events to/from conversations. */

  communicate(request: Observable<ChatMessage>): Observable<ChatMessage>;
}

export interface MessageServiceController {
  /** Return the history of messages for an conversation in DESC order. */

  getHistory(
    request: GetHistoryRequest,
  ): Promise<GetHistoryResponse> | Observable<GetHistoryResponse> | GetHistoryResponse;

  /** Notifies the reading of messages from a channel or a user. */

  readHistory(request: ReadHistoryRequest): void;

  /** Send and receive Messages or events to/from conversations. */

  communicate(request: Observable<ChatMessage>): Observable<ChatMessage>;
}

export function MessageServiceControllerMethods() {
  return function (constructor: Function) {
    const grpcMethods: string[] = ["getHistory", "readHistory"];
    for (const method of grpcMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcMethod("MessageService", method)(constructor.prototype[method], method, descriptor);
    }
    const grpcStreamMethods: string[] = ["communicate"];
    for (const method of grpcStreamMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcStreamMethod("MessageService", method)(constructor.prototype[method], method, descriptor);
    }
  };
}

export const MESSAGE_SERVICE_NAME = "MessageService";

export interface ConversationServiceClient {
  /** Return conversations ordered by last update date */

  get(request: GetConversationsRequest): Observable<GetConversationsResponse>;

  create(request: CreateConversationRequest): Observable<CreateConversationResponse>;

  leave(request: LeaveConversationRequest): Observable<Empty>;

  addMember(request: MemberRequest): Observable<Empty>;

  removeMember(request: MemberRequest): Observable<Empty>;
}

export interface ConversationServiceController {
  /** Return conversations ordered by last update date */

  get(
    request: GetConversationsRequest,
  ): Promise<GetConversationsResponse> | Observable<GetConversationsResponse> | GetConversationsResponse;

  create(
    request: CreateConversationRequest,
  ): Promise<CreateConversationResponse> | Observable<CreateConversationResponse> | CreateConversationResponse;

  leave(request: LeaveConversationRequest): void;

  addMember(request: MemberRequest): void;

  removeMember(request: MemberRequest): void;
}

export function ConversationServiceControllerMethods() {
  return function (constructor: Function) {
    const grpcMethods: string[] = ["get", "create", "leave", "addMember", "removeMember"];
    for (const method of grpcMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcMethod("ConversationService", method)(constructor.prototype[method], method, descriptor);
    }
    const grpcStreamMethods: string[] = [];
    for (const method of grpcStreamMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcStreamMethod("ConversationService", method)(constructor.prototype[method], method, descriptor);
    }
  };
}

export const CONVERSATION_SERVICE_NAME = "ConversationService";
