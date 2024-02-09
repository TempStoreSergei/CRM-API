import { NestFactory } from '@nestjs/core';
import { ApigetewayModule } from './apigeteway.module';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import * as cookieParser from 'cookie-parser';

async function bootstrap() {
  const app = await NestFactory.create(ApigetewayModule);

  app.enableCors();
  app.use(cookieParser());

  const config = new DocumentBuilder()
    .setTitle('CRM public endpoint')
    .setDescription('The CRM API description')
    .setVersion('1.0')
    .build();
  const document = SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('api', app, document);

  await app.listen(3000);
}
bootstrap();
