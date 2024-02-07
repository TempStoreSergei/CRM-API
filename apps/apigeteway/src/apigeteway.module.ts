import { Module } from '@nestjs/common';
import { ApigetewayController } from './apigeteway.controller';
import { ApigetewayService } from './apigeteway.service';

@Module({
  imports: [],
  controllers: [ApigetewayController],
  providers: [ApigetewayService],
})
export class ApigetewayModule {}
