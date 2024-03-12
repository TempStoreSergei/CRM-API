import { Inject, Injectable, Logger } from '@nestjs/common';
import {
  Event,
  EventCreate,
  EventDelete,
  EventUpdate,
  EventUser,
  ResponseEvents,
  ResponseOK,
} from '@app/common';
import { addMinutes, endOfDay, fromUnixTime, startOfDay } from 'date-fns';
import { PrismaService } from '../prisma/prisma.service';
import { CACHE_MANAGER, Cache } from '@nestjs/cache-manager';

@Injectable()
export class EventService {
  private readonly logger = new Logger(EventService.name);

  constructor(
    private readonly prisma: PrismaService,
    @Inject(CACHE_MANAGER) private cacheManager: Cache,
  ) {}

  async create(eventCreate: EventCreate): Promise<ResponseOK> {
    const startTimeDate = fromUnixTime(eventCreate.startTime.seconds);

    const endTime = addMinutes(startTimeDate, Number(eventCreate.duration));
    const newEvent = await this.prisma.event.create({
      data: {
        userId: eventCreate.user,
        name: eventCreate.title,
        description: eventCreate.body,
        startTime: startTimeDate,
        endTime: endTime,
        categoryId: '1',
        priorityId: '1',
      },
    });

    return { result: newEvent.id };
  }

  async delete(eventDelete: EventDelete): Promise<ResponseOK> {
    const eventId = eventDelete.id;

    try {
      await this.prisma.event.delete({
        where: { id: eventId },
      });

      return { result: 'ok' };
    } catch (error) {
      throw new Error(`Failed to delete event: ${error.message}`);
    }
  }

  async update(eventUpdate: EventUpdate): Promise<ResponseOK> {
    const eventId = eventUpdate.id;
    const endTime = addMinutes(
      eventUpdate.startTime,
      Number(eventUpdate.duration) * 60000,
    );

    try {
      await this.prisma.event.update({
        where: { id: eventId },
        data: {
          title: eventUpdate.title,
          body: eventUpdate.body,
          startTime: eventUpdate.startTime,
          endTime,
        },
      });

      return { result: 'ok' };
    } catch (error) {
      throw new Error(`Failed to update event: ${error.message}`);
    }
  }

  async getAll(): Promise<ResponseEvents> {
    const events: Event[] = await this.prisma.event.findMany();
    return this.mapEventsToResponse(events);
  }

  async getDayEvent(user: EventUser): Promise<ResponseEvents> {
    const todayStart = startOfDay(new Date());
    const todayEnd = endOfDay(new Date());

    const events: Event[] = await this.prisma.event.findMany({
      where: {
        userId: user.user,
        startTime: {
          gte: todayStart,
          lt: todayEnd,
        },
      },
    });

    return this.mapEventsToResponse(events);
  }

  async getWeekEvent(user: EventUser): Promise<ResponseEvents> {
    const today = new Date();
    const startOfWeek = new Date(today);
    startOfWeek.setDate(
      today.getDate() - today.getDay() + (today.getDay() === 0 ? -6 : 1),
    );

    const endOfWeek = new Date(startOfWeek);
    endOfWeek.setDate(startOfWeek.getDate() + 6);

    const events: Event[] = await this.prisma.event.findMany({
      where: {
        userId: user.user,
        startTime: {
          gte: startOfWeek,
          lt: endOfWeek,
        },
      },
    });

    return this.mapEventsToResponse(events);
  }

  async getMonthEvent(user: EventUser): Promise<ResponseEvents> {
    const today = new Date();
    const startOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);

    const lastDayOfMonth = new Date(
      today.getFullYear(),
      today.getMonth() + 1,
      0,
    );
    const endOfMonth = new Date(lastDayOfMonth);
    endOfMonth.setDate(lastDayOfMonth.getDate() + 1);

    const events: Event[] = await this.prisma.event.findMany({
      where: {
        userId: user.user,
        startTime: {
          gte: startOfMonth,
          lt: endOfMonth,
        },
      },
    });

    return this.mapEventsToResponse(events);
  }

  private mapEventsToResponse(events: Event[]): ResponseEvents {
    const result = events.map((e) => ({
      id: e.id,
      user: e.user,
      title: e.title,
      body: e.body,
      startTime: e.startTime,
      endTime: e.endTime,
    }));
    // @ts-ignore
    return { result };
  }
}
