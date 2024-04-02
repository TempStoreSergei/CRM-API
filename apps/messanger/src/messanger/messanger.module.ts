import { Module } from '@nestjs/common';
import { MessangerService } from './messanger.service';
import { MessangerController } from './messanger.controller';
import { PrismaModule } from '../prisma/prisma.module';

@Module({
  providers: [MessangerService],
  controllers: [MessangerController],
  imports: [PrismaModule],
})
export class MessangerModuleInner {}
