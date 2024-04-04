import { Module } from '@nestjs/common';
import { MessengerModuleMessage } from './message/message.module';
import { MessengerModuleConversation } from './conversation/conversation.module';

@Module({
  imports: [MessengerModuleMessage, MessengerModuleConversation],
  controllers: [],
  providers: [],
})
export class MessengerModule {}
