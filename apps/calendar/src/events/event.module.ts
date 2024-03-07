import { Module } from '@nestjs/common';
import { EventService } from './event.service';
import { EventController } from './event.controller';
import { CacheModule } from '@nestjs/cache-manager';
import { PrismaModule } from '../prisma/prisma.module';

@Module({
  controllers: [EventController],
  imports: [CacheModule.register(), PrismaModule],
  providers: [EventService],
})
export class EventModule {}
