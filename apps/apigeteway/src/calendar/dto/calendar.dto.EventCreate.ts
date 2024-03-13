import { ApiProperty } from '@nestjs/swagger';
import { EventCreate } from '@app/common';

export class EventCreateDto implements EventCreate {
  @ApiProperty()
  duration: number;

  @ApiProperty({ format: 'timestamp' })
  startTime: {
    seconds: number;
    nanos: number;
  };

  @ApiProperty()
  name: string;

  @ApiProperty()
  userId: string;

  @ApiProperty()
  categoryId: string;

  @ApiProperty()
  priorityId: string;

  @ApiProperty()
  description: string;
}
