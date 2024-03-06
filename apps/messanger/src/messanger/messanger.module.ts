import { Module } from '@nestjs/common';
import { MessangerService } from './messanger.service';
import { MessangerController } from './messanger.controller';
import { MessangerGateway } from './messanger.gateway';

@Module({
  providers: [MessangerService, MessangerGateway],
  controllers: [MessangerController],
})
export class MessangerModuleInner {}
