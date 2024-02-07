import { Injectable } from '@nestjs/common';

@Injectable()
export class ApigetewayService {
  getHello(): string {
    return 'Hello World!';
  }
}
