import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Put,
} from '@nestjs/common';
import { CalendarService } from './calendar.service';
import { ApiResponse, ApiTags } from '@nestjs/swagger';
import { Observable } from 'rxjs';
import {
  ResponseCategories,
  ResponseEvents,
  ResponseOK,
  ResponsePriorities,
} from '@app/common';
import {
  EventCreateDto,
  EventDeleteDto,
  EventUpdateDto,
  ResponseEventsDto,
  ResponseOKDto,
} from './dto';

@ApiTags('Events')
@Controller('events')
export class CalendarController {
  constructor(private readonly calendarService: CalendarService) {}

  @Get()
  @ApiResponse({
    status: 200,
    description: 'Get all events',
    type: ResponseEventsDto,
  })
  getAll(): Observable<ResponseEvents> {
    return this.calendarService.getAll();
  }

  @Get('categories')
  @ApiResponse({
    status: 200,
    description: 'Get all categories',
    type: ResponseEventsDto,
  })
  getCategories(): Observable<ResponseCategories> {
    return this.calendarService.getCategories({});
  }

  @Get('priorities')
  @ApiResponse({
    status: 200,
    description: 'Get all priorities',
    type: ResponseEventsDto,
  })
  getPriorities(): Observable<ResponsePriorities> {
    return this.calendarService.getPriorities({});
  }

  @Get('day/:user/:timestamp')
  @ApiResponse({
    status: 200,
    description: 'Get events for a specific day',
    type: ResponseEventsDto,
  })
  getDayEvent(
    @Param('user') request: string,
    @Param('timestamp') timestamp: number,
  ): Observable<ResponseEvents> {
    const currentTime = {
      seconds: timestamp,
      nanos: 0,
    };
    return this.calendarService.getDayEvent({ userId: request, currentTime });
  }

  @Get('week/:user/:timestamp')
  @ApiResponse({
    status: 200,
    description: 'Get events for a specific week',
    type: ResponseEventsDto,
  })
  getWeekEvent(
    @Param('user') request: string,
    @Param('timestamp') timestamp: number,
  ): Observable<ResponseEvents> {
    const currentTime = {
      seconds: timestamp,
      nanos: 0,
    };

    return this.calendarService.getWeekEvent({
      userId: request,
      currentTime,
    });
  }

  @Get('month/:user/:timestamp')
  @ApiResponse({
    status: 200,
    description: 'Get events for a specific month',
    type: ResponseEventsDto,
  })
  getMonthEvent(
    @Param('user') request: string,
    @Param('timestamp') timestamp: number,
  ): Observable<ResponseEvents> {
    const currentTime = {
      seconds: timestamp,
      nanos: 0,
    };

    return this.calendarService.getMonthEvent({ userId: request, currentTime });
  }

  @Post()
  @ApiResponse({
    status: 201,
    description: 'Event created successfully',
    type: ResponseOKDto,
  })
  create(@Body() request: EventCreateDto): Observable<ResponseOK> {
    return this.calendarService.create(request);
  }

  @Put()
  @ApiResponse({
    status: 200,
    description: 'Event updated successfully',
    type: ResponseOKDto,
  })
  update(@Body() request: EventUpdateDto): Observable<ResponseOK> {
    return this.calendarService.update(request);
  }

  @Delete()
  @ApiResponse({
    status: 200,
    description: 'Event deleted successfully',
    type: ResponseOKDto,
  })
  delete(@Body() request: EventDeleteDto): Observable<ResponseOK> {
    return this.calendarService.delete(request);
  }
}
