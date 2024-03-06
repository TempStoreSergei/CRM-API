import { Module } from '@nestjs/common';
import { MessangerModuleInner } from './messanger/messanger.module';

@Module({
  imports: [MessangerModuleInner],
  controllers: [],
  providers: [],
})
export class MessangerModule {}
