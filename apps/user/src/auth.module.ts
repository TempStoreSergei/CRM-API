import { Module } from '@nestjs/common';
import { AuthModuleInner } from './auth/auth.module';

@Module({
  imports: [AuthModuleInner],
  controllers: [],
  providers: [],
})
export class AuthModule {}
