import { Module } from '@nestjs/common';
import { EventModule } from './events/event.module';
@Module({
  imports: [EventModule],
  controllers: [],
  providers: [],
})
export class CalendarModule {}
