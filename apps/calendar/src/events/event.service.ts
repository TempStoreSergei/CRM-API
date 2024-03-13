import { Inject, Injectable, Logger } from '@nestjs/common';
import {
  addMinutes,
  fromUnixTime,
  formatISO,
  startOfWeek,
  endOfWeek,
  startOfMonth,
  endOfMonth,
  startOfDay,
  endOfDay,
} from 'date-fns';
import { PrismaService } from '../prisma/prisma.service';
import { CACHE_MANAGER, Cache } from '@nestjs/cache-manager';

import {
  EventCreate,
  EventDelete,
  EventUpdate,
  EventUser,
  ResponseEvents,
  ResponseOK,
  ResponseCategories,
  ResponsePriorities,
} from '@app/common';

@Injectable()
export class EventService {
  private readonly logger = new Logger(EventService.name);

  constructor(
    private readonly prisma: PrismaService,
    @Inject(CACHE_MANAGER) private cacheManager: Cache,
  ) {}

  async create(eventCreate: EventCreate): Promise<ResponseOK> {
    const startTimeDate = fromUnixTime(eventCreate.startTime.seconds);
    const endTime = addMinutes(startTimeDate, eventCreate.duration);
    try {
      const newEvent = await this.prisma.event.create({
        data: {
          userId: eventCreate.userId,
          name: eventCreate.name,
          description: eventCreate.description,
          startTime: startTimeDate,
          endTime: endTime,
          categoryId: eventCreate.categoryId,
          priorityId: eventCreate.priorityId,
        },
      });
      this.logger.log(`Event created: ${newEvent.id}`);
      return { result: newEvent.id };
    } catch (error) {
      this.logger.error(`Error creating event: ${error.message}`);
      throw error;
    }
  }

  async delete(eventDelete: EventDelete): Promise<ResponseOK> {
    try {
      await this.prisma.event.delete({
        where: { id: eventDelete.id },
      });
      this.logger.log(`Event deleted: ${eventDelete.id}`);
      return { result: 'ok' };
    } catch (error) {
      this.logger.error(`Error deleting event: ${error.message}`);
      throw error;
    }
  }

  async update(eventUpdate: EventUpdate): Promise<ResponseOK> {
    const startTimeDate = fromUnixTime(eventUpdate.startTime.seconds);
    const endTime = addMinutes(startTimeDate, eventUpdate.duration);
    try {
      await this.prisma.event.update({
        where: { id: eventUpdate.id },
        data: {
          userId: eventUpdate.userId,
          name: eventUpdate.name,
          description: eventUpdate.description,
          startTime: startTimeDate,
          endTime: endTime,
          categoryId: eventUpdate.categoryId,
          priorityId: eventUpdate.priorityId,
        },
      });
      this.logger.log(`Event updated: ${eventUpdate.id}`);
      return { result: 'ok' };
    } catch (error) {
      this.logger.error(`Error updating event: ${error.message}`);
      throw error;
    }
  }

  async getAll(): Promise<ResponseEvents> {
    try {
      const events = await this.prisma.event.findMany();
      return this.mapEventsToResponse(events);
    } catch (error) {
      this.logger.error(`Error retrieving events: ${error.message}`);
      throw error;
    }
  }

  async getDayEvent(eventUser: EventUser): Promise<ResponseEvents> {
    if (!eventUser.currentTime) {
      throw new Error('Current time is required.');
    }

    const referenceTime = fromUnixTime(eventUser.currentTime.seconds);

    // Setting the start and end of the day based on the reference time
    const todayStart = startOfDay(referenceTime);
    const todayEnd = endOfDay(referenceTime);
    try {
      const events = await this.prisma.event.findMany({
        where: {
          userId: eventUser.userId,
          startTime: {
            gte: todayStart,
            lte: todayEnd,
          },
        },
      });
      return this.mapEventsToResponse(events);
    } catch (error) {
      this.logger.error(`Error retrieving day events: ${error.message}`);
      throw error;
    }
  }

  async getWeekEvent(eventUser: EventUser): Promise<ResponseEvents> {
    if (!eventUser.currentTime) {
      throw new Error('Current time is required.');
    }

    const referenceTime = fromUnixTime(eventUser.currentTime.seconds);

    const weekStart = startOfWeek(referenceTime);
    const weekEnd = endOfWeek(referenceTime);
    try {
      const events = await this.prisma.event.findMany({
        where: {
          userId: eventUser.userId,
          startTime: {
            gte: weekStart,
            lte: weekEnd,
          },
        },
      });
      return this.mapEventsToResponse(events);
    } catch (error) {
      this.logger.error(`Error retrieving week events: ${error.message}`);
      throw error;
    }
  }

  async getMonthEvent(eventUser: EventUser): Promise<ResponseEvents> {
    if (!eventUser.currentTime) {
      throw new Error('Current time is required.');
    }

    const referenceTime = fromUnixTime(eventUser.currentTime.seconds);
    const monthStart = startOfMonth(referenceTime);
    const monthEnd = endOfMonth(referenceTime);
    try {
      const events = await this.prisma.event.findMany({
        where: {
          userId: eventUser.userId,
          startTime: {
            gte: monthStart,
            lte: monthEnd,
          },
        },
      });
      return this.mapEventsToResponse(events);
    } catch (error) {
      this.logger.error(`Error retrieving month events: ${error.message}`);
      throw error;
    }
  }

  async getCategories(): Promise<ResponseCategories> {
    try {
      const categories = await this.prisma.category.findMany();
      const result = categories.map((category) => ({
        id: category.id,
        name: category.name,
      }));
      return { categories: result }; // Note: Make sure this matches the structure expected by your protobuf definition.
    } catch (error) {
      this.logger.error(`Error retrieving categories: ${error.message}`);
      throw error;
    }
  }

  async getPriorities(): Promise<ResponsePriorities> {
    try {
      const priorities = await this.prisma.priority.findMany();
      const result = priorities.map((priority) => ({
        id: priority.id,
        name: priority.name,
      }));
      return { priorities: result };
    } catch (error) {
      this.logger.error(`Error retrieving priorities: ${error.message}`);
      throw error;
    }
  }

  private mapEventsToResponse(events) {
    const result = events.map((event) => ({
      id: event.id,
      userId: event.userId,
      name: event.name,
      description: event.description,
      startTime: formatISO(event.startTime),
      endTime: formatISO(event.endTime),
      categoryId: event.categoryId,
      priorityId: event.priorityId,
    }));
    return { result };
  }
}
