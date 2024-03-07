import { ApiProperty } from '@nestjs/swagger';
import { EventDelete } from '@app/common';

export class EventDeleteDto implements EventDelete {
  @ApiProperty()
  id: string;
}
