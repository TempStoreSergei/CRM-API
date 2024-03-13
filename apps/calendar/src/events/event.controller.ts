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
  ResponseCategories,
  ResponsePriorities,
} from '@app/common';
import { Observable } from 'rxjs';

@Controller()
@CalendarServiceControllerMethods()
export class EventController implements CalendarServiceController {
  constructor(private readonly usersService: EventService) {}

  getWeekEvent(
    request: EventUser,
  ): Promise<ResponseEvents> | Observable<ResponseEvents> | ResponseEvents {
    return this.usersService.getWeekEvent(request);
  }

  getDayEvent(
    request: EventUser,
  ): Promise<ResponseEvents> | Observable<ResponseEvents> | ResponseEvents {
    return this.usersService.getDayEvent(request);
  }

  getMonthEvent(
    request: EventUser,
  ): Promise<ResponseEvents> | Observable<ResponseEvents> | ResponseEvents {
    return this.usersService.getMonthEvent(request);
  }

  getAll(
    request: Empty,
  ): Promise<ResponseEvents> | Observable<ResponseEvents> | ResponseEvents {
    return this.usersService.getAll();
  }

  create(
    request: EventCreate,
  ): Promise<ResponseOK> | Observable<ResponseOK> | ResponseOK {
    return this.usersService.create(request);
  }

  delete(
    request: EventDelete,
  ): Promise<ResponseOK> | Observable<ResponseOK> | ResponseOK {
    return this.usersService.delete(request);
  }

  update(
    request: EventUpdate,
  ): Promise<ResponseOK> | Observable<ResponseOK> | ResponseOK {
    return this.usersService.update(request);
  }

  getCategories(
    request: Empty,
  ):
    | Promise<ResponseCategories>
    | Observable<ResponseCategories>
    | ResponseCategories {
    return this.usersService.getCategories();
  }

  getPriorities(
    request: Empty,
  ):
    | Promise<ResponsePriorities>
    | Observable<ResponsePriorities>
    | ResponsePriorities {
    return this.usersService.getPriorities();
  }
}
