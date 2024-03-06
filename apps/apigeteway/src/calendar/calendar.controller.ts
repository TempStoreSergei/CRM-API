import { Controller, Get } from '@nestjs/common';
import { CalendarService } from './calendar.service';
import { ApiTags } from '@nestjs/swagger';

@ApiTags('Calendar')
@Controller('events')
export class CalendarController {
  constructor(private readonly calendarService: CalendarService) {}

  @Get()
  findAll() {
    return this.calendarService.findAll();
  }
}
