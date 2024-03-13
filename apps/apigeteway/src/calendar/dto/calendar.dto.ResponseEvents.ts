import { ApiProperty } from '@nestjs/swagger';
import { ResponseEvents, EventResponse } from '@app/common';

export class ResponseEventsDto implements ResponseEvents {
  @ApiProperty()
  result: EventResponse[];
}
