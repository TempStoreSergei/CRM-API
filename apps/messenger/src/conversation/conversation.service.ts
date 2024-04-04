import { Injectable } from '@nestjs/common';
import { Metadata, ServerUnaryCall } from '@grpc/grpc-js';
import { ConversationManager } from '../cassandra/conversationManager';
import { v4 as uuidv4 } from 'uuid';

@Injectable()
export class ConversationService {
  constructor(private readonly conversationsService: ConversationManager) {}

  async getConversations(
    call: ServerUnaryCall<any, any>,
    metadata: Metadata,
  ): Promise<any> {
    const authToken = call.request.token; // Assuming the token is part of the request
    const userId = authToken;
    const { limit, offset } = call.request;
    const conversations = await this.conversationsService.getByUserID(
      userId,
      limit,
      offset,
    );
    return { conversations };
  }

  async createConversation(
    call: ServerUnaryCall<any, any>,
    metadata: Metadata,
  ): Promise<any> {
    const { title, type, memberIds, image } = call;
    const userId = 'serega-to';
    const conversation = {
      id: uuidv4(),
      title,
      type,
      memberIds: [...memberIds, userId], // Add the creator as a member
      image,
      creationTime: new Date().toISOString(),
      updateTime: new Date().toISOString(),
    };
    await this.conversationsService.create(conversation);
    console.log('Created conversation:', conversation);
    return { conversation };
  }

  async addMember(
    call: ServerUnaryCall<any, any>,
    metadata: Metadata,
  ): Promise<any> {
    const { conversationId, userId } = call.request;
    const authToken = metadata.get('authorization')[0] as string;
    const requestorId = authToken;

    await this.conversationsService.addMember(conversationId, userId);
    return {};
  }

  async removeMember(
    call: ServerUnaryCall<any, any>,
    metadata: Metadata,
  ): Promise<any> {
    const { conversationId, userId } = call.request;
    const authToken = metadata.get('authorization')[0] as string;
    const requestorId = authToken;

    await this.conversationsService.removeMember(conversationId, userId);
    return {};
  }

  async createSchema(): Promise<void> {
    await this.conversationsService.createSchema('conversations_new');
  }
}
