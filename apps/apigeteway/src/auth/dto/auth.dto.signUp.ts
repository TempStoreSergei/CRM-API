import { ApiProperty } from '@nestjs/swagger';
import { RegisterDto } from '@app/common';

export class SignUpDtoClass implements RegisterDto {
  @ApiProperty()
  readonly email: string;

  @ApiProperty()
  readonly password: string;
}
