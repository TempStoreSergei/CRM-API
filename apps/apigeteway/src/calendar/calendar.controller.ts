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
import { ResponseEvents, ResponseOK } from '@app/common';
import {
  EventCreateDto,
  EventDeleteDto,
  EventUpdateDto,
  ResponseEventsDto,
  ResponseOKDto,
} from './dto';

@ApiTags('Calendar')
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

  @Get('day/:user')
  @ApiResponse({
    status: 200,
    description: 'Get events for a specific day',
    type: ResponseEventsDto,
  })
  getDayEvent(@Param('user') request: string): Observable<ResponseEvents> {
    return this.calendarService.getDayEvent({ user: request });
  }

  @Get('week/:user')
  @ApiResponse({
    status: 200,
    description: 'Get events for a specific week',
    type: ResponseEventsDto,
  })
  getWeekEvent(@Param('user') request: string): Observable<ResponseEvents> {
    return this.calendarService.getWeekEvent({ user: request });
  }

  @Get('month/:user')
  @ApiResponse({
    status: 200,
    description: 'Get events for a specific month',
    type: ResponseEventsDto,
  })
  getMonthEvent(@Param('user') request: string): Observable<ResponseEvents> {
    return this.calendarService.getMonthEvent({ user: request });
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
