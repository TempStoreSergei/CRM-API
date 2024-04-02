import { NestFactory } from '@nestjs/core';
import { MessangerModule } from './messanger.module';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { join } from 'path';
import { MESSENGER_PACKAGE_NAME } from '@app/common/types/messanegr';

async function bootstrap() {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    MessangerModule,
    {
      transport: Transport.GRPC,
      options: {
        package: MESSENGER_PACKAGE_NAME,
        protoPath: join(__dirname, '../messanegr.proto'),
        url: 'localhost:5007',
      },
    },
  );
  await app.listen();
}
bootstrap();
