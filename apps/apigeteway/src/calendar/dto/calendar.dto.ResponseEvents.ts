import { ApiProperty } from '@nestjs/swagger';
import { ResponseEvents, Event } from '@app/common';

export class ResponseEventsDto implements ResponseEvents {
  @ApiProperty()
  result: Event[];
}
