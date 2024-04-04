import { Controller, OnModuleInit } from '@nestjs/common';
import { Observable } from 'rxjs';
import { ConversationService } from './conversation.service';
import {
  ConversationServiceController,
  CreateConversationRequest,
  CreateConversationResponse,
  GetConversationsRequest,
  GetConversationsResponse,
  LeaveConversationRequest,
  MemberRequest,
} from '@app/common/types/messenger';
import { GrpcMethod } from '@nestjs/microservices'

@Controller()
export class ConversationController
  implements ConversationServiceController, OnModuleInit
{
  constructor(private readonly conversationService: ConversationService) {}

  onModuleInit(): any {
    this.conversationService.createSchema();
  }

  @GrpcMethod('ConversationService', 'addMember')
  addMember(request: MemberRequest) {
    return this.conversationService.addMember(request);
  }

  @GrpcMethod('ConversationService', 'removeMember')
  removeMember(request: MemberRequest) {
    return this.conversationService.removeMember(request);
  }

  @GrpcMethod('ConversationService', 'create')
  create(
    request: CreateConversationRequest,
  ):
    | Promise<CreateConversationResponse>
    | Observable<CreateConversationResponse>
    | CreateConversationResponse {
    return this.conversationService.createConversation(request);
  }

  @GrpcMethod('ConversationService', 'get')
  get(
    request: GetConversationsRequest,
  ):
    | Promise<GetConversationsResponse>
    | Observable<GetConversationsResponse>
    | GetConversationsResponse {
    return this.conversationService.getConversations(request);
  }

  @GrpcMethod('ConversationService', 'leave')
  leave(request: LeaveConversationRequest) {
    return 'Not implemented';
  }
}
