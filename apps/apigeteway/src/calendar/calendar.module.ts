import { Module } from '@nestjs/common';
import { CalendarService } from './calendar.service';
import { CalendarController } from './calendar.controller';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { CALENDAR_SERVICE } from './constatnst';
import { CALENDAR_PACKAGE_NAME } from '@app/common';
import { join } from 'path';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: CALENDAR_SERVICE,
        transport: Transport.GRPC,
        options: {
          package: CALENDAR_PACKAGE_NAME,
          protoPath: join(__dirname, '../calendar.proto'),
          url: 'localhost:5004',
        },
      },
    ]),
  ],
  controllers: [CalendarController],
  providers: [CalendarService],
})
export class CalendarModule {}
