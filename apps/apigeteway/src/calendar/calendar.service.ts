import { Inject, Injectable, OnModuleInit } from '@nestjs/common';
import {
  CALENDAR_PACKAGE_NAME,
  CALENDAR_SERVICE_NAME,
  CalendarServiceClient,
} from '@app/common';
import { ClientGrpc } from '@nestjs/microservices';

@Injectable()
export class CalendarService implements OnModuleInit {
  private calendarService: CalendarServiceClient;

  constructor(@Inject(CALENDAR_PACKAGE_NAME) private client: ClientGrpc) {}
  onModuleInit() {
    this.calendarService =
      this.client.getService<CalendarServiceClient>(CALENDAR_SERVICE_NAME);
  }
  findAll() {
    return this.calendarService.getAll({});
  }
}
