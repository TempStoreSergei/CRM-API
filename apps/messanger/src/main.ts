import { NestFactory } from '@nestjs/core';
import { MessangerModule } from './messanger.module';

async function bootstrap() {
  const app = await NestFactory.create(MessangerModule);
  await app.listen(3004);
}
bootstrap();
