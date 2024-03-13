/* eslint-disable */
import { GrpcMethod, GrpcStreamMethod } from "@nestjs/microservices";
import { Observable } from "rxjs";
import { Timestamp } from "../../../../google/protobuf/timestamp";

const protobufPackage = "calendar";

interface Empty {
}

/** Создание события в календаре */
export interface EventCreate {
  /** ID пользователя */
  userId: string;
  /** Название события */
  name: string;
  /** Описание события */
  description: string;
  /** Время начала */
  startTime:
    | Timestamp
    | undefined;
  /** // Сколько по времени идет событие */
  duration: number;
  /** ID категории */
  categoryId: string;
  /** ID приоритета */
  priorityId: string;
}

/** Обновление события в календаре */
export interface EventUpdate {
  /** ID события */
  id: string;
  /** ID пользователя */
  userId: string;
  /** Название события */
  name: string;
  /** Описание события */
  description: string;
  /** Время начала */
  startTime:
    | Timestamp
    | undefined;
  /** // Сколько по времени идет событие */
  duration: number;
  /** ID категории */
  categoryId: string;
  /** ID приоритета */
  priorityId: string;
}

/** Ответ с одним событием */
export interface EventResponse {
  /** ID события */
  id: string;
  /** ID пользователя */
  userId: string;
  /** Название события */
  name: string;
  /** Описание события */
  description: string;
  /** Время начала */
  startTime: string;
  /** Время окончания */
  endTime: string;
  /** ID категории */
  categoryId: string;
  /** ID приоритета */
  priorityId: string;
}

/** Id события - для удаления */
export interface EventDelete {
  id: string;
}

/** Пользователь */
export interface EventUser {
  userId: string;
  currentTime: Timestamp | undefined;
}

/** Успешный ответ - список категорий */
export interface ResponseCategories {
  categories: Category[];
}

/** Успешный ответ - список приоритетов */
export interface ResponsePriorities {
  priorities: Priority[];
}

/** Категория */
export interface Category {
  id: string;
  name: string;
}

/** Приоритет */
export interface Priority {
  id: string;
  name: string;
}

/** Успешный ответ - список Cобытий */
export interface ResponseEvents {
  result: EventResponse[];
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

  getCategories(request: Empty): Observable<ResponseCategories>;

  getPriorities(request: Empty): Observable<ResponsePriorities>;
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

  getCategories(request: Empty): Promise<ResponseCategories> | Observable<ResponseCategories> | ResponseCategories;

  getPriorities(request: Empty): Promise<ResponsePriorities> | Observable<ResponsePriorities> | ResponsePriorities;
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
      "getCategories",
      "getPriorities",
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
