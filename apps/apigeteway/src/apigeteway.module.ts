import { Module } from '@nestjs/common';
import { UsersModule } from './users/users.module';
import { AuthModule } from './auth/auth.module';
import { CalendarModule } from './events/calendar.module';

@Module({
  imports: [UsersModule, AuthModule, CalendarModule],
  controllers: [],
  providers: [],
})
export class ApigetewayModule {}
