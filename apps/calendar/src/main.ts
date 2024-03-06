import { NestFactory } from '@nestjs/core';
import { CalendarModule } from './calendar.module';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { join } from 'path';
import { CALENDAR_PACKAGE_NAME } from '@app/common';

async function bootstrap() {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    CalendarModule,
    {
      transport: Transport.GRPC,
      options: {
        package: CALENDAR_PACKAGE_NAME,
        protoPath: join(__dirname, '../calendar.proto'),
        url: 'localhost:5004',
      },
    },
  );
  await app.listen();
}
bootstrap();
