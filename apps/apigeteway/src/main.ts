import { NestFactory } from '@nestjs/core';
import { ApigetewayModule } from './apigeteway.module';

async function bootstrap() {
  const app = await NestFactory.create(ApigetewayModule);
  await app.listen(3000);
}
bootstrap();
