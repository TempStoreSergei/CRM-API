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

@Controller()
export class ConversationController
  implements ConversationServiceController, OnModuleInit
{
  constructor(private readonly conversationService: ConversationService) {}

  onModuleInit(): any {
    this.conversationService.createSchema();
  }

  addMember(request: MemberRequest) {
    return this.conversationService.addMember(request);
  }

  removeMember(request: MemberRequest) {
    return this.conversationService.removeMember(request);
  }

  create(
    request: CreateConversationRequest,
  ):
    | Promise<CreateConversationResponse>
    | Observable<CreateConversationResponse>
    | CreateConversationResponse {
    return this.conversationService.createConversation(request);
  }

  get(
    request: GetConversationsRequest,
  ):
    | Promise<GetConversationsResponse>
    | Observable<GetConversationsResponse>
    | GetConversationsResponse {
    return this.conversationService.getConversations(request);
  }

  leave(request: LeaveConversationRequest) {
    return 'Not implemented';
  }
}
