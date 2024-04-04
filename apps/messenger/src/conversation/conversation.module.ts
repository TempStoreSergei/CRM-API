import { Module } from '@nestjs/common';
import { ConversationService } from './conversation.service';
import { ConversationController } from './conversation.controller';
import { ConversationManager } from '../cassandra/conversationManager';

@Module({
  providers: [ConversationService, ConversationManager],
  controllers: [ConversationController],
  imports: [],
})
export class MessengerModuleConversation {}
