import { Inject, Injectable, OnModuleInit } from '@nestjs/common';
import {
  CALENDAR_PACKAGE_NAME,
  CALENDAR_SERVICE_NAME,
  CalendarServiceClient,
  Empty,
  EventCreate,
  EventDelete,
  EventUpdate,
  EventUser,
  ResponseCategories,
  ResponseEvents,
  ResponseOK,
  ResponsePriorities,
} from '@app/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Observable } from 'rxjs';

@Injectable()
export class CalendarService implements OnModuleInit {
  private calendarService: CalendarServiceClient;

  constructor(@Inject(CALENDAR_PACKAGE_NAME) private client: ClientGrpc) {}
  onModuleInit() {
    this.calendarService = this.client.getService<CalendarServiceClient>(
      CALENDAR_SERVICE_NAME,
    );
  }

  getAll(): Observable<ResponseEvents> {
    return this.calendarService.getAll({});
  }

  getDayEvent(request: EventUser): Observable<ResponseEvents> {
    return this.calendarService.getDayEvent(request);
  }

  getWeekEvent(request: EventUser): Observable<ResponseEvents> {
    return this.calendarService.getWeekEvent(request);
  }

  getMonthEvent(request: EventUser): Observable<ResponseEvents> {
    return this.calendarService.getMonthEvent(request);
  }

  create(request: EventCreate): Observable<ResponseOK> {
    return this.calendarService.create(request);
  }

  update(request: EventUpdate): Observable<ResponseOK> {
    return this.calendarService.update(request);
  }

  delete(request: EventDelete): Observable<ResponseOK> {
    return this.calendarService.delete(request);
  }

  getPriorities(request: Empty): Observable<ResponsePriorities> {
    return this.calendarService.getPriorities(request);
  }

  getCategories(request: Empty): Observable<ResponseCategories> {
    return this.calendarService.getCategories(request);
  }
}
