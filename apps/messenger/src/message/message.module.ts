import { Module } from '@nestjs/common';
import { MessengerController } from './message.controller';
import { MessageManager } from '../cassandra/messageManager';
import { ConversationManager } from '../cassandra/conversationManager';

@Module({
  providers: [MessageManager, ConversationManager],
  controllers: [MessengerController],
  imports: [],
})
export class MessengerModuleMessage {}
