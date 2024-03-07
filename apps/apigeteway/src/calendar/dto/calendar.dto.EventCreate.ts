import { ApiProperty } from '@nestjs/swagger';
import { EventCreate } from '@app/common';

export class EventCreateDto implements EventCreate {
  @ApiProperty()
  duration: string;

  @ApiProperty({ format: 'date-time' })
  eventDate: Date;

  @ApiProperty({ format: 'date-time' })
  startTime: {
    seconds: number;
    nanos: number;
  };

  @ApiProperty()
  body: string;

  @ApiProperty()
  user: string;

  @ApiProperty()
  title: string;
}
