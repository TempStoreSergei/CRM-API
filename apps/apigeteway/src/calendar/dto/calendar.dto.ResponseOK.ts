import { ApiProperty } from '@nestjs/swagger';
import { ResponseOK } from '@app/common';

export class ResponseOKDto implements ResponseOK {
  @ApiProperty()
  result: string;
}
