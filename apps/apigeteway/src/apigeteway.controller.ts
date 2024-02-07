import { Controller, Get } from '@nestjs/common';
import { ApigetewayService } from './apigeteway.service';

@Controller()
export class ApigetewayController {
  constructor(private readonly apigetewayService: ApigetewayService) {}

  @Get()
  getHello(): string {
    return this.apigetewayService.getHello();
  }
}
