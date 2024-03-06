import { Controller } from '@nestjs/common';
import { EventService } from './event.service';
import {
  CalendarServiceController,
  EventUser,
  ResponseEvents,
  Empty,
  EventCreate,
  ResponseOK,
  EventDelete,
  EventUpdate,
  CalendarServiceControllerMethods,
} from '@app/common';
import { Observable } from 'rxjs';

@Controller()
@CalendarServiceControllerMethods()
export class EventController implements CalendarServiceController {
  constructor(private readonly usersService: EventService) {}

  getWeekEvent(
    request: EventUser,
  ): Promise<ResponseEvents> | Observable<ResponseEvents> | ResponseEvents {
    return {
      result: [],
    };
  }

  getDayEvent(
    request: EventUser,
  ): Promise<ResponseEvents> | Observable<ResponseEvents> | ResponseEvents {
    return {
      result: [],
    };
  }

  getMonthEvent(
    request: EventUser,
  ): Promise<ResponseEvents> | Observable<ResponseEvents> | ResponseEvents {
    return {
      result: [],
    };
  }

  getAll(
    request: Empty,
  ): Promise<ResponseEvents> | Observable<ResponseEvents> | ResponseEvents {
    return this.usersService.findAll(request);
  }

  create(
    request: EventCreate,
  ): Promise<ResponseOK> | Observable<ResponseOK> | ResponseOK {
    return {
      result: 'Ok',
    };
  }

  delete(
    request: EventDelete,
  ): Promise<ResponseOK> | Observable<ResponseOK> | ResponseOK {
    return {
      result: 'Ok',
    };
  }

  update(
    request: EventUpdate,
  ): Promise<ResponseOK> | Observable<ResponseOK> | ResponseOK {
    return {
      result: 'Ok',
    };
  }
}
