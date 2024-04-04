import { NestFactory } from '@nestjs/core';
import { MessengerModule } from './messengerModule';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { join } from 'path';
import { MESSENGER_PACKAGE_NAME } from '@app/common/types/messenger';

async function bootstrap() {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    MessengerModule,
    {
      transport: Transport.GRPC,
      options: {
        package: MESSENGER_PACKAGE_NAME,
        protoPath: join(__dirname, '../messenger.proto'),
        url: 'localhost:5007',
      },
    },
  );
  await app.listen();
}
bootstrap();
