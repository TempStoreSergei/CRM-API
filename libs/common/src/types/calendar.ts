/* eslint-disable */
import { GrpcMethod, GrpcStreamMethod } from "@nestjs/microservices";
import { Observable } from "rxjs";

export const protobufPackage = "calendar";

export interface Empty {
}

/** Создание события в календаре */
export interface EventCreate {
  user: string;
  title: string;
  body: string;
  starttime: string;
  duration: string;
}

/** Обновление события в календаре */
export interface EventUpdate {
  id: string;
  user: string;
  title: string;
  body: string;
  starttime: string;
  duration: string;
}

/** Id события - для удаления */
export interface EventDelete {
  id: string;
}

/** Пользователь */
export interface EventUser {
  user: string;
}

/** Cобытие в календаре */
export interface Event {
  id: string;
  user: string;
  title: string;
  body: string;
  starttime: string;
  endtime: string;
}

/** Успешный ответ - список Cобытий */
export interface ResponseEvents {
  result: Event[];
}

/** Успешный ответ от сервиса */
export interface ResponseOK {
  result: string;
}

export const CALENDAR_PACKAGE_NAME = "calendar";

/** Методы Сервиса */

export interface CalendarServiceClient {
  create(request: EventCreate): Observable<ResponseOK>;

  update(request: EventUpdate): Observable<ResponseOK>;

  delete(request: EventDelete): Observable<ResponseOK>;

  getAll(request: Empty): Observable<ResponseEvents>;

  getDayEvent(request: EventUser): Observable<ResponseEvents>;

  getWeekEvent(request: EventUser): Observable<ResponseEvents>;

  getMonthEvent(request: EventUser): Observable<ResponseEvents>;
}

/** Методы Сервиса */

export interface CalendarServiceController {
  create(request: EventCreate): Promise<ResponseOK> | Observable<ResponseOK> | ResponseOK;

  update(request: EventUpdate): Promise<ResponseOK> | Observable<ResponseOK> | ResponseOK;

  delete(request: EventDelete): Promise<ResponseOK> | Observable<ResponseOK> | ResponseOK;

  getAll(request: Empty): Promise<ResponseEvents> | Observable<ResponseEvents> | ResponseEvents;

  getDayEvent(request: EventUser): Promise<ResponseEvents> | Observable<ResponseEvents> | ResponseEvents;

  getWeekEvent(request: EventUser): Promise<ResponseEvents> | Observable<ResponseEvents> | ResponseEvents;

  getMonthEvent(request: EventUser): Promise<ResponseEvents> | Observable<ResponseEvents> | ResponseEvents;
}

export function CalendarServiceControllerMethods() {
  return function (constructor: Function) {
    const grpcMethods: string[] = [
      "create",
      "update",
      "delete",
      "getAll",
      "getDayEvent",
      "getWeekEvent",
      "getMonthEvent",
    ];
    for (const method of grpcMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcMethod("CalendarService", method)(constructor.prototype[method], method, descriptor);
    }
    const grpcStreamMethods: string[] = [];
    for (const method of grpcStreamMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcStreamMethod("CalendarService", method)(constructor.prototype[method], method, descriptor);
    }
  };
}

export const CALENDAR_SERVICE_NAME = "CalendarService";
