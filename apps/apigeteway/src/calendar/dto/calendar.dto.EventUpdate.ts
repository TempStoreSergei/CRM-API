import { ApiProperty } from '@nestjs/swagger';
import { EventUpdate } from '@app/common';
export class EventUpdateDto implements EventUpdate {
  @ApiProperty()
  eventId: string;

  @ApiProperty()
  updatedFields: Record<string, any>;

  // Add other properties from the EventUpdate interface
}
